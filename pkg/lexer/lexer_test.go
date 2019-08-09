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
		{LET, "let", 0, 0},
		{IDENT, "five", 0, 0},
		{ASSIGN, "=", 0, 0},
		{INT, "5", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{LET, "let", 0, 0},
		{IDENT, "ten", 0, 0},
		{ASSIGN, "=", 0, 0},
		{INT, "10", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{LET, "let", 0, 0},
		{IDENT, "add", 0, 0},
		{ASSIGN, "=", 0, 0},
		{FUNCTION, "fn", 0, 0},
		{LPAREN, "(", 0, 0},
		{IDENT, "x", 0, 0},
		{COMMA, ",", 0, 0},
		{IDENT, "y", 0, 0},
		{RPAREN, ")", 0, 0},
		{LBRACE, "{", 0, 0},
		{IDENT, "x", 0, 0},
		{PLUS, "+", 0, 0},
		{IDENT, "y", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{RBRACE, "}", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{LET, "let", 0, 0},
		{IDENT, "result", 0, 0},
		{ASSIGN, "=", 0, 0},
		{IDENT, "add", 0, 0},
		{LPAREN, "(", 0, 0},
		{IDENT, "five", 0, 0},
		{COMMA, ",", 0, 0},
		{IDENT, "ten", 0, 0},
		{RPAREN, ")", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{BANG, "!", 0, 0},
		{MINUS, "-", 0, 0},
		{SLASH, "/", 0, 0},
		{ASTERISK, "*", 0, 0},
		{INT, "5", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{INT, "5", 0, 0},
		{LT, "<", 0, 0},
		{INT, "10", 0, 0},
		{GT, ">", 0, 0},
		{INT, "5", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{IF, "if", 0, 0},
		{LPAREN, "(", 0, 0},
		{INT, "5", 0, 0},
		{LT, "<", 0, 0},
		{INT, "10", 0, 0},
		{RPAREN, ")", 0, 0},
		{LBRACE, "{", 0, 0},
		{RETURN, "return", 0, 0},
		{TRUE, "true", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{RBRACE, "}", 0, 0},
		{ELSE, "else", 0, 0},
		{LBRACE, "{", 0, 0},
		{RETURN, "return", 0, 0},
		{FALSE, "false", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{RBRACE, "}", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{LET, "let", 0, 0},
		{IDENT, "żółwik", 0, 0},
		{ASSIGN, "=", 0, 0},
		{STRING, "zażółć gęślą jaźń", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{INT, "12", 0, 0},
		{LE, "<=", 0, 0},
		{INT, "46", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{INT, "43", 0, 0},
		{GE, ">=", 0, 0},
		{INT, "17", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{INT, "10", 0, 0},
		{EQ, "==", 0, 0},
		{INT, "10", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{INT, "10", 0, 0},
		{NOT_EQ, "!=", 0, 0},
		{INT, "9", 0, 0},
		{SEMICOLON, ";", 0, 0},
		{EOF, "", 0, 0},
	}

	l := NewLexerFromString(input)

	for _, expected := range tests {
		got := l.ReadToken()
		compareTokens(t, got, expected)
	}
}

func TestUnreadToken(t *testing.T) {
	input := `!-/*5;`

	tests := []Token{
		{BANG, "!", 0, 0},
		{MINUS, "-", 0, 0},
		{SLASH, "/", 0, 0},
		{ASTERISK, "*", 0, 0},
		{INT, "5", 0, 0},
		{SEMICOLON, ";", 0, 0},
	}

	l := NewLexerFromString(input)

	for _, expected := range tests {
		got := l.ReadToken()
		compareTokens(t, got, expected)
		l.UnreadToken()
		got = l.ReadToken()
		compareTokens(t, got, expected)
	}
}
