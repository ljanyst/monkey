package parser

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/token"
)

type Node interface {
	String(string) string
	Children() []Node
	Token() token.Token
}

type BlockNode struct {
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

type ConditionalNode struct {
	token       token.Token
	condition   Node
	consequent  Node
	alternative Node
}

type StatementNode struct {
	token      token.Token
	expression Node
}

type FunctionNode struct {
	token  token.Token
	params []Node
	body   Node
}

type FunctionCallNode struct {
	token token.Token
	name  Node
	args  []Node
}

func (n *BlockNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString(padding)
	sb.WriteString("{\n")
	for _, node := range n.children {
		sb.WriteString(padding)
		sb.WriteString("  ")
		sb.WriteString(node.String(padding + "  "))
		sb.WriteString("\n")
	}
	sb.WriteString(padding)
	sb.WriteString("}")
	return sb.String()
}

func (n *BlockNode) Children() []Node {
	return n.children
}

func (n *BlockNode) Token() token.Token {
	return token.Token{"", "BLOCK", 0, 0}
}

func (n *IntNode) String(padding string) string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *IntNode) Children() []Node {
	return []Node{}
}

func (n *IntNode) Token() token.Token {
	return n.token
}

func (n *StringNode) String(padding string) string {
	return fmt.Sprintf("%q", n.Value)
}

func (n *StringNode) Children() []Node {
	return []Node{}
}

func (n *StringNode) Token() token.Token {
	return n.token
}

func (n *IdentifierNode) String(padding string) string {
	return n.Value
}

func (n *IdentifierNode) Children() []Node {
	return []Node{}
}

func (n *IdentifierNode) Token() token.Token {
	return n.token
}

func (n *BoolNode) String(padding string) string {
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

func (n *PrefixNode) String(padding string) string {
	return fmt.Sprintf("(%s %s)", n.token.Literal, n.expression.String(padding))
}

func (n *PrefixNode) Children() []Node {
	return []Node{n.expression}
}

func (n *PrefixNode) Token() token.Token {
	return n.token
}

func (n *InfixNode) String(padding string) string {
	return fmt.Sprintf("(%s %s %s)", n.left.String(padding), n.token.Literal, n.right.String(padding))
}

func (n *InfixNode) Children() []Node {
	return []Node{n.left, n.right}
}

func (n *InfixNode) Token() token.Token {
	return n.token
}

func (n *ConditionalNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("if %s\n", n.condition.String(padding)))
	sb.WriteString(n.consequent.String(padding))
	if n.alternative != nil {
		sb.WriteString(fmt.Sprintf("\n%selse\n", padding))
		sb.WriteString(n.alternative.String(padding))
	}
	return sb.String()
}

func (n *ConditionalNode) Children() []Node {
	return []Node{n.condition, n.consequent, n.alternative}
}

func (n *ConditionalNode) Token() token.Token {
	return n.token
}

func (n *StatementNode) String(padding string) string {
	return fmt.Sprintf("%s %s", n.token.Literal, n.expression.String(padding))
}

func (n *StatementNode) Children() []Node {
	return []Node{n.expression}
}

func (n *StatementNode) Token() token.Token {
	return n.token
}

func (n *FunctionNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString("fn(")
	for i, param := range n.params {
		sb.WriteString(param.String(padding))
		if i < len(n.params)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")\n")
	sb.WriteString(n.body.String(padding))
	return sb.String()
}

func (n *FunctionNode) Children() []Node {
	return append(n.params, n.body)
}

func (n *FunctionNode) Token() token.Token {
	return n.token
}

func (n *FunctionCallNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s(", n.name.String(padding)))
	for i, arg := range n.args {
		sb.WriteString(arg.String(padding))
		if i < len(n.args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (n *FunctionCallNode) Children() []Node {
	return append([]Node{n.name}, n.args...)
}

func (n *FunctionCallNode) Token() token.Token {
	return n.token
}
