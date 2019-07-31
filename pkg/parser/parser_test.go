package parser

import (
	"testing"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/token"
)

func compareAst(t *testing.T, got, expected Node) bool {
	gt := got.Token()
	et := expected.Token()
	if gt.Type != et.Type || gt.Literal != et.Literal {
		t.Errorf("Wrong token: expected %s(%q), got %s(%q), at %d:%d",
			et.Type, et.Literal, gt.Type, gt.Literal, gt.Line, gt.Column)
		return false
	}

	gotChildren := got.Children()
	expectedChildren := expected.Children()

	if len(gotChildren) != len(expectedChildren) {
		t.Errorf("Wrong number of children: expected %d, got %d, for token %s(%q), at %d:%d",
			len(expectedChildren), len(gotChildren), gt.Type, gt.Literal, gt.Line, gt.Column)
		return false
	}

	for i := 0; i < len(gotChildren); i++ {
		if !compareAst(t, gotChildren[i], expectedChildren[i]) {
			return false
		}
	}
	return false
}

func parseAndCompareAst(t *testing.T, input string, expected Node) bool {
	l := lexer.NewLexerFromString(input)
	p := NewParser(l)
	parsed, err := p.Parse()
	if err != nil {
		t.Errorf("Parser error: %s", err)
		return false
	}
	return compareAst(t, parsed, expected)
}

func TestLiteralsAndIdentifiers(t *testing.T) {
	input := `10;
"zażółć gęślą jaźń";
test;
true;
false;
!true;
-10;
`

	expected := ProgramNode{
		[]Node{
			&IntNode{
				token.Token{token.INT, "10", 1, 1},
				10,
			},
			&StringNode{
				token.Token{token.STRING, "zażółć gęślą jaźń", 2, 1},
				"zażółć gęślą jaźń",
			},
			&IdentifierNode{
				token.Token{token.IDENT, "test", 3, 1},
				"test",
			},
			&BoolNode{
				token.Token{token.TRUE, "true", 4, 1},
				true,
			},
			&BoolNode{
				token.Token{token.FALSE, "false", 5, 1},
				false,
			},
			&PrefixNode{
				token.Token{token.BANG, "!", 6, 1},
				&BoolNode{
					token.Token{token.TRUE, "true", 6, 2},
					true,
				},
			},
			&PrefixNode{
				token.Token{token.MINUS, "-", 7, 1},
				&IntNode{
					token.Token{token.INT, "10", 7, 2},
					10,
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}
