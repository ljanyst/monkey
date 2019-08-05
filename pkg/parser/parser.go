package parser

import (
	"fmt"
	"strconv"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/token"
)

const (
	LOWEST = iota
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

	infixParsers  map[token.TokenType]infixParseFn
	prefixParsers map[token.TokenType]prefixParseFn
	priorities    map[token.TokenType]int
}

func (p *Parser) nextToken() token.Token {
	tok := p.lexer.ReadToken()
	p.lexer.UnreadToken()
	return tok
}

func (p *Parser) mkErrWrongToken(expected string, got token.Token) error {
	lit := got.Literal
	if got.Type == token.EOF {
		lit = "end of input"
	}
	return fmt.Errorf("Parsing error: expected %s, got %s at %d:%d",
		expected, lit,
		got.Line, got.Column)
}

func (p *Parser) mkErrUnexpectedToken(got token.Token) error {
	return fmt.Errorf("Parsing error: don't know what to do with %q at %d:%d",
		got.Literal, got.Line, got.Column)
}

func (p *Parser) parseInt() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != token.INT {
		return nil, p.mkErrWrongToken("integer", tok)
	}

	i64, err := strconv.ParseInt(tok.Literal, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Parsing error: %s is not an integer literal", tok.Literal)
	}

	return &IntNode{tok, i64}, nil
}

func (p *Parser) parseString() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != token.STRING {
		return nil, p.mkErrWrongToken("string", tok)
	}
	return &StringNode{tok, tok.Literal}, nil
}

func (p *Parser) parseIdent() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != token.IDENT {
		return nil, p.mkErrWrongToken("identifier", tok)
	}
	return &IdentifierNode{tok, tok.Literal}, nil
}

func (p *Parser) parseBool() (Node, error) {
	tok := p.lexer.ReadToken()
	if tok.Type != token.TRUE && tok.Type != token.FALSE {
		return nil, p.mkErrWrongToken("boolean", tok)
	}
	if tok.Type == token.TRUE {
		return &BoolNode{tok, true}, nil
	}
	return &BoolNode{tok, false}, nil
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

func (p *Parser) getPriority(token token.Token) int {
	if prio, ok := p.priorities[token.Type]; ok {
		return prio
	}
	return LOWEST
}

func (p *Parser) parseExpression(priority int) (Node, error) {
	tok := p.nextToken()

	prefixParser, hasParser := p.prefixParsers[tok.Type]
	if !hasParser {
		return nil, p.mkErrUnexpectedToken(tok)
	}

	left, err := prefixParser()
	if err != nil {
		return nil, err
	}

	for {
		tok = p.nextToken()
		if priority >= p.getPriority(tok) || tok.Type == token.SEMICOLON {
			break
		}
		infixParser, hasParser := p.infixParsers[tok.Type]
		if !hasParser {
			return nil, p.mkErrUnexpectedToken(tok)
		}

		left, err = infixParser(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parsePrimaryExpression() (Node, error) {
	node, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	tok := p.lexer.ReadToken()
	if tok.Type != token.SEMICOLON {
		return nil, p.mkErrWrongToken("semicolon", tok)
	}
	return node, nil
}

func (p *Parser) Parse() (Node, error) {
	n := ProgramNode{}
	for {
		if p.nextToken().Type == token.EOF {
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

func NewParser(lexer *lexer.Lexer) *Parser {
	p := new(Parser)
	p.lexer = lexer

	p.prefixParsers = make(map[token.TokenType]prefixParseFn)
	p.prefixParsers[token.INT] = p.parseInt
	p.prefixParsers[token.STRING] = p.parseString
	p.prefixParsers[token.IDENT] = p.parseIdent
	p.prefixParsers[token.TRUE] = p.parseBool
	p.prefixParsers[token.FALSE] = p.parseBool
	p.prefixParsers[token.BANG] = p.parsePrefix
	p.prefixParsers[token.MINUS] = p.parsePrefix

	p.infixParsers = make(map[token.TokenType]infixParseFn)
	p.infixParsers[token.MINUS] = p.parseInfix
	p.infixParsers[token.PLUS] = p.parseInfix
	p.infixParsers[token.ASTERISK] = p.parseInfix
	p.infixParsers[token.SLASH] = p.parseInfix

	p.priorities = make(map[token.TokenType]int)
	p.priorities[token.MINUS] = SUM
	p.priorities[token.PLUS] = SUM
	p.priorities[token.ASTERISK] = PRODUCT
	p.priorities[token.SLASH] = PRODUCT
	return p
}
