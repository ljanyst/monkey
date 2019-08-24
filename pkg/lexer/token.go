package lexer

import (
	"fmt"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type TokenType token.go

type TokenType int

const (
	LET TokenType = iota
	IDENT
	ASSIGN
	INT
	SEMICOLON
	FUNCTION
	LPAREN
	COMMA
	RPAREN
	LBRACE
	PLUS
	RBRACE
	BANG
	MINUS
	SLASH
	ASTERISK
	LT
	LE
	GT
	GE
	IF
	RETURN
	TRUE
	ELSE
	FALSE
	STRING
	EQ
	NOT_EQ
	INVALID
	BLOCK
	EOF
)

type Token struct {
	Type     TokenType
	Literal  string
	Line     uint32
	Column   uint32
	FileName *string
}

func (t Token) Location() string {
	return fmt.Sprintf("[%s:%d:%d]", *t.FileName, t.Line, t.Column)
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupKeyword(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
