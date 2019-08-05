package parser

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/token"
)

type Node interface {
	String() string
	Children() []Node
	Token() token.Token
}

type ProgramNode struct {
	children []Node
}

type IntNode struct {
	token token.Token
	Value int64
}

type StringNode struct {
	token token.Token
	Value string
}

type IdentifierNode struct {
	token token.Token
	Value string
}

type BoolNode struct {
	token token.Token
	Value bool
}

type PrefixNode struct {
	token      token.Token
	expression Node
}

type InfixNode struct {
	token token.Token
	left  Node
	right Node
}

func (n *ProgramNode) String() string {
	var sb strings.Builder
	for _, node := range n.children {
		sb.WriteString(node.String())
	}
	return sb.String()
}

func (n *ProgramNode) Children() []Node {
	return n.children
}

func (n *ProgramNode) Token() token.Token {
	return token.Token{"", "PROGRAM", 0, 0}
}

func (n *IntNode) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *IntNode) Children() []Node {
	return []Node{}
}

func (n *IntNode) Token() token.Token {
	return n.token
}

func (n *StringNode) String() string {
	return n.Value
}

func (n *StringNode) Children() []Node {
	return []Node{}
}

func (n *StringNode) Token() token.Token {
	return n.token
}

func (n *IdentifierNode) String() string {
	return n.Value
}

func (n *IdentifierNode) Children() []Node {
	return []Node{}
}

func (n *IdentifierNode) Token() token.Token {
	return n.token
}

func (n *BoolNode) String() string {
	if n.Value {
		return "true"
	}
	return "false"
}

func (n *BoolNode) Children() []Node {
	return []Node{}
}

func (n *BoolNode) Token() token.Token {
	return n.token
}

func (n *PrefixNode) String() string {
	return fmt.Sprintf("(%s %s)", n.token.Literal, n.expression)
}

func (n *PrefixNode) Children() []Node {
	return []Node{n.expression}
}

func (n *PrefixNode) Token() token.Token {
	return n.token
}

func (n *InfixNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.left, n.token.Literal, n.right)
}

func (n *InfixNode) Children() []Node {
	return []Node{n.left, n.right}
}

func (n *InfixNode) Token() token.Token {
	return n.token
}
