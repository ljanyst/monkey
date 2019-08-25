package parser

import (
	"fmt"
	"strconv"

	"github.com/ljanyst/monkey/pkg/lexer"
)

const (
	LOWEST = iota
	ASSIGN
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() (Node, error)
	infixParseFn  func(Node) (Node, error)
)

type Parser struct {
	lexer *lexer.Lexer

	infixParsers  map[lexer.TokenType]infixParseFn
	prefixParsers map[lexer.TokenType]prefixParseFn
	priorities    map[lexer.TokenType]int
}

func (p *Parser) nextToken() lexer.Token {
	tok := p.lexer.ReadToken()
	p.lexer.UnreadToken()
	return tok
}

func mkErrWrongToken(expected string, got lexer.Token) error {
	lit := got.Literal
	if got.Type == lexer.EOF {
		lit = "end of input"
	}
	return fmt.Errorf("%s Parsing error: expected %s, got %q", got.Location(), expected, lit)
}

func mkErrUnexpectedToken(got lexer.Token) error {
	return fmt.Errorf("%s Parsing error: Unexpected token %q", got.Location(), got.Literal)
}

func (p *Parser) parseInt() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.INT {
		return nil, mkErrWrongToken("integer", tok)
	}

	i64, err := strconv.ParseInt(tok.Literal, 10, 64)
	if err != nil {
		return nil, mkErrWrongToken("integer literal", tok)
	}

	return &IntNode{tok, i64}, nil
}

func (p *Parser) parseString() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.STRING {
		return nil, mkErrWrongToken("string", tok)
	}
	return &StringNode{tok, tok.Literal}, nil
}

func (p *Parser) parseRune() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.RUNE {
		return nil, mkErrWrongToken("rune", tok)
	}
	return &RuneNode{tok, []rune(tok.Literal)[0]}, nil
}

func (p *Parser) parseIdent() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.IDENT {
		return nil, mkErrWrongToken("identifier", tok)
	}
	return &IdentifierNode{tok, tok.Literal}, nil
}

func (p *Parser) parseBool() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.TRUE && tok.Type != lexer.FALSE {
		return nil, mkErrWrongToken("boolean", tok)
	}
	if tok.Type == lexer.TRUE {
		return &BoolNode{tok, true}, nil
	}
	return &BoolNode{tok, false}, nil
}

func (p *Parser) parseNil() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != lexer.NIL {
		return nil, mkErrWrongToken("nil", tok)
	}
	return &NilNode{tok}, nil
}

func (p *Parser) parsePrefix() (Node, error) {
	tok := p.lexer.ReadToken()
	exp, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	return &PrefixNode{tok, exp}, nil
}

func (p *Parser) parseInfix(left Node) (Node, error) {
	tok := p.lexer.ReadToken()
	right, err := p.parseExpression(p.getPriority(tok))
	if err != nil {
		return nil, err
	}
	return &InfixNode{tok, left, right}, nil
}

func (p *Parser) parseParen() (Node, error) {
	tok := p.lexer.ReadToken()
	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	tok = p.lexer.ReadToken()
	if tok.Type != lexer.RPAREN {
		return nil, mkErrWrongToken(")", tok)
	}
	return exp, nil
}

func (p *Parser) parseBlock() (Node, error) {
	n := BlockNode{}

	tok := p.lexer.ReadToken()
	if tok.Type != lexer.LBRACE {
		return nil, mkErrWrongToken("{", tok)
	}

	for {
		if p.nextToken().Type == lexer.RBRACE {
			break
		}

		node, err := p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
		n.children = append(n.children, node)
	}

	tok = p.lexer.ReadToken()
	if tok.Type != lexer.RBRACE {
		return nil, mkErrWrongToken("}", tok)
	}

	return &n, nil
}

func (p *Parser) parseConditional() (Node, error) {
	ifTok := p.lexer.ReadToken()
	tok := p.lexer.ReadToken()

	if tok.Type != lexer.LPAREN {
		return nil, mkErrWrongToken("(", tok)
	}

	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	tok = p.lexer.ReadToken()
	if tok.Type != lexer.RPAREN {
		return nil, mkErrWrongToken(")", tok)
	}

	consequent, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	var alternative Node
	if p.nextToken().Type == lexer.ELSE {
		tok = p.lexer.ReadToken()
		alternative, err = p.parseBlock()
	}

	exp := &ConditionalNode{ifTok, condition, consequent, alternative}

	return exp, nil
}

func (p *Parser) parseAssign(left Node) (Node, error) {
	if left.Token().Type != lexer.IDENT {
		return nil, mkErrWrongToken("identifier", left.Token())
	}

	tok := p.lexer.ReadToken()

	right, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	return &InfixNode{tok, left, right}, nil
}

func (p *Parser) parseStatement() (Node, error) {
	tok := p.lexer.ReadToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	return &StatementNode{tok, exp}, nil
}

func (p *Parser) parseFunction() (Node, error) {
	fnTok := p.lexer.ReadToken()

	tok := p.lexer.ReadToken()
	if tok.Type != lexer.LPAREN {
		return nil, mkErrWrongToken("(", tok)
	}

	params := []Node{}

	for {
		tok = p.nextToken()
		if tok.Type == lexer.RPAREN {
			p.lexer.ReadToken()
			break
		}

		ident, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		params = append(params, ident)

		tok = p.nextToken()
		if tok.Type != lexer.COMMA && tok.Type != lexer.RPAREN {
			return nil, mkErrWrongToken(", or )", tok)
		}
		if tok.Type == lexer.COMMA {
			p.lexer.ReadToken()
		}
	}

	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &FunctionNode{fnTok, params, body}, nil
}

func (p *Parser) parseFunctionCall(left Node) (Node, error) {
	parenTok := p.lexer.ReadToken()

	args := []Node{}

	for {
		tok := p.nextToken()
		if tok.Type == lexer.RPAREN {
			p.lexer.ReadToken()
			break
		}

		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		args = append(args, exp)

		tok = p.nextToken()
		if tok.Type == lexer.COMMA {
			p.lexer.ReadToken()
		}
	}

	return &FunctionCallNode{parenTok, left, args}, nil
}

func (p *Parser) getPriority(token lexer.Token) int {
	if prio, ok := p.priorities[token.Type]; ok {
		return prio
	}
	return LOWEST
}

func isTerminator(tok lexer.Token) bool {
	switch tok.Type {
	case lexer.RPAREN:
		return true
	default:
		return false
	}
}

func (p *Parser) parseExpression(priority int) (Node, error) {
	tok := p.nextToken()

	prefixParser, hasParser := p.prefixParsers[tok.Type]
	if !hasParser {
		return nil, mkErrUnexpectedToken(tok)
	}

	left, err := prefixParser()
	if err != nil {
		return nil, err
	}

	for {
		tok = p.nextToken()
		if priority >= p.getPriority(tok) || tok.Type == lexer.SEMICOLON {
			break
		}
		infixParser, hasParser := p.infixParsers[tok.Type]
		if !hasParser {
			if isTerminator(tok) {
				p.lexer.ReadToken()
				return left, nil
			}
			return nil, mkErrUnexpectedToken(tok)
		}

		left, err = infixParser(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpression() (Node, error) {
	var node Node
	var err error
	switch p.nextToken().Type {
	case lexer.LET, lexer.RETURN:
		node, err = p.parseStatement()
	default:
		node, err = p.parseExpression(LOWEST)
	}

	if err != nil {
		return nil, err
	}

	tok := p.lexer.ReadToken()
	if tok.Type != lexer.SEMICOLON {
		return nil, mkErrWrongToken("semicolon", tok)
	}
	return node, nil
}

func (p *Parser) Parse() (Node, error) {
	n := BlockNode{true, []Node{}}
	for {
		if p.nextToken().Type == lexer.EOF {
			break
		}

		node, err := p.parsePrimaryExpression()
		if err != nil {
			return nil, err
		}
		n.children = append(n.children, node)
	}
	return &n, nil
}

func NewParser(lex *lexer.Lexer) *Parser {
	p := new(Parser)
	p.lexer = lex

	p.prefixParsers = make(map[lexer.TokenType]prefixParseFn)
	p.prefixParsers[lexer.INT] = p.parseInt
	p.prefixParsers[lexer.STRING] = p.parseString
	p.prefixParsers[lexer.RUNE] = p.parseRune
	p.prefixParsers[lexer.IDENT] = p.parseIdent
	p.prefixParsers[lexer.TRUE] = p.parseBool
	p.prefixParsers[lexer.FALSE] = p.parseBool
	p.prefixParsers[lexer.NIL] = p.parseNil
	p.prefixParsers[lexer.BANG] = p.parsePrefix
	p.prefixParsers[lexer.MINUS] = p.parsePrefix
	p.prefixParsers[lexer.LPAREN] = p.parseParen
	p.prefixParsers[lexer.IF] = p.parseConditional
	p.prefixParsers[lexer.FUNCTION] = p.parseFunction

	p.infixParsers = make(map[lexer.TokenType]infixParseFn)
	for _, t := range []lexer.TokenType{
		lexer.MINUS, lexer.PLUS, lexer.ASTERISK, lexer.SLASH, lexer.EQ,
		lexer.NOT_EQ, lexer.LT, lexer.LE, lexer.GT, lexer.GE,
	} {
		p.infixParsers[t] = p.parseInfix
	}
	p.infixParsers[lexer.ASSIGN] = p.parseAssign
	p.infixParsers[lexer.LPAREN] = p.parseFunctionCall

	p.priorities = make(map[lexer.TokenType]int)
	p.priorities[lexer.MINUS] = SUM
	p.priorities[lexer.PLUS] = SUM
	p.priorities[lexer.ASTERISK] = PRODUCT
	p.priorities[lexer.SLASH] = PRODUCT
	p.priorities[lexer.EQ] = COMPARISON
	p.priorities[lexer.NOT_EQ] = COMPARISON
	p.priorities[lexer.LT] = COMPARISON
	p.priorities[lexer.LE] = COMPARISON
	p.priorities[lexer.GT] = COMPARISON
	p.priorities[lexer.GE] = COMPARISON
	p.priorities[lexer.ASSIGN] = ASSIGN
	p.priorities[lexer.LPAREN] = CALL
	return p
}
