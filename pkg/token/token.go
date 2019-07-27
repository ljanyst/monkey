package token

type TokenType string

const (
	LET       = "LET"
	IDENT     = "IDENT"
	ASSIGN    = "="
	INT       = "INT"
	SEMICOLON = "SEMICOLON"
	FUNCTION  = "FUNCTION"
	LPAREN    = "LPAREN"
	COMMA     = "COMMA"
	RPAREN    = "RPAREN"
	LBRACE    = "LBRACE"
	PLUS      = "PLUS"
	RBRACE    = "RBRACE"
	BANG      = "BANG"
	MINUS     = "MINUS"
	SLASH     = "SLASH"
	ASTERISK  = "ASTERISK"
	LT        = "LT"
	GT        = "GT"
	IF        = "IF"
	RETURN    = "RETURN"
	TRUE      = "TRUE"
	ELSE      = "ELSE"
	FALSE     = "FALSE"
	STRING    = "STRING"
	EQ        = "EQ"
	NOT_EQ    = "NOT_EQ"
	EOF       = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    uint32
	Column  uint32
}
