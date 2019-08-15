package evaluator

import (
	"fmt"
)

type ObjectType string

const (
	INT    = "INT"
	BOOL   = "BOOL"
	STRING = "STRING"
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
