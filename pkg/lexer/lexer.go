package lexer

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/ej-shafran/gompiler/pkg/location"
	"github.com/ej-shafran/gompiler/pkg/token"
)

type ParseError struct {
	location location.Location
	error    error
}

func (err *ParseError) Error() string {
	line, offset := err.location.LineAndOffset()
	return fmt.Sprintf("%s:%d:%d: %s", err.location.FileInfo.FileName, line, offset, err.error)
}

func NewParseError(location location.Location, err error) *ParseError {
	return &ParseError{location: location, error: err}
}

var UNEXPECTED_END_OF_FILE = errors.New("Unexpected end of file")

type Lexer struct {
	cursor   int
	fileInfo location.FileInfo
}

func NewLexer(fileName, contents string) *Lexer {
	return &Lexer{cursor: 0, fileInfo: location.FileInfo{FileName: fileName, Contents: contents}}
}

func (l *Lexer) location() location.Location {
	return location.Location{Cursor: l.cursor, FileInfo: l.fileInfo}
}

func (l *Lexer) todo(s string) *ParseError {
	return NewParseError(l.location(), fmt.Errorf("TODO: %s", s))
}

func (l *Lexer) peekCharacter() (c rune, eof bool) {
	if l.cursor >= len(l.fileInfo.Contents) {
		return 0, true
	}

	return rune(l.fileInfo.Contents[l.cursor]), false
}

func (l *Lexer) consumeCharacter() {
	l.cursor++
}

func (l *Lexer) ConsumeToken() (*token.Token, *ParseError) {
	start := l.location()
	for {
		c, eof := l.peekCharacter()
		if eof {
			return token.NewToken(token.TOKEN_END_OF_FILE, start, start), nil
		}

		l.consumeCharacter()

		if unicode.IsSpace(c) {
			start = l.location()
			continue
		}

		// Symbols which can only appear on their own
		if c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}' || c == ';' || c == ',' || c == '.' {
			return token.NewToken(token.TOKEN_SYMBOL, start, l.location()), nil
		}

		// Symbols which can appear on their own or be followed with `=`
		if c == '!' || c == '%' || c == '^' || c == '*' || c == '/' || c == '=' {
			c, eof = l.peekCharacter()
			if eof {
				return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
			}

			if c == '=' {
				l.consumeCharacter()
			}

			return token.NewToken(token.TOKEN_SYMBOL, start, l.location()), nil
		}

		// Symbols which can appear:
		// - on their own
		// - doubled
		// - followed by `=`
		if c == '&' || c == '-' || c == '+' {
			c2, eof := l.peekCharacter()
			if eof {
				return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
			}

			// Special case for `->`
			if c2 == '=' || c2 == c || (c == '-' && c2 == '>') {
				l.consumeCharacter()
			}

			return token.NewToken(token.TOKEN_SYMBOL, start, l.location()), nil
		}

		// Symbols which can appear:
		// - on their own
		// - followed by `=`
		// - doubled
		// - doubled AND followed by `=`
		if c == '<' || c == '>' {
			c2, eof := l.peekCharacter()
			if eof {
				return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
			}

			switch c2 {
			case '=':
				l.consumeCharacter()
			case c:
				l.consumeCharacter()
				c2, eof = l.peekCharacter()
				if c2 == '=' {
					l.consumeCharacter()
				}
			}

			return token.NewToken(token.TOKEN_SYMBOL, start, l.location()), nil
		}

		return nil, l.todo("ConsumeToken")
	}
}
