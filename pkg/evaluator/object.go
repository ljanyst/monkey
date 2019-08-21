package evaluator

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/parser"
)

type ObjectType string

const (
	INT      = "INT"
	BOOL     = "BOOL"
	STRING   = "STRING"
	RETURN   = "RETURN"
	FUNCTION = "FUNCTION"
)

type Object interface {
	Inspect() string
	Type() ObjectType
	Value() interface{}
}

type IntObject struct {
	value int64
}

type BoolObject struct {
	value bool
}

type StringObject struct {
	value string
}

type ReturnObject struct {
	value Object
}

type FunctionObject struct {
	Params        []string
	ParentContext *Context
	Block         parser.Node
}

func (o *IntObject) Inspect() string {
	return fmt.Sprintf("%d", o.value)
}

func (o *IntObject) Type() ObjectType {
	return INT
}

func (o *IntObject) Value() interface{} {
	return o.value
}

func (o *BoolObject) Inspect() string {
	if o.value {
		return "true"
	}
	return "false"
}

func (o *BoolObject) Type() ObjectType {
	return BOOL
}

func (o *BoolObject) Value() interface{} {
	return o.value
}

func (o *StringObject) Inspect() string {
	return fmt.Sprintf("%q", o.value)
}

func (o *StringObject) Type() ObjectType {
	return STRING
}

func (o *StringObject) Value() interface{} {
	return o.value
}

func (o *ReturnObject) Inspect() string {
	return fmt.Sprintf("return %q", o.value.Inspect())
}

func (o *ReturnObject) Type() ObjectType {
	return RETURN
}

func (o *ReturnObject) Value() interface{} {
	return o.value
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

func (o *FunctionObject) Value() interface{} {
	return o.Block
}
