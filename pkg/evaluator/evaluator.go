package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/ljanyst/edu-interp/pkg/lexer"
	"github.com/ljanyst/edu-interp/pkg/token"
)

func EvalReader(reader io.Reader) error {
	l := lexer.NewLexerFromReader(reader)
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%v %s (%d:%d)\n", tok.Type, tok.Literal, tok.Line, tok.Column)
	}
	return nil
}

func EvalString(code string) error {
	return EvalReader(strings.NewReader(code))
}
