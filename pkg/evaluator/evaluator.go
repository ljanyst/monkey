package evaluator

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/parser"
)

func EvalReader(reader io.Reader) (Object, error) {
	l := lexer.NewLexerFromReader(reader)
	p := parser.NewParser(l)
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return EvalNode(program)
}

func EvalString(code string) (Object, error) {
	return EvalReader(strings.NewReader(code))
}

func mkErrUnexpectedType(exp, got ObjectType, node parser.Node) error {
	return fmt.Errorf("Expected type %s got %s for expression %q at line %d",
		exp, got, node.String(""), node.Token().Line,
	)
}

func evalBlock(node parser.Node) (Object, error) {
	var obj Object
	var err error
	for _, n := range node.Children() {
		obj, err = EvalNode(n)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func evalInt(node parser.Node) (Object, error) {
	return &IntObject{node.(*parser.IntNode).Value}, nil
}

func evalString(node parser.Node) (Object, error) {
	return &StringObject{node.(*parser.StringNode).Value}, nil
}

func evalBool(node parser.Node) (Object, error) {
	return &BoolObject{node.(*parser.BoolNode).Value}, nil
}

func evalPrefix(node parser.Node) (Object, error) {
	exp := node.(*parser.PrefixNode).Expression
	obj, err := EvalNode(exp)
	if err != nil {
		return nil, err
	}

	if node.Token().Type == lexer.BANG {
		if obj.Type() != BOOL {
			return nil, mkErrUnexpectedType(BOOL, obj.Type(), exp)
		}
		return &BoolObject{!obj.Value().(bool)}, nil
	}

	if node.Token().Type == lexer.MINUS {
		if obj.Type() != INT {
			return nil, mkErrUnexpectedType(INT, obj.Type(), exp)
		}
		return &IntObject{-obj.Value().(int64)}, nil
	}

	return nil, fmt.Errorf("Unrecognized token for prefix expression: %s", node.Token().Literal)
}

func evalInfix(node parser.Node) (Object, error) {
	iNode := node.(*parser.InfixNode)

	left, err := EvalNode(iNode.Left)
	if err != nil {
		return nil, err
	}

	right, err := EvalNode(iNode.Right)
	if err != nil {
		return nil, err
	}

	if left.Type() != INT {
		return nil, mkErrUnexpectedType(INT, left.Type(), iNode.Left)
	}

	if right.Type() != INT {
		return nil, mkErrUnexpectedType(INT, right.Type(), iNode.Right)
	}

	lVal := left.Value().(int64)
	rVal := right.Value().(int64)

	switch node.Token().Type {
	case lexer.PLUS:
		return &IntObject{lVal + rVal}, nil
	case lexer.MINUS:
		return &IntObject{lVal - rVal}, nil
	case lexer.SLASH:
		return &IntObject{lVal / rVal}, nil
	case lexer.ASTERISK:
		return &IntObject{lVal * rVal}, nil
	case lexer.LT:
		return &BoolObject{lVal < rVal}, nil
	case lexer.LE:
		return &BoolObject{lVal <= rVal}, nil
	case lexer.GT:
		return &BoolObject{lVal > rVal}, nil
	case lexer.GE:
		return &BoolObject{lVal >= rVal}, nil
	case lexer.EQ:
		return &BoolObject{lVal == rVal}, nil
	case lexer.NOT_EQ:
		return &BoolObject{lVal != rVal}, nil
	}

	return nil, fmt.Errorf("Unrecognized token for infix expression: %s", node.Token().Literal)
}

func EvalNode(node parser.Node) (Object, error) {
	switch node.(type) {
	case *parser.BlockNode:
		return evalBlock(node)
	case *parser.IntNode:
		return evalInt(node)
	case *parser.StringNode:
		return evalString(node)
	case *parser.BoolNode:
		return evalBool(node)
	case *parser.PrefixNode:
		return evalPrefix(node)
	case *parser.InfixNode:
		return evalInfix(node)
	default:
		return nil,
			fmt.Errorf("Evaluator not implemented for node type %s created for %s at (%d:%d)",
				reflect.ValueOf(node).Elem().Type(),
				node.Token().Literal,
				node.Token().Line,
				node.Token().Column,
			)
	}
}
