package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/parser"
)

func EvalReader(reader io.Reader) error {
	l := lexer.NewLexerFromReader(reader)
	p := parser.NewParser(l)
	program, err := p.Parse()
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", program.String(""))

	return nil
}

func EvalString(code string) error {
	return EvalReader(strings.NewReader(code))
}
