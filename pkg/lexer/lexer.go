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
var EXPECTED_EXPERSSION = errors.New("Expected experssion")

type Lexer struct {
	cursor   int
	fileInfo location.FileInfo
}

type LexerSnapshot struct {
	cursor int
}

func NewLexer(fileName, contents string) *Lexer {
	return &Lexer{cursor: 0, fileInfo: location.FileInfo{FileName: fileName, Contents: contents}}
}

func (l *Lexer) location() location.Location {
	return location.Location{Cursor: l.cursor, FileInfo: l.fileInfo}
}

func (l *Lexer) todo(s string, location *location.Location) *ParseError {
	if location == nil {
		loc := l.location()
		location = &loc
	}
	return NewParseError(*location, fmt.Errorf("TODO: %s", s))
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

func (l *Lexer) SaveSnapshot() LexerSnapshot {
	return LexerSnapshot{cursor: l.cursor}
}

func (l *Lexer) RestoreSnapshot(snap LexerSnapshot) {
	l.cursor = snap.cursor
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

		// Macros
		if c == '#' {
			escaping := false
			for {
				c, eof = l.peekCharacter()
				if eof {
					return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
				}

				l.consumeCharacter()

				if escaping {
					escaping = false
					continue
				} else if c != '\n' {
					escaping = c == '\\'
				} else {
					return token.NewToken(token.TOKEN_MACRO, start, l.location()), nil
				}
			}
		}

		// Comments
		if c == '/' {
			c2, eof := l.peekCharacter()
			if eof {
				return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
			}

			// Single-line
			if c2 == '/' {
				l.consumeCharacter()
				for {
					c2, eof = l.peekCharacter()
					if eof {
						return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
					}

					l.consumeCharacter()

					if c2 == '\n' {
						return token.NewToken(token.TOKEN_SINGLE_LINE_COMMENT, start, l.location()), nil
					}
				}
			}

			// Multi-line
			if c2 == '*' {
				l.consumeCharacter()

				lastStar := false
				for {
					c2, eof = l.peekCharacter()
					if eof {
						return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
					}

					l.consumeCharacter()

					if lastStar && c2 == '/' {
						return token.NewToken(token.TOKEN_MULTI_LINE_COMMENT, start, l.location()), nil
					} else {
						lastStar = c2 == '*'
					}
				}
			}
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

		// Double/single quotes
		if c == '"' || c == '\'' {
			escaping := false

			for {
				c2, eof := l.peekCharacter()
				if eof {
					return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
				}

				l.consumeCharacter()

				if escaping {
					escaping = false
					continue
				} else {
					escaping = c2 == '\\'
				}

				if c2 == c {
					var kind token.TokenKind
					if c == '"' {
						kind = token.TOKEN_QUOTED_STRING
					} else {
						kind = token.TOKEN_QUOTED_CHARACTER
					}

					return token.NewToken(kind, start, l.location()), nil
				}
			}
		}

		// Number literals
		if unicode.IsDigit(c) {
			// `0x` and `0b` prefixes
			if c == '0' {
				c, eof = l.peekCharacter()
				if eof {
					return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
				}

				if c == 'x' || c == 'b' {
					l.consumeCharacter()
				}
			}

			for {
				c, eof = l.peekCharacter()
				if eof {
					return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
				}

				if c == '.' || unicode.IsDigit(c) {
					l.consumeCharacter()
					continue
				}

				return token.NewToken(token.TOKEN_NUMBER_LITERAL, start, l.location()), nil
			}
		}

		// Identifiers
		if unicode.IsLetter(c) || c == '_' {
			for {
				c, eof = l.peekCharacter()
				if eof {
					return nil, NewParseError(l.location(), UNEXPECTED_END_OF_FILE)
				}

				if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
					l.consumeCharacter()
					continue
				}

				// TODO: logic for builtin keywords?
				return token.NewToken(token.TOKEN_IDENTIFIER, start, l.location()), nil
			}
		}

		return nil, l.todo("ConsumeToken", nil)
	}
}

func (l *Lexer) TokenValue(t *token.Token) string {
	return l.fileInfo.Contents[t.Start.Cursor:t.End.Cursor]
}

func (l *Lexer) ExpectTokenKind(kind token.TokenKind) (*token.Token, *ParseError) {
	t, err := l.ConsumeToken()
	if err != nil {
		return nil, err
	}

	if t.Kind != kind {
		return nil, l.todo("unexpected kind error", &t.Start)
	}

	return t, nil
}

func (l *Lexer) ExpectSymbol(symbol string) (*token.Token, *ParseError) {
	t, err := l.ExpectTokenKind(token.TOKEN_SYMBOL)
	if err != nil {
		return nil, err
	}

	if l.TokenValue(t) != symbol {
		return nil, l.todo("unexpected symbol error", &t.Start)
	}

	return t, nil
}

func (l *Lexer) ExpectIdentifier(identifier string) (*token.Token, *ParseError) {
	t, err := l.ExpectTokenKind(token.TOKEN_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	if l.TokenValue(t) != identifier {
		return nil, l.todo("unexpected identifier error", &t.Start)
	}

	return t, nil
}
