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

func mkErrWrongType(exp, got ObjectType, node parser.Node) error {
	return fmt.Errorf("%s Eval error: Expected type %s, got %s for expression %q",
		node.Token().Location(), exp, got, node.String(""),
	)
}

func mkErrWrongTypeStr(exp string, got ObjectType, node parser.Node) error {
	return fmt.Errorf("%s Eval error: Expected type %s, got %s for expression %q",
		node.Token().Location(), exp, got, node.String(""),
	)
}

func mkErrWrongOpForType(tok lexer.Token, objType ObjectType) error {
	return fmt.Errorf("%s Eval error: Invalid operator %s for type %s", tok.Location(), tok.Literal, objType)
}

func mkErrIndexOutOfBounds(node parser.Node, value, first, last int64) error {
	tok := node.Token()
	return fmt.Errorf("%s Eval error: Index %q is out of bounds: %d, valid range [%d:%d]",
		tok.Location(), node.String(""), value, first, last,
	)
}

func mkErrSliceEmpty(node parser.Node) error {
	tok := node.Token()
	return fmt.Errorf("%s Eval error: Slicing empty container %q", tok.Location(), node.String(""))
}

func mkErrWrongToken(expected string, got lexer.Token) error {
	lit := got.Literal
	return fmt.Errorf("%s Eval error: Expected %s, got %q", got.Location(), expected, lit)
}

func evalBlock(node parser.Node, c *Context) (Object, error) {
	var obj Object
	var err error

	implicit := node.(*parser.BlockNode).Implicit
	if !implicit {
		c = c.ChildContext()
	}

	obj = &NilObject{}
	for _, n := range node.Children() {
		obj, err = EvalNode(n, c)
		if err != nil {
			return nil, err
		}

		if obj.Type() == EXIT {
			break
		}
	}

	if implicit && obj.Type() == EXIT {
		exitObj := obj.(*ExitObject)
		if exitObj.Kind == RETURN {
			return exitObj.Value, nil
		}
		return nil, fmt.Errorf("%s Eval error: %s exit statement outside of a loop context",
			node.Token().Location(), exitObj.Kind)
	}

	return obj, nil
}

func evalInt(node parser.Node, c *Context) (Object, error) {
	return &IntObject{node.(*parser.IntNode).Value}, nil
}

func evalString(node parser.Node, c *Context) (Object, error) {
	return &StringObject{[]rune(node.(*parser.StringNode).Value)}, nil
}

func evalBool(node parser.Node, c *Context) (Object, error) {
	return &BoolObject{node.(*parser.BoolNode).Value}, nil
}

func evalRune(node parser.Node, c *Context) (Object, error) {
	return &RuneObject{node.(*parser.RuneNode).Value}, nil
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

	tok := node.Token()

	if tok.Type == lexer.BANG {
		if obj.Type() != BOOL {
			return nil, mkErrWrongType(BOOL, obj.Type(), exp)
		}
		return &BoolObject{!obj.(*BoolObject).Value}, nil
	}

	if tok.Type == lexer.MINUS {
		if obj.Type() != INT {
			return nil, mkErrWrongType(INT, obj.Type(), exp)
		}
		return &IntObject{-obj.(*IntObject).Value}, nil
	}

	return nil, fmt.Errorf("%s Unrecognized token for prefix expression: %s", tok.Location(), tok.Literal)
}

func assignIdent(node *parser.InfixNode, c *Context) (Object, error) {
	identNode := node.Left.(*parser.IdentifierNode)

	obj, err := EvalNode(node.Right, c)
	if err != nil {
		return nil, err
	}

	err = c.Set(identNode.Value, obj)
	if err != nil {
		return nil, fmt.Errorf("%s Eval error: %s", node.Token().Location(), err)
	}

	return obj, nil
}

func assignSlice(node *parser.InfixNode, c *Context) (Object, error) {
	sliceTok := node.Left.Token()
	slice, ok := node.Left.(*parser.SliceNode)
	if !ok {
		return nil, fmt.Errorf("%s Left hand side is not a slice: %s", sliceTok.Location(),
			node.Left.String(""))
	}

	if slice.End != nil {
		return nil, fmt.Errorf("%s Can only assign to a single element", sliceTok.Location())
	}

	subject, err := EvalNode(slice.Subject, c)
	if err != nil {
		return nil, err
	}

	indexObj, err := EvalNode(slice.Start, c)
	if err != nil {
		return nil, err
	}

	if subject.Type() != STRING && subject.Type() != ARRAY {
		return nil, mkErrWrongTypeStr("STRING or ARRAY", subject.Type(), slice.Subject)
	}

	if indexObj.Type() != INT {
		return nil, mkErrWrongType(INT, indexObj.Type(), slice.Start)
	}

	index := indexObj.(*IntObject).Value
	length := objLen(subject)
	if index < 0 || index >= length {
		return nil, mkErrIndexOutOfBounds(slice.Start, index, 0, length-1)
	}

	rhs, err := EvalNode(node.Right, c)
	if err != nil {
		return nil, err
	}

	switch subject.Type() {
	case STRING:
		if rhs.Type() != RUNE {
			return nil, mkErrWrongType(RUNE, rhs.Type(), node.Right)
		}
		str := subject.(*StringObject)
		str.Value[index] = rhs.(*RuneObject).Value
	case ARRAY:
		array := subject.(*ArrayObject)
		array.Value[index] = rhs
	}

	return rhs, nil
}

func evalAssign(node parser.Node, c *Context) (Object, error) {
	assignNode := node.(*parser.InfixNode)

	tok := assignNode.Left.Token()
	if tok.Type == lexer.IDENT {
		return assignIdent(assignNode, c)
	}
	return assignSlice(assignNode, c)
}

func evalInfixString(op lexer.Token, lVal, rVal []rune) (Object, error) {
	switch op.Type {
	case lexer.PLUS:
		return &StringObject{append(append([]rune{}, lVal...), rVal...)}, nil
	case lexer.LT:
		return &BoolObject{compareStrings(lVal, rVal) < 0}, nil
	case lexer.LE:
		return &BoolObject{compareStrings(lVal, rVal) <= 0}, nil
	case lexer.GT:
		return &BoolObject{compareStrings(lVal, rVal) > 0}, nil
	case lexer.GE:
		return &BoolObject{compareStrings(lVal, rVal) >= 0}, nil
	case lexer.EQ:
		return &BoolObject{compareStrings(lVal, rVal) == 0}, nil
	case lexer.NOT_EQ:
		return &BoolObject{compareStrings(lVal, rVal) != 0}, nil
	}

	return nil, mkErrWrongOpForType(op, STRING)
}

func evalInfixRune(op lexer.Token, lVal, rVal rune) (Object, error) {
	switch op.Type {
	case lexer.LT:
		return &BoolObject{compareRunes(lVal, rVal) < 0}, nil
	case lexer.LE:
		return &BoolObject{compareRunes(lVal, rVal) <= 0}, nil
	case lexer.GT:
		return &BoolObject{compareRunes(lVal, rVal) > 0}, nil
	case lexer.GE:
		return &BoolObject{compareRunes(lVal, rVal) >= 0}, nil
	case lexer.EQ:
		return &BoolObject{compareRunes(lVal, rVal) == 0}, nil
	case lexer.NOT_EQ:
		return &BoolObject{compareRunes(lVal, rVal) != 0}, nil
	}

	return nil, mkErrWrongOpForType(op, STRING)
}

func evalInfixArray(op lexer.Token, lVal, rVal []Object) (Object, error) {
	if op.Type == lexer.PLUS {
		return &ArrayObject{
			append(append([]Object{}, lVal...), rVal...),
		}, nil
	}

	return nil, mkErrWrongOpForType(op, STRING)
}

func evalInfixBool(op lexer.Token, lVal, rVal bool) (Object, error) {
	switch op.Type {
	case lexer.AND:
		return &BoolObject{lVal && rVal}, nil
	case lexer.OR:
		return &BoolObject{lVal || rVal}, nil
	default:
		return nil, mkErrWrongOpForType(op, BOOL)
	}
}

func evalInfixInt(op lexer.Token, lVal, rVal int64) (Object, error) {
	switch op.Type {
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

	return nil, mkErrWrongOpForType(op, INT)
}

func evalInfix(node parser.Node, c *Context) (Object, error) {
	iNode := node.(*parser.InfixNode)
	tok := node.Token()

	if tok.Type == lexer.ASSIGN {
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

	if left.Type() != INT && left.Type() != STRING && left.Type() != ARRAY &&
		left.Type() != BOOL && left.Type() != RUNE {
		return nil, mkErrWrongTypeStr("INT or STRING or ARRAY or BOOL or RUNE", left.Type(), iNode.Left)
	}

	if right.Type() != right.Type() {
		return nil, mkErrWrongType(left.Type(), right.Type(), iNode.Right)
	}

	switch left.Type() {
	case STRING:
		return evalInfixString(tok, left.(*StringObject).Value, right.(*StringObject).Value)
	case ARRAY:
		return evalInfixArray(tok, left.(*ArrayObject).Value, right.(*ArrayObject).Value)
	case INT:
		return evalInfixInt(tok, left.(*IntObject).Value, right.(*IntObject).Value)
	case BOOL:
		return evalInfixBool(tok, left.(*BoolObject).Value, right.(*BoolObject).Value)
	case RUNE:
		return evalInfixRune(tok, left.(*RuneObject).Value, right.(*RuneObject).Value)
	default:
		return nil, fmt.Errorf("%s Eval error: No infix eval function for type %s",
			tok.Location(), left.Type())
	}
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
	return &ExitObject{RETURN, obj}, nil
}

func evalStatement(node parser.Node, c *Context) (Object, error) {
	tok := node.Token()
	switch tok.Type {
	case lexer.LET:
		return evalLet(node, c)
	case lexer.RETURN:
		return evalReturn(node, c)
	case lexer.BREAK:
		return &ExitObject{BREAK, nil}, nil
	case lexer.CONTINUE:
		return &ExitObject{CONTINUE, nil}, nil
	}
	return nil, fmt.Errorf("%s Eval error: Unrecognized statement: %s", tok.Location(), tok.Literal)
}

func evalConditional(node parser.Node, c *Context) (Object, error) {
	condNode := node.(*parser.ConditionalNode)

	condObj, err := EvalNode(condNode.Condition, c)
	if err != nil {
		return nil, err
	}

	if condObj.Type() != BOOL {
		return nil, mkErrWrongType(BOOL, condObj.Type(), condNode.Condition)
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
	return &FunctionObject{params, c, funcNode.Body, nil}, nil
}

func evalFunctionCall(node parser.Node, c *Context) (Object, error) {
	funcCallNode := node.(*parser.FunctionCallNode)

	fObj, err := EvalNode(funcCallNode.Function, c)
	if err != nil {
		return nil, err
	}

	if fObj.Type() != FUNCTION {
		return nil, mkErrWrongType(FUNCTION, fObj.Type(), funcCallNode.Function)
	}

	f := fObj.(*FunctionObject)

	if f.Params != nil && len(f.Params) != len(funcCallNode.Args) {
		return nil, fmt.Errorf("%s Eval error: Expected %d params, got %d",
			funcCallNode.Token().Location(), len(f.Params), len(funcCallNode.Args))
	}

	params := []Object{}
	for _, paramNode := range funcCallNode.Args {
		paramObj, err := EvalNode(paramNode, c)
		if err != nil {
			return nil, err
		}
		params = append(params, paramObj)
	}

	if f.BuiltIn != nil {
		obj, err := f.BuiltIn(params)
		if err != nil {
			return nil, fmt.Errorf("%s Eval error: Expression %q: %s", node.Token().Location(),
				node.String(""), err)
		}
		return obj, nil
	}

	paramContext := f.ParentContext.ChildContext()
	for i, paramName := range f.Params {
		paramContext.Create(paramName, params[i])
	}

	funcCallContext := paramContext.ChildContext()

	retObj, err := EvalNode(f.Value, funcCallContext)
	if err != nil {
		return nil, err
	}

	if retObj.Type() == EXIT {
		exitObj := retObj.(*ExitObject)
		if exitObj.Kind == RETURN {
			return exitObj.Value, nil
		}
		return nil, fmt.Errorf("%s Eval error: %s exit statement outside of a loop context",
			f.Value.Token().Location(), exitObj.Kind)
	}
	return retObj, nil
}

func objLen(obj Object) int64 {
	if obj.Type() == STRING {
		return int64(len(obj.(*StringObject).Value))
	}

	if obj.Type() == ARRAY {
		return int64(len(obj.(*ArrayObject).Value))
	}
	return 0
}

func objRange(obj Object, start, end int64) Object {
	if obj.Type() == STRING {
		return &StringObject{obj.(*StringObject).Value[start:end]}
	}

	if obj.Type() == ARRAY {
		return &ArrayObject{obj.(*ArrayObject).Value[start:end]}
	}

	return &NilObject{}
}

func objItem(obj Object, index int64) Object {
	if obj.Type() == STRING {
		return &RuneObject{obj.(*StringObject).Value[index]}
	}

	if obj.Type() == ARRAY {
		return obj.(*ArrayObject).Value[index]
	}

	return &NilObject{}
}

func evalSlice(node parser.Node, c *Context) (Object, error) {
	sliceNode := node.(*parser.SliceNode)

	sliceObj, err := EvalNode(sliceNode.Subject, c)
	if err != nil {
		return nil, err
	}

	if sliceObj.Type() != STRING && sliceObj.Type() != ARRAY {
		return nil, mkErrWrongTypeStr("STRING or ARRAY", sliceObj.Type(), sliceNode.Subject)
	}

	length := objLen(sliceObj)

	if length == 0 {
		return nil, mkErrSliceEmpty(sliceNode.Subject)
	}

	startObj, err := EvalNode(sliceNode.Start, c)
	if err != nil {
		return nil, err
	}

	if startObj.Type() != INT {
		return nil, mkErrWrongType(INT, startObj.Type(), sliceNode.Start)
	}
	start := startObj.(*IntObject).Value

	if start < 0 || start >= length {
		return nil, mkErrIndexOutOfBounds(sliceNode.Start, start, 0, length-1)
	}

	if sliceNode.End != nil {
		endObj, err := EvalNode(sliceNode.End, c)
		if err != nil {
			return nil, err
		}

		if endObj.Type() != INT {
			return nil, mkErrWrongType(INT, endObj.Type(), sliceNode.End)
		}

		end := endObj.(*IntObject).Value
		if end < 0 || end > length {
			return nil, mkErrIndexOutOfBounds(sliceNode.End, end, 0, length)
		}

		if start > end {
			return nil, mkErrIndexOutOfBounds(sliceNode.End, end, start, length)
		}

		return objRange(sliceObj, start, end), nil
	}

	return objItem(sliceObj, start), nil
}

func evalArray(node parser.Node, c *Context) (Object, error) {
	arrayNode := node.(*parser.ArrayNode)
	var objects []Object
	for _, node := range arrayNode.Items {
		obj, err := EvalNode(node, c)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}
	return &ArrayObject{objects}, nil
}

func evalLoop(node parser.Node, c *Context) (Object, error) {
	loopNode := node.(*parser.LoopNode)
	cLoop := c.ChildContext()
	cInner := cLoop.ChildContext()

	if loopNode.Initializer != nil {
		_, err := EvalNode(loopNode.Initializer, cLoop)
		if err != nil {
			return nil, err
		}
	}

	var retObject Object
	retObject = &NilObject{}

	for {
		condObj, err := EvalNode(loopNode.Condition, cLoop)
		if err != nil {
			return nil, err
		}

		if condObj.Type() != BOOL {
			return nil, mkErrWrongType(BOOL, condObj.Type(), loopNode.Condition)
		}

		if condObj.(*BoolObject).Value == false {
			return retObject, nil
		}

		retObject, err = EvalNode(loopNode.Body, cInner)
		if err != nil {
			return nil, err
		}

		if retObject.Type() == EXIT {
			obj := retObject.(*ExitObject)
			if obj.Kind == CONTINUE || obj.Kind == BREAK {
				retObject = &NilObject{}
			}

			if obj.Kind == BREAK || obj.Kind == RETURN {
				break
			}
		}

		if loopNode.Modifier != nil {
			_, err := EvalNode(loopNode.Modifier, cLoop)
			if err != nil {
				return nil, err
			}
		}
	}

	return retObject, nil
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
	case *parser.RuneNode:
		return evalRune(node, c)
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
	case *parser.SliceNode:
		return evalSlice(node, c)
	case *parser.ArrayNode:
		return evalArray(node, c)
	case *parser.LoopNode:
		return evalLoop(node, c)
	default:
		return nil,
			fmt.Errorf("%s Eval error: Evaluator not implemented for %s",
				node.Token().Location(),
				reflect.ValueOf(node).Elem().Type(),
			)
	}
}
