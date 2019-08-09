package lexer

type TokenType string

const (
	LET       = "LET"
	IDENT     = "IDENT"
	ASSIGN    = "ASSIGN"
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
	LE        = "LE"
	GT        = "GT"
	GE        = "GE"
	IF        = "IF"
	RETURN    = "RETURN"
	TRUE      = "TRUE"
	ELSE      = "ELSE"
	FALSE     = "FALSE"
	STRING    = "STRING"
	EQ        = "EQ"
	NOT_EQ    = "NOT_EQ"
	INVALID   = "INVALID"
	EOF       = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    uint32
	Column  uint32
}

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupKeyword(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
