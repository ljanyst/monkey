package evaluator

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/parser"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type ObjectType object.go
//go:generate go run golang.org/x/tools/cmd/stringer -type ExitType object.go

type ObjectType int

const (
	INT ObjectType = iota
	BOOL
	STRING
	EXIT
	FUNCTION
	NIL
	RUNE
	ARRAY
)

type Object interface {
	Inspect() string
	Type() ObjectType
}

type IntObject struct {
	Value int64
}

type BoolObject struct {
	Value bool
}

type StringObject struct {
	Value []rune
}

type RuneObject struct {
	Value rune
}

type ExitType int

const (
	RETURN ExitType = iota
	BREAK
	CONTINUE
)

type ExitObject struct {
	Kind  ExitType
	Value Object
}

type BuiltInFunction func([]Object) (Object, error)

type FunctionObject struct {
	Params        []string
	ParentContext *Context
	Value         parser.Node
	BuiltIn       BuiltInFunction
}

type NilObject struct {
}

type ArrayObject struct {
	Value []Object
}

func (o *IntObject) Inspect() string {
	return fmt.Sprintf("%d", o.Value)
}

func (o *IntObject) Type() ObjectType {
	return INT
}

func (o *BoolObject) Inspect() string {
	if o.Value {
		return "true"
	}
	return "false"
}

func (o *BoolObject) Type() ObjectType {
	return BOOL
}

func (o *StringObject) Inspect() string {
	return fmt.Sprintf("%q", string(o.Value))
}

func (o *StringObject) Type() ObjectType {
	return STRING
}

func (o *RuneObject) Inspect() string {
	return fmt.Sprintf("%q", o.Value)
}

func (o *RuneObject) Type() ObjectType {
	return RUNE
}

func (o *ExitObject) Inspect() string {
	return fmt.Sprintf("return %q", o.Value.Inspect())
}

func (o *ExitObject) Type() ObjectType {
	return EXIT
}

func (o *FunctionObject) Inspect() string {
	if o.BuiltIn != nil {
		return "builtin(...)"
	}
	var sb strings.Builder
	sb.WriteString("fn(")
	for i, param := range o.Params {
		sb.WriteString(param)
		if i < len(o.Params)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (o *FunctionObject) Type() ObjectType {
	return FUNCTION
}

func (o *NilObject) Inspect() string {
	return "nil"
}

func (o *NilObject) Type() ObjectType {
	return NIL
}

func (o *ArrayObject) Inspect() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, item := range o.Value {
		sb.WriteString(item.Inspect())
		if i < len(o.Value)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (o *ArrayObject) Type() ObjectType {
	return ARRAY
}
