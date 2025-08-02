package token

import "github.com/ej-shafran/gompiler/pkg/location"

type TokenKind int

const (
	TOKEN_END_OF_FILE TokenKind = iota
	TOKEN_SYMBOL
	TOKEN_IDENTIFIER
	TOKEN_QUOTED_STRING
	TOKEN_QUOTED_CHARACTER
	TOKEN_NUMBER_LITERAL
)

type Token struct {
	Kind  TokenKind
	Start location.Location
	End   location.Location
}

func NewToken(kind TokenKind, start, end location.Location) *Token {
	return &Token{Kind: kind, Start: start, End: end}
}
