package evaluator

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/parser"
)

func EvalReader(reader io.Reader, c *Context, name string) (Object, error) {
	l := lexer.NewLexerFromReader(reader, name)
	p := parser.NewParser(l)
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	return EvalNode(program, c)
}

func EvalString(code string, c *Context, name string) (Object, error) {
	return EvalReader(strings.NewReader(code), c, name)
}

func mkErrUnexpectedType(exp, got ObjectType, node parser.Node) error {
	return fmt.Errorf("%s Eval error: Expected type %s, got %s for expression %q",
		node.Token().Location(), exp, got, node.String(""),
	)
}

func mkErrWrongToken(expected string, got lexer.Token) error {
	lit := got.Literal
	return fmt.Errorf("%s Eval error: expected %s, got %q", got.Location(), expected, lit)
}

func evalBlock(node parser.Node, c *Context) (Object, error) {
	var obj Object
	var err error

	implicit := node.(*parser.BlockNode).Implicit
	if !implicit {
		c = c.ChildContext()
	}

	for _, n := range node.Children() {
		obj, err = EvalNode(n, c)
		if err != nil {
			return nil, err
		}

		if obj.Type() == RETURN {
			break
		}
	}

	if implicit && obj.Type() == RETURN {
		return obj.(*ReturnObject).Value, nil
	}

	return obj, nil
}

func evalInt(node parser.Node, c *Context) (Object, error) {
	return &IntObject{node.(*parser.IntNode).Value}, nil
}

func evalString(node parser.Node, c *Context) (Object, error) {
	return &StringObject{node.(*parser.StringNode).Value}, nil
}

func evalBool(node parser.Node, c *Context) (Object, error) {
	return &BoolObject{node.(*parser.BoolNode).Value}, nil
}

func evalIdentifier(node parser.Node, c *Context) (Object, error) {
	identNode := node.(*parser.IdentifierNode)
	obj, err := c.Resolve(identNode.Value)
	if err != nil {
		tok := node.Token()
		return nil, fmt.Errorf("%s Eval error: %s", tok.Location(), err)
	}
	return obj, nil
}

func evalPrefix(node parser.Node, c *Context) (Object, error) {
	exp := node.(*parser.PrefixNode).Expression
	obj, err := EvalNode(exp, c)
	if err != nil {
		return nil, err
	}

	if node.Token().Type == lexer.BANG {
		if obj.Type() != BOOL {
			return nil, mkErrUnexpectedType(BOOL, obj.Type(), exp)
		}
		return &BoolObject{!obj.(*BoolObject).Value}, nil
	}

	if node.Token().Type == lexer.MINUS {
		if obj.Type() != INT {
			return nil, mkErrUnexpectedType(INT, obj.Type(), exp)
		}
		return &IntObject{-obj.(*IntObject).Value}, nil
	}

	return nil, fmt.Errorf("Unrecognized token for prefix expression: %s", node.Token().Literal)
}

func evalAssign(node parser.Node, c *Context) (Object, error) {
	assignNode := node.(*parser.InfixNode)
	tok := assignNode.Left.Token()
	if tok.Type != lexer.IDENT {
		return nil, mkErrWrongToken("identifier", tok)
	}

	identNode := assignNode.Left.(*parser.IdentifierNode)

	obj, err := EvalNode(assignNode.Right, c)
	if err != nil {
		return nil, err
	}

	err = c.Set(identNode.Value, obj)
	if err != nil {
		return nil, fmt.Errorf("%s Eval error: %s", tok.Location(), err)
	}

	return obj, nil
}

func evalInfix(node parser.Node, c *Context) (Object, error) {
	iNode := node.(*parser.InfixNode)

	if node.Token().Type == lexer.ASSIGN {
		return evalAssign(node, c)
	}

	left, err := EvalNode(iNode.Left, c)
	if err != nil {
		return nil, err
	}

	right, err := EvalNode(iNode.Right, c)
	if err != nil {
		return nil, err
	}

	if left.Type() != INT {
		return nil, mkErrUnexpectedType(INT, left.Type(), iNode.Left)
	}

	if right.Type() != INT {
		return nil, mkErrUnexpectedType(INT, right.Type(), iNode.Right)
	}

	lVal := left.(*IntObject).Value
	rVal := right.(*IntObject).Value

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

func evalLet(node parser.Node, c *Context) (Object, error) {
	child := node.Children()[0]
	tok := child.Token()
	if tok.Type != lexer.ASSIGN {
		return nil, mkErrWrongToken("assignment", tok)
	}

	assignNode := child.(*parser.InfixNode)
	tok = assignNode.Left.Token()
	if tok.Type != lexer.IDENT {
		return nil, mkErrWrongToken("identifier", tok)
	}

	identNode := assignNode.Left.(*parser.IdentifierNode)

	obj, err := EvalNode(assignNode.Right, c)
	if err != nil {
		return nil, err
	}

	err = c.Create(identNode.Value, obj)
	if err != nil {
		return nil, fmt.Errorf("%s Eval error: %s", tok.Location(), err)
	}

	return obj, nil
}

func evalReturn(node parser.Node, c *Context) (Object, error) {
	obj, err := EvalNode(node.Children()[0], c)
	if err != nil {
		return nil, err
	}
	return &ReturnObject{obj}, nil
}

func evalStatement(node parser.Node, c *Context) (Object, error) {
	if node.Token().Type == lexer.LET {
		return evalLet(node, c)
	} else if node.Token().Type == lexer.RETURN {
		return evalReturn(node, c)
	}
	return nil, fmt.Errorf("Unrecognized statement: %s", node.Token().Literal)
}

func evalConditional(node parser.Node, c *Context) (Object, error) {
	condNode := node.(*parser.ConditionalNode)

	condObj, err := EvalNode(condNode.Condition, c)
	if err != nil {
		return nil, err
	}

	if condObj.Type() != BOOL {
		return nil, mkErrUnexpectedType(BOOL, condObj.Type(), condNode.Condition)
	}

	if condObj.(*BoolObject).Value {
		return EvalNode(condNode.Consequent, c)
	}

	return EvalNode(condNode.Alternative, c)
}

func evalFunction(node parser.Node, c *Context) (Object, error) {
	funcNode := node.(*parser.FunctionNode)
	params := []string{}

	for _, param := range funcNode.Params {
		tok := param.Token()
		if tok.Type != lexer.IDENT {
			return nil, mkErrWrongToken("identifier", tok)
		}
		params = append(params, param.(*parser.IdentifierNode).Value)
	}
	return &FunctionObject{params, c, funcNode.Body}, nil
}

func evalFunctionCall(node parser.Node, c *Context) (Object, error) {
	funcCallNode := node.(*parser.FunctionCallNode)

	tok := funcCallNode.Name.Token()
	if tok.Type != lexer.IDENT {
		return nil, mkErrWrongToken("identifier", tok)
	}
	name := funcCallNode.Name.(*parser.IdentifierNode).Value

	fObj, err := c.Resolve(name)
	if err != nil {
		return nil, err
	}

	if fObj.Type() != FUNCTION {
		return nil, mkErrUnexpectedType(FUNCTION, fObj.Type(), funcCallNode.Name)
	}

	f := fObj.(*FunctionObject)

	if len(f.Params) != len(funcCallNode.Args) {
		return nil, fmt.Errorf("%s Eval error: Expected %d params for %q, got %d",
			funcCallNode.Token().Location(), len(f.Params), name, len(funcCallNode.Args))
	}

	paramContext := f.ParentContext.ChildContext()
	for i, paramName := range f.Params {
		paramObj, err := EvalNode(funcCallNode.Args[i], c)
		if err != nil {
			return nil, err
		}
		paramContext.Create(paramName, paramObj)
	}

	funcCallContext := paramContext.ChildContext()

	retObj, err := EvalNode(f.Value, funcCallContext)
	if err != nil {
		return nil, err
	}

	if retObj.Type() == RETURN {
		return retObj.(*ReturnObject).Value, nil
	}
	return retObj, nil
}

func EvalNode(node parser.Node, c *Context) (Object, error) {
	if node == nil {
		return &NilObject{}, nil
	}

	switch node.(type) {
	case *parser.BlockNode:
		return evalBlock(node, c)
	case *parser.IntNode:
		return evalInt(node, c)
	case *parser.StringNode:
		return evalString(node, c)
	case *parser.BoolNode:
		return evalBool(node, c)
	case *parser.NilNode:
		return &NilObject{}, nil
	case *parser.IdentifierNode:
		return evalIdentifier(node, c)
	case *parser.PrefixNode:
		return evalPrefix(node, c)
	case *parser.InfixNode:
		return evalInfix(node, c)
	case *parser.StatementNode:
		return evalStatement(node, c)
	case *parser.ConditionalNode:
		return evalConditional(node, c)
	case *parser.FunctionNode:
		return evalFunction(node, c)
	case *parser.FunctionCallNode:
		return evalFunctionCall(node, c)
	default:
		return nil,
			fmt.Errorf("%s Eval error: Evaluator not implemented for %s",
				node.Token().Location(),
				reflect.ValueOf(node).Elem().Type(),
			)
	}
}
