package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/token"
)

func EvalReader(reader io.Reader) error {
	l := lexer.NewLexerFromReader(reader)
	for tok := l.ReadToken(); tok.Type != token.EOF; tok = l.ReadToken() {
		fmt.Printf("%v %s (%d:%d)\n", tok.Type, tok.Literal, tok.Line, tok.Column)
	}
	return nil
}

func EvalString(code string) error {
	return EvalReader(strings.NewReader(code))
}
