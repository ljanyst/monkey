package lexer

import (
	"testing"
)

func compareTokens(t *testing.T, got, expected Token) bool {
	if got.Type != expected.Type || got.Literal != expected.Literal {
		t.Errorf("Wrong token: expected %s(%q), got %s(%q), at %d:%d",
			expected.Type, expected.Literal, got.Type, got.Literal, got.Line, got.Column)
		return false
	}
	return true
}

func TestReadToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
};

let żółwik = "zażółć gęślą jaźń";

12 <= 46;
43 >= 17;
10 == 10;
10 != 9;
`

	tests := []Token{
		{LET, "let", 0, 0, nil},
		{IDENT, "five", 0, 0, nil},
		{ASSIGN, "=", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{LET, "let", 0, 0, nil},
		{IDENT, "ten", 0, 0, nil},
		{ASSIGN, "=", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{LET, "let", 0, 0, nil},
		{IDENT, "add", 0, 0, nil},
		{ASSIGN, "=", 0, 0, nil},
		{FUNCTION, "fn", 0, 0, nil},
		{LPAREN, "(", 0, 0, nil},
		{IDENT, "x", 0, 0, nil},
		{COMMA, ",", 0, 0, nil},
		{IDENT, "y", 0, 0, nil},
		{RPAREN, ")", 0, 0, nil},
		{LBRACE, "{", 0, 0, nil},
		{IDENT, "x", 0, 0, nil},
		{PLUS, "+", 0, 0, nil},
		{IDENT, "y", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{RBRACE, "}", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{LET, "let", 0, 0, nil},
		{IDENT, "result", 0, 0, nil},
		{ASSIGN, "=", 0, 0, nil},
		{IDENT, "add", 0, 0, nil},
		{LPAREN, "(", 0, 0, nil},
		{IDENT, "five", 0, 0, nil},
		{COMMA, ",", 0, 0, nil},
		{IDENT, "ten", 0, 0, nil},
		{RPAREN, ")", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{BANG, "!", 0, 0, nil},
		{MINUS, "-", 0, 0, nil},
		{SLASH, "/", 0, 0, nil},
		{ASTERISK, "*", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{LT, "<", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{GT, ">", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{IF, "if", 0, 0, nil},
		{LPAREN, "(", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{LT, "<", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{RPAREN, ")", 0, 0, nil},
		{LBRACE, "{", 0, 0, nil},
		{RETURN, "return", 0, 0, nil},
		{TRUE, "true", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{RBRACE, "}", 0, 0, nil},
		{ELSE, "else", 0, 0, nil},
		{LBRACE, "{", 0, 0, nil},
		{RETURN, "return", 0, 0, nil},
		{FALSE, "false", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{RBRACE, "}", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{LET, "let", 0, 0, nil},
		{IDENT, "żółwik", 0, 0, nil},
		{ASSIGN, "=", 0, 0, nil},
		{STRING, "zażółć gęślą jaźń", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{INT, "12", 0, 0, nil},
		{LE, "<=", 0, 0, nil},
		{INT, "46", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{INT, "43", 0, 0, nil},
		{GE, ">=", 0, 0, nil},
		{INT, "17", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{EQ, "==", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{INT, "10", 0, 0, nil},
		{NOT_EQ, "!=", 0, 0, nil},
		{INT, "9", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
		{EOF, "", 0, 0, nil},
	}

	l := NewLexerFromString(input, "input")

	for _, expected := range tests {
		got := l.ReadToken()
		compareTokens(t, got, expected)
	}
}

func TestUnreadToken(t *testing.T) {
	input := `!-/*5;`

	tests := []Token{
		{BANG, "!", 0, 0, nil},
		{MINUS, "-", 0, 0, nil},
		{SLASH, "/", 0, 0, nil},
		{ASTERISK, "*", 0, 0, nil},
		{INT, "5", 0, 0, nil},
		{SEMICOLON, ";", 0, 0, nil},
	}

	l := NewLexerFromString(input, "input")

	for _, expected := range tests {
		got := l.ReadToken()
		compareTokens(t, got, expected)
		l.UnreadToken()
		got = l.ReadToken()
		compareTokens(t, got, expected)
	}
}
