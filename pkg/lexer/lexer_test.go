package lexer

import (
	"testing"

	"github.com/ljanyst/monkey/pkg/token"
)

func compareTokens(t *testing.T, got, expected token.Token) bool {
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

	tests := []token.Token{
		{token.LET, "let", 0, 0},
		{token.IDENT, "five", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.INT, "5", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.LET, "let", 0, 0},
		{token.IDENT, "ten", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.INT, "10", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.LET, "let", 0, 0},
		{token.IDENT, "add", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.FUNCTION, "fn", 0, 0},
		{token.LPAREN, "(", 0, 0},
		{token.IDENT, "x", 0, 0},
		{token.COMMA, ",", 0, 0},
		{token.IDENT, "y", 0, 0},
		{token.RPAREN, ")", 0, 0},
		{token.LBRACE, "{", 0, 0},
		{token.IDENT, "x", 0, 0},
		{token.PLUS, "+", 0, 0},
		{token.IDENT, "y", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.RBRACE, "}", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.LET, "let", 0, 0},
		{token.IDENT, "result", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.IDENT, "add", 0, 0},
		{token.LPAREN, "(", 0, 0},
		{token.IDENT, "five", 0, 0},
		{token.COMMA, ",", 0, 0},
		{token.IDENT, "ten", 0, 0},
		{token.RPAREN, ")", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.BANG, "!", 0, 0},
		{token.MINUS, "-", 0, 0},
		{token.SLASH, "/", 0, 0},
		{token.ASTERISK, "*", 0, 0},
		{token.INT, "5", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.INT, "5", 0, 0},
		{token.LT, "<", 0, 0},
		{token.INT, "10", 0, 0},
		{token.GT, ">", 0, 0},
		{token.INT, "5", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.IF, "if", 0, 0},
		{token.LPAREN, "(", 0, 0},
		{token.INT, "5", 0, 0},
		{token.LT, "<", 0, 0},
		{token.INT, "10", 0, 0},
		{token.RPAREN, ")", 0, 0},
		{token.LBRACE, "{", 0, 0},
		{token.RETURN, "return", 0, 0},
		{token.TRUE, "true", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.RBRACE, "}", 0, 0},
		{token.ELSE, "else", 0, 0},
		{token.LBRACE, "{", 0, 0},
		{token.RETURN, "return", 0, 0},
		{token.FALSE, "false", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.RBRACE, "}", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.LET, "let", 0, 0},
		{token.IDENT, "żółwik", 0, 0},
		{token.ASSIGN, "=", 0, 0},
		{token.STRING, "zażółć gęślą jaźń", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.INT, "12", 0, 0},
		{token.LE, "<=", 0, 0},
		{token.INT, "46", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.INT, "43", 0, 0},
		{token.GE, ">=", 0, 0},
		{token.INT, "17", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.INT, "10", 0, 0},
		{token.EQ, "==", 0, 0},
		{token.INT, "10", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.INT, "10", 0, 0},
		{token.NOT_EQ, "!=", 0, 0},
		{token.INT, "9", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
		{token.EOF, "", 0, 0},
	}

	l := NewLexerFromString(input)

	for _, expected := range tests {
		got := l.ReadToken()
		compareTokens(t, got, expected)
	}
}

func TestUnreadToken(t *testing.T) {
	input := `!-/*5;`

	tests := []token.Token{
		{token.BANG, "!", 0, 0},
		{token.MINUS, "-", 0, 0},
		{token.SLASH, "/", 0, 0},
		{token.ASTERISK, "*", 0, 0},
		{token.INT, "5", 0, 0},
		{token.SEMICOLON, ";", 0, 0},
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
