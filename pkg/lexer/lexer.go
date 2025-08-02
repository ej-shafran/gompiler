package lexer

import (
	"fmt"

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

		_ = c

		return nil, l.todo("ConsumeToken")
	}
}
