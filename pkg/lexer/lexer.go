package lexer

import (
	"io"

	"github.com/ljanyst/edu-interp/pkg/token"
)

type Lexer struct {
}

func NewLexerFromString(input string) *Lexer {
	return new(Lexer)
}

func NewLexerFromReader(input io.Reader) *Lexer {
	return new(Lexer)
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{token.EOF, "EOF", 0, 0}
}
