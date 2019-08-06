package parser

import (
	"fmt"
	"testing"

	"github.com/ljanyst/monkey/pkg/lexer"
	"github.com/ljanyst/monkey/pkg/token"
)

func compareAst(t *testing.T, got, expected Node, print bool, depth string) bool {
	gt := got.Token()
	et := expected.Token()

	if print {
		fmt.Printf("%s %s\n", depth, gt.Literal)
	}

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
		if !compareAst(t, gotChildren[i], expectedChildren[i], print, depth+"  ") {
			return false
		}
	}
	return true
}

func parseAndCompareAst(t *testing.T, input string, expected Node) bool {
	l := lexer.NewLexerFromString(input)
	p := NewParser(l)
	parsed, err := p.Parse()
	if err != nil {
		t.Errorf("Parser error: %s", err)
		return false
	}
	if !compareAst(t, parsed, expected, false, "") {
		t.Errorf("ASTs differ")
		return false
	}
	return true
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

func TestInfixPriority(t *testing.T) {
	input := `10 + 2;
3 * 20;
10 + 2 * 6;
12 * 7 + 12;
12 * 7 + 12 * 8;
2 + 4 * 5 * 6 * 7;
-12 * 7 + 12 * -8;
-12 * 7 == 12 + -8;
`
	expected := ProgramNode{
		[]Node{
			&InfixNode{
				token.Token{token.PLUS, "+", 1, 4},
				&IntNode{
					token.Token{token.INT, "10", 1, 1},
					10,
				},
				&IntNode{
					token.Token{token.INT, "2", 1, 6},
					2,
				},
			},
			&InfixNode{
				token.Token{token.ASTERISK, "*", 2, 3},
				&IntNode{
					token.Token{token.INT, "3", 2, 1},
					3,
				},
				&IntNode{
					token.Token{token.INT, "20", 2, 5},
					20,
				},
			},
			&InfixNode{
				token.Token{token.PLUS, "+", 3, 4},
				&IntNode{
					token.Token{token.INT, "10", 3, 1},
					10,
				},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 3, 8},
					&IntNode{
						token.Token{token.INT, "2", 3, 6},
						2,
					},
					&IntNode{
						token.Token{token.INT, "6", 3, 10},
						6,
					},
				},
			},
			&InfixNode{
				token.Token{token.PLUS, "+", 4, 8},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 4, 4},
					&IntNode{
						token.Token{token.INT, "12", 4, 1},
						12,
					},
					&IntNode{
						token.Token{token.INT, "7", 4, 8},
						7,
					},
				},
				&IntNode{
					token.Token{token.INT, "12", 4, 10},
					12,
				},
			},
			&InfixNode{
				token.Token{token.PLUS, "+", 5, 8},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 5, 4},
					&IntNode{
						token.Token{token.INT, "12", 5, 1},
						12,
					},
					&IntNode{
						token.Token{token.INT, "7", 5, 8},
						7,
					},
				},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 5, 13},
					&IntNode{
						token.Token{token.INT, "12", 5, 10},
						12,
					},
					&IntNode{
						token.Token{token.INT, "8", 5, 15},
						8,
					},
				},
			},
			&InfixNode{
				token.Token{token.PLUS, "+", 6, 3},
				&IntNode{
					token.Token{token.INT, "2", 6, 1},
					2,
				},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 6, 7},
					&IntNode{
						token.Token{token.INT, "4", 6, 5},
						4,
					},
					&InfixNode{
						token.Token{token.ASTERISK, "*", 6, 11},
						&IntNode{
							token.Token{token.INT, "5", 6, 9},
							5,
						},
						&InfixNode{
							token.Token{token.ASTERISK, "*", 6, 15},
							&IntNode{
								token.Token{token.INT, "6", 6, 13},
								6,
							},
							&IntNode{
								token.Token{token.INT, "7", 6, 17},
								7,
							},
						},
					},
				},
			},
			&InfixNode{
				token.Token{token.PLUS, "+", 7, 9},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 7, 5},
					&PrefixNode{
						token.Token{token.MINUS, "-", 7, 1},
						&IntNode{
							token.Token{token.INT, "12", 7, 2},
							12,
						},
					},
					&IntNode{
						token.Token{token.INT, "7", 7, 9},
						7,
					},
				},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 7, 14},
					&IntNode{
						token.Token{token.INT, "12", 7, 11},
						12,
					},
					&PrefixNode{
						token.Token{token.MINUS, "-", 7, 16},
						&IntNode{
							token.Token{token.INT, "8", 7, 17},
							8,
						},
					},
				},
			},
			&InfixNode{
				token.Token{token.EQ, "==", 8, 9},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 8, 5},
					&PrefixNode{
						token.Token{token.MINUS, "-", 8, 1},
						&IntNode{
							token.Token{token.INT, "12", 8, 2},
							12,
						},
					},
					&IntNode{
						token.Token{token.INT, "7", 8, 9},
						7,
					},
				},
				&InfixNode{
					token.Token{token.ASTERISK, "*", 8, 15},
					&IntNode{
						token.Token{token.INT, "12", 8, 12},
						12,
					},
					&PrefixNode{
						token.Token{token.MINUS, "-", 8, 17},
						&IntNode{
							token.Token{token.INT, "8", 8, 18},
							8,
						},
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}
