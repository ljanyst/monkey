package evaluator

import (
	"fmt"
)

type Context struct {
	bindings map[string]Object
	parent   *Context
}

func (c *Context) Resolve(name string) (Object, error) {
	if obj, ok := c.bindings[name]; ok {
		return obj, nil
	}
	if c.parent != nil {
		return c.parent.Resolve(name)
	}
	return nil, fmt.Errorf("Variable %q not defined", name)
}

func (c *Context) Create(name string, obj Object) error {
	if _, ok := c.bindings[name]; ok {
		return fmt.Errorf("Unable to create variable: %q already exist", name)
	}
	c.bindings[name] = obj
	return nil
}

func (c *Context) Set(name string, obj Object) error {
	if _, ok := c.bindings[name]; !ok {
		if c.parent != nil {
			return c.parent.Set(name, obj)
		}
		return fmt.Errorf("Unable to set variable: %q does not exist", name)
	}
	c.bindings[name] = obj
	return nil
}

func (c *Context) ChildContext() *Context {
	child := NewContext()
	child.parent = c
	return child
}

func NewContext() *Context {
	c := new(Context)
	c.bindings = make(map[string]Object)
	return c
}
