package lexer

import (
	"fmt"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type TokenType token.go

type TokenType int

const (
	NONE TokenType = iota
	LET
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
	NIL
	RUNE
	LBRACKET
	RBRACKET
	COLON
	FOR
	BREAK
	CONTINUE
	AND
	OR
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
	"let":      LET,
	"fn":       FUNCTION,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupKeyword(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
