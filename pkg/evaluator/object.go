package evaluator

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/parser"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type ObjectType object.go

type ObjectType int

const (
	INT ObjectType = iota
	BOOL
	STRING
	RETURN
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
	Value string
}

type RuneObject struct {
	Value rune
}

type ReturnObject struct {
	Value Object
}

type FunctionObject struct {
	Params        []string
	ParentContext *Context
	Value         parser.Node
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
	return fmt.Sprintf("%q", o.Value)
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

func (o *ReturnObject) Inspect() string {
	return fmt.Sprintf("return %q", o.Value.Inspect())
}

func (o *ReturnObject) Type() ObjectType {
	return RETURN
}

func (o *FunctionObject) Inspect() string {
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
