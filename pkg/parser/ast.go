package parser

import (
	"fmt"
	"strings"

	"github.com/ljanyst/monkey/pkg/lexer"
)

type Node interface {
	String(string) string
	Children() []Node
	Token() lexer.Token
}

type BlockNode struct {
	Implicit bool
	children []Node
}

type IntNode struct {
	token lexer.Token
	Value int64
}

type StringNode struct {
	token lexer.Token
	Value string
}

type IdentifierNode struct {
	token lexer.Token
	Value string
}

type BoolNode struct {
	token lexer.Token
	Value bool
}

type NilNode struct {
	token lexer.Token
}

type PrefixNode struct {
	token      lexer.Token
	Expression Node
}

type InfixNode struct {
	token lexer.Token
	Left  Node
	Right Node
}

type ConditionalNode struct {
	token       lexer.Token
	Condition   Node
	Consequent  Node
	Alternative Node
}

type StatementNode struct {
	token      lexer.Token
	expression Node
}

type FunctionNode struct {
	token  lexer.Token
	Params []Node
	Body   Node
}

type FunctionCallNode struct {
	token lexer.Token
	Name  Node
	Args  []Node
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

func (n *BlockNode) Token() lexer.Token {
	return lexer.Token{lexer.BLOCK, "BLOCK", 0, 0, nil}
}

func (n *IntNode) String(padding string) string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *IntNode) Children() []Node {
	return []Node{}
}

func (n *IntNode) Token() lexer.Token {
	return n.token
}

func (n *StringNode) String(padding string) string {
	return fmt.Sprintf("%q", n.Value)
}

func (n *StringNode) Children() []Node {
	return []Node{}
}

func (n *StringNode) Token() lexer.Token {
	return n.token
}

func (n *IdentifierNode) String(padding string) string {
	return n.Value
}

func (n *IdentifierNode) Children() []Node {
	return []Node{}
}

func (n *IdentifierNode) Token() lexer.Token {
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

func (n *BoolNode) Token() lexer.Token {
	return n.token
}

func (n *PrefixNode) String(padding string) string {
	return fmt.Sprintf("(%s %s)", n.token.Literal, n.Expression.String(padding))
}

func (n *PrefixNode) Children() []Node {
	return []Node{n.Expression}
}

func (n *PrefixNode) Token() lexer.Token {
	return n.token
}

func (n *InfixNode) String(padding string) string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(padding), n.token.Literal, n.Right.String(padding))
}

func (n *InfixNode) Children() []Node {
	return []Node{n.Left, n.Right}
}

func (n *InfixNode) Token() lexer.Token {
	return n.token
}

func (n *ConditionalNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("if %s\n", n.Condition.String(padding)))
	sb.WriteString(n.Consequent.String(padding))
	if n.Alternative != nil {
		sb.WriteString(fmt.Sprintf("\n%selse\n", padding))
		sb.WriteString(n.Alternative.String(padding))
	}
	return sb.String()
}

func (n *ConditionalNode) Children() []Node {
	return []Node{n.Condition, n.Consequent, n.Alternative}
}

func (n *ConditionalNode) Token() lexer.Token {
	return n.token
}

func (n *StatementNode) String(padding string) string {
	return fmt.Sprintf("%s %s", n.token.Literal, n.expression.String(padding))
}

func (n *StatementNode) Children() []Node {
	return []Node{n.expression}
}

func (n *StatementNode) Token() lexer.Token {
	return n.token
}

func (n *FunctionNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString("fn(")
	for i, param := range n.Params {
		sb.WriteString(param.String(padding))
		if i < len(n.Params)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")\n")
	sb.WriteString(n.Body.String(padding))
	return sb.String()
}

func (n *FunctionNode) Children() []Node {
	return append(n.Params, n.Body)
}

func (n *FunctionNode) Token() lexer.Token {
	return n.token
}

func (n *FunctionCallNode) String(padding string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s(", n.Name.String(padding)))
	for i, arg := range n.Args {
		sb.WriteString(arg.String(padding))
		if i < len(n.Args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

func (n *FunctionCallNode) Children() []Node {
	return append([]Node{n.Name}, n.Args...)
}

func (n *FunctionCallNode) Token() lexer.Token {
	return n.token
}

func (n *NilNode) String(padding string) string {
	return "nil"
}

func (n *NilNode) Children() []Node {
	return []Node{}
}

func (n *NilNode) Token() lexer.Token {
	return n.token
}
