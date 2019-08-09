package parser

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ljanyst/monkey/pkg/lexer"
)

var printAst = flag.Bool("print-ast", false, "print the AST")
var printProgram = flag.Bool("print-program", false, "print the parsed program")

func compareAst(t *testing.T, got, expected Node, print bool, depth string) bool {
	if got == nil || expected == nil {
		if expected == got {
			return true
		}
		return false
	}

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

	if *printProgram {
		fmt.Printf("%s\n", parsed.String(""))
	}

	if !compareAst(t, parsed, expected, *printAst, "") {
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

	expected := BlockNode{
		[]Node{
			&IntNode{
				lexer.Token{lexer.INT, "10", 1, 1},
				10,
			},
			&StringNode{
				lexer.Token{lexer.STRING, "zażółć gęślą jaźń", 2, 1},
				"zażółć gęślą jaźń",
			},
			&IdentifierNode{
				lexer.Token{lexer.IDENT, "test", 3, 1},
				"test",
			},
			&BoolNode{
				lexer.Token{lexer.TRUE, "true", 4, 1},
				true,
			},
			&BoolNode{
				lexer.Token{lexer.FALSE, "false", 5, 1},
				false,
			},
			&PrefixNode{
				lexer.Token{lexer.BANG, "!", 6, 1},
				&BoolNode{
					lexer.Token{lexer.TRUE, "true", 6, 2},
					true,
				},
			},
			&PrefixNode{
				lexer.Token{lexer.MINUS, "-", 7, 1},
				&IntNode{
					lexer.Token{lexer.INT, "10", 7, 2},
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
-12 * (7 + 12) * -8;
-(12 + 4);
`
	expected := BlockNode{
		[]Node{
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 1, 4},
				&IntNode{
					lexer.Token{lexer.INT, "10", 1, 1},
					10,
				},
				&IntNode{
					lexer.Token{lexer.INT, "2", 1, 6},
					2,
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASTERISK, "*", 2, 3},
				&IntNode{
					lexer.Token{lexer.INT, "3", 2, 1},
					3,
				},
				&IntNode{
					lexer.Token{lexer.INT, "20", 2, 5},
					20,
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 3, 4},
				&IntNode{
					lexer.Token{lexer.INT, "10", 3, 1},
					10,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 3, 8},
					&IntNode{
						lexer.Token{lexer.INT, "2", 3, 6},
						2,
					},
					&IntNode{
						lexer.Token{lexer.INT, "6", 3, 10},
						6,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 4, 8},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 4, 4},
					&IntNode{
						lexer.Token{lexer.INT, "12", 4, 1},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 4, 8},
						7,
					},
				},
				&IntNode{
					lexer.Token{lexer.INT, "12", 4, 10},
					12,
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 5, 8},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 5, 4},
					&IntNode{
						lexer.Token{lexer.INT, "12", 5, 1},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 5, 8},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 5, 13},
					&IntNode{
						lexer.Token{lexer.INT, "12", 5, 10},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "8", 5, 15},
						8,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 6, 3},
				&IntNode{
					lexer.Token{lexer.INT, "2", 6, 1},
					2,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 6, 15},
					&InfixNode{
						lexer.Token{lexer.ASTERISK, "*", 6, 11},
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 6, 7},
							&IntNode{
								lexer.Token{lexer.INT, "4", 6, 5},
								4,
							},
							&IntNode{
								lexer.Token{lexer.INT, "5", 6, 9},
								5,
							},
						},
						&IntNode{
							lexer.Token{lexer.INT, "6", 6, 13},
							6,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 6, 17},
						7,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 7, 9},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 7, 5},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 7, 1},
						&IntNode{
							lexer.Token{lexer.INT, "12", 7, 2},
							12,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 7, 9},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 7, 14},
					&IntNode{
						lexer.Token{lexer.INT, "12", 7, 11},
						12,
					},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 7, 16},
						&IntNode{
							lexer.Token{lexer.INT, "8", 7, 17},
							8,
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.EQ, "==", 8, 9},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 8, 5},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 8, 1},
						&IntNode{
							lexer.Token{lexer.INT, "12", 8, 2},
							12,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 8, 9},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 8, 15},
					&IntNode{
						lexer.Token{lexer.INT, "12", 8, 12},
						12,
					},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 8, 17},
						&IntNode{
							lexer.Token{lexer.INT, "8", 8, 18},
							8,
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASTERISK, "*", 9, 16},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 9, 5},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 9, 1},
						&IntNode{
							lexer.Token{lexer.INT, "12", 9, 2},
							12,
						},
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 9, 10},
						&IntNode{
							lexer.Token{lexer.INT, "7", 9, 8},
							7,
						},
						&IntNode{
							lexer.Token{lexer.INT, "12", 9, 12},
							12,
						},
					},
				},
				&PrefixNode{
					lexer.Token{lexer.MINUS, "-", 9, 18},
					&IntNode{
						lexer.Token{lexer.INT, "8", 9, 19},
						8,
					},
				},
			},
			&PrefixNode{
				lexer.Token{lexer.MINUS, "-", 10, 1},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 10, 6},
					&IntNode{
						lexer.Token{lexer.INT, "12", 10, 3},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "4", 10, 8},
						4,
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestIfElse(t *testing.T) {
	input := `
if (12 < 4) {
  3 * 20;
  23 >= 20;
};

if (!flag) {
  false;
} else {
  10;
  "test";
};
`
	expected := BlockNode{
		[]Node{
			&ConditionalNode{
				lexer.Token{lexer.IF, "if", 2, 1},
				&InfixNode{
					lexer.Token{lexer.LT, "<", 2, 8},
					&IntNode{
						lexer.Token{lexer.INT, "12", 2, 5},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "4", 2, 10},
						4,
					},
				},
				&BlockNode{
					[]Node{
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 3, 5},
							&IntNode{
								lexer.Token{lexer.INT, "3", 3, 3},
								3,
							},
							&IntNode{
								lexer.Token{lexer.INT, "20", 3, 7},
								20,
							},
						},
						&InfixNode{
							lexer.Token{lexer.GE, ">=", 4, 6},
							&IntNode{
								lexer.Token{lexer.INT, "23", 4, 3},
								23,
							},
							&IntNode{
								lexer.Token{lexer.INT, "20", 4, 9},
								20,
							},
						},
					},
				},
				nil,
			},
			&ConditionalNode{
				lexer.Token{lexer.IF, "if", 7, 1},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 7, 5},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "flag", 7, 6},
						"flag",
					},
				},
				&BlockNode{
					[]Node{
						&BoolNode{
							lexer.Token{lexer.FALSE, "false", 8, 3},
							false,
						},
					},
				},
				&BlockNode{
					[]Node{
						&IntNode{
							lexer.Token{lexer.INT, "10", 10, 3},
							10,
						},
						&StringNode{
							lexer.Token{lexer.STRING, "test", 11, 3},
							"test",
						},
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestLetReturnAssign(t *testing.T) {
	input := `
let test = 10 + 2 * 6;
return !true;
test = !false;
`
	expected := BlockNode{
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5},
						"test",
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 2, 15},
						&IntNode{
							lexer.Token{lexer.INT, "10", 2, 12},
							10,
						},
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 2, 19},
							&IntNode{
								lexer.Token{lexer.INT, "2", 2, 17},
								2,
							},
							&IntNode{
								lexer.Token{lexer.INT, "6", 2, 21},
								6,
							},
						},
					},
				},
			},
			&StatementNode{
				lexer.Token{lexer.RETURN, "return", 3, 1},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 3, 8},
					&BoolNode{
						lexer.Token{lexer.TRUE, "true", 3, 9},
						true,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASSIGN, "=", 4, 6},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 4, 1},
					"test",
				},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 4, 8},
					&BoolNode{
						lexer.Token{lexer.FALSE, "false", 4, 9},
						false,
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestFunction(t *testing.T) {
	input := `
let test = fn(a, b, c) {
  return a * b + c;
};
test = fn() {
  !true;
};
fn(b) {
  return b;
};
`
	expected := BlockNode{
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5},
						"test",
					},
					&FunctionNode{
						lexer.Token{lexer.FUNCTION, "fn", 2, 12},
						[]Node{
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "a", 2, 15},
								"a",
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "b", 2, 18},
								"b",
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "c", 2, 21},
								"c",
							},
						},
						&BlockNode{
							[]Node{
								&StatementNode{
									lexer.Token{lexer.RETURN, "return", 3, 1},
									&InfixNode{
										lexer.Token{lexer.PLUS, "+", 3, 12},
										&InfixNode{
											lexer.Token{lexer.ASTERISK, "*", 3, 4},
											&IdentifierNode{
												lexer.Token{lexer.IDENT, "a", 3, 10},
												"a",
											},
											&IdentifierNode{
												lexer.Token{lexer.IDENT, "b", 3, 14},
												"b",
											},
										},
										&IdentifierNode{
											lexer.Token{lexer.IDENT, "c", 3, 18},
											"c",
										},
									},
								},
							},
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASSIGN, "=", 5, 6},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 5, 1},
					"test",
				},
				&FunctionNode{
					lexer.Token{lexer.FUNCTION, "fn", 5, 8},
					[]Node{},
					&BlockNode{
						[]Node{
							&PrefixNode{
								lexer.Token{lexer.BANG, "!", 6, 3},
								&BoolNode{
									lexer.Token{lexer.TRUE, "true", 6, 4},
									true,
								},
							},
						},
					},
				},
			},
			&FunctionNode{
				lexer.Token{lexer.FUNCTION, "fn", 8, 1},
				[]Node{
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "b", 8, 4},
						"b",
					},
				},
				&BlockNode{
					[]Node{
						&StatementNode{
							lexer.Token{lexer.RETURN, "return", 9, 3},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "b", 9, 10},
								"b",
							},
						},
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestFunctionCall(t *testing.T) {
	input := `
let test = testFn(10 + 2 * 6, test);
17 + 4 * test(25 + 2);
`
	expected := BlockNode{
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5},
						"test",
					},
					&FunctionCallNode{
						lexer.Token{lexer.LPAREN, "(", 2, 18},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "testFn", 2, 12},
							"testFn",
						},
						[]Node{
							&InfixNode{
								lexer.Token{lexer.PLUS, "+", 2, 22},
								&IntNode{
									lexer.Token{lexer.INT, "10", 2, 19},
									10,
								},
								&InfixNode{
									lexer.Token{lexer.ASTERISK, "*", 2, 26},
									&IntNode{
										lexer.Token{lexer.INT, "2", 2, 24},
										2,
									},
									&IntNode{
										lexer.Token{lexer.INT, "6", 2, 28},
										6,
									},
								},
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "test", 2, 31},
								"test",
							},
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 3, 4},
				&IntNode{
					lexer.Token{lexer.INT, "17", 3, 1},
					17,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 3, 8},
					&IntNode{
						lexer.Token{lexer.INT, "4", 3, 6},
						4,
					},
					&FunctionCallNode{
						lexer.Token{lexer.LPAREN, "(", 3, 14},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "test", 3, 10},
							"test",
						},
						[]Node{
							&InfixNode{
								lexer.Token{lexer.PLUS, "+", 3, 18},
								&IntNode{
									lexer.Token{lexer.INT, "25", 3, 15},
									25,
								},
								&IntNode{
									lexer.Token{lexer.INT, "2", 3, 20},
									2,
								},
							},
						},
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}
