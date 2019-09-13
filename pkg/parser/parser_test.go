package parser

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ljanyst/monkey/pkg/lexer"
)

var printAst = flag.Bool("print-ast", false, "print the AST")
var printProgram = flag.Bool("print-program", false, "print the parsed program")

const input = "input"

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
	l := lexer.NewLexerFromString(input, "input")
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
nil;
'ć';
`

	expected := BlockNode{
		true,
		[]Node{
			&IntNode{
				lexer.Token{lexer.INT, "10", 1, 1, &input},
				10,
			},
			&StringNode{
				lexer.Token{lexer.STRING, "zażółć gęślą jaźń", 2, 1, &input},
				"zażółć gęślą jaźń",
			},
			&IdentifierNode{
				lexer.Token{lexer.IDENT, "test", 3, 1, &input},
				"test",
			},
			&BoolNode{
				lexer.Token{lexer.TRUE, "true", 4, 1, &input},
				true,
			},
			&BoolNode{
				lexer.Token{lexer.FALSE, "false", 5, 1, &input},
				false,
			},
			&PrefixNode{
				lexer.Token{lexer.BANG, "!", 6, 1, &input},
				&BoolNode{
					lexer.Token{lexer.TRUE, "true", 6, 2, &input},
					true,
				},
			},
			&PrefixNode{
				lexer.Token{lexer.MINUS, "-", 7, 1, &input},
				&IntNode{
					lexer.Token{lexer.INT, "10", 7, 2, &input},
					10,
				},
			},
			&NilNode{
				lexer.Token{lexer.NIL, "nil", 8, 1, &input},
			},
			&RuneNode{
				lexer.Token{lexer.RUNE, "ć", 9, 1, &input},
				'ć',
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
-12 * 7 == 12 + -8 && -32 > 3;
`
	expected := BlockNode{
		true,
		[]Node{
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 1, 4, &input},
				&IntNode{
					lexer.Token{lexer.INT, "10", 1, 1, &input},
					10,
				},
				&IntNode{
					lexer.Token{lexer.INT, "2", 1, 6, &input},
					2,
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASTERISK, "*", 2, 3, &input},
				&IntNode{
					lexer.Token{lexer.INT, "3", 2, 1, &input},
					3,
				},
				&IntNode{
					lexer.Token{lexer.INT, "20", 2, 5, &input},
					20,
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 3, 4, &input},
				&IntNode{
					lexer.Token{lexer.INT, "10", 3, 1, &input},
					10,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 3, 8, &input},
					&IntNode{
						lexer.Token{lexer.INT, "2", 3, 6, &input},
						2,
					},
					&IntNode{
						lexer.Token{lexer.INT, "6", 3, 10, &input},
						6,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 4, 8, &input},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 4, 4, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 4, 1, &input},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 4, 8, &input},
						7,
					},
				},
				&IntNode{
					lexer.Token{lexer.INT, "12", 4, 10, &input},
					12,
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 5, 8, &input},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 5, 4, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 5, 1, &input},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 5, 8, &input},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 5, 13, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 5, 10, &input},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "8", 5, 15, &input},
						8,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 6, 3, &input},
				&IntNode{
					lexer.Token{lexer.INT, "2", 6, 1, &input},
					2,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 6, 15, &input},
					&InfixNode{
						lexer.Token{lexer.ASTERISK, "*", 6, 11, &input},
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 6, 7, &input},
							&IntNode{
								lexer.Token{lexer.INT, "4", 6, 5, &input},
								4,
							},
							&IntNode{
								lexer.Token{lexer.INT, "5", 6, 9, &input},
								5,
							},
						},
						&IntNode{
							lexer.Token{lexer.INT, "6", 6, 13, &input},
							6,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 6, 17, &input},
						7,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 7, 9, &input},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 7, 5, &input},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 7, 1, &input},
						&IntNode{
							lexer.Token{lexer.INT, "12", 7, 2, &input},
							12,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 7, 9, &input},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 7, 14, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 7, 11, &input},
						12,
					},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 7, 16, &input},
						&IntNode{
							lexer.Token{lexer.INT, "8", 7, 17, &input},
							8,
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.EQ, "==", 8, 9, &input},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 8, 5, &input},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 8, 1, &input},
						&IntNode{
							lexer.Token{lexer.INT, "12", 8, 2, &input},
							12,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "7", 8, 9, &input},
						7,
					},
				},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 8, 15, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 8, 12, &input},
						12,
					},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 8, 17, &input},
						&IntNode{
							lexer.Token{lexer.INT, "8", 8, 18, &input},
							8,
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASTERISK, "*", 9, 16, &input},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 9, 5, &input},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 9, 1, &input},
						&IntNode{
							lexer.Token{lexer.INT, "12", 9, 2, &input},
							12,
						},
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 9, 10, &input},
						&IntNode{
							lexer.Token{lexer.INT, "7", 9, 8, &input},
							7,
						},
						&IntNode{
							lexer.Token{lexer.INT, "12", 9, 12, &input},
							12,
						},
					},
				},
				&PrefixNode{
					lexer.Token{lexer.MINUS, "-", 9, 18, &input},
					&IntNode{
						lexer.Token{lexer.INT, "8", 9, 19, &input},
						8,
					},
				},
			},
			&PrefixNode{
				lexer.Token{lexer.MINUS, "-", 10, 1, &input},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 10, 6, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 10, 3, &input},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "4", 10, 8, &input},
						4,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.AND, "&&", 11, 20, &input},
				&InfixNode{
					lexer.Token{lexer.EQ, "==", 11, 9, &input},
					&InfixNode{
						lexer.Token{lexer.ASTERISK, "*", 11, 5, &input},
						&PrefixNode{
							lexer.Token{lexer.MINUS, "-", 11, 1, &input},
							&IntNode{
								lexer.Token{lexer.INT, "12", 11, 2, &input},
								12,
							},
						},
						&IntNode{
							lexer.Token{lexer.INT, "7", 11, 9, &input},
							7,
						},
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 11, 15, &input},
						&IntNode{
							lexer.Token{lexer.INT, "12", 11, 12, &input},
							12,
						},
						&PrefixNode{
							lexer.Token{lexer.MINUS, "-", 11, 17, &input},
							&IntNode{
								lexer.Token{lexer.INT, "8", 11, 18, &input},
								8,
							},
						},
					},
				},
				&InfixNode{
					lexer.Token{lexer.GT, ">", 11, 27, &input},
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 11, 23, &input},
						&IntNode{
							lexer.Token{lexer.INT, "32", 11, 24, &input},
							32,
						},
					},
					&IntNode{
						lexer.Token{lexer.INT, "3", 11, 29, &input},
						3,
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
		true,
		[]Node{
			&ConditionalNode{
				lexer.Token{lexer.IF, "if", 2, 1, &input},
				&InfixNode{
					lexer.Token{lexer.LT, "<", 2, 8, &input},
					&IntNode{
						lexer.Token{lexer.INT, "12", 2, 5, &input},
						12,
					},
					&IntNode{
						lexer.Token{lexer.INT, "4", 2, 10, &input},
						4,
					},
				},
				&BlockNode{
					false,
					[]Node{
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 3, 5, &input},
							&IntNode{
								lexer.Token{lexer.INT, "3", 3, 3, &input},
								3,
							},
							&IntNode{
								lexer.Token{lexer.INT, "20", 3, 7, &input},
								20,
							},
						},
						&InfixNode{
							lexer.Token{lexer.GE, ">=", 4, 6, &input},
							&IntNode{
								lexer.Token{lexer.INT, "23", 4, 3, &input},
								23,
							},
							&IntNode{
								lexer.Token{lexer.INT, "20", 4, 9, &input},
								20,
							},
						},
					},
				},
				nil,
			},
			&ConditionalNode{
				lexer.Token{lexer.IF, "if", 7, 1, &input},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 7, 5, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "flag", 7, 6, &input},
						"flag",
					},
				},
				&BlockNode{
					false,
					[]Node{
						&BoolNode{
							lexer.Token{lexer.FALSE, "false", 8, 3, &input},
							false,
						},
					},
				},
				&BlockNode{
					false,
					[]Node{
						&IntNode{
							lexer.Token{lexer.INT, "10", 10, 3, &input},
							10,
						},
						&StringNode{
							lexer.Token{lexer.STRING, "test", 11, 3, &input},
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
		true,
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1, &input},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5, &input},
						"test",
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 2, 15, &input},
						&IntNode{
							lexer.Token{lexer.INT, "10", 2, 12, &input},
							10,
						},
						&InfixNode{
							lexer.Token{lexer.ASTERISK, "*", 2, 19, &input},
							&IntNode{
								lexer.Token{lexer.INT, "2", 2, 17, &input},
								2,
							},
							&IntNode{
								lexer.Token{lexer.INT, "6", 2, 21, &input},
								6,
							},
						},
					},
				},
			},
			&StatementNode{
				lexer.Token{lexer.RETURN, "return", 3, 1, &input},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 3, 8, &input},
					&BoolNode{
						lexer.Token{lexer.TRUE, "true", 3, 9, &input},
						true,
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.ASSIGN, "=", 4, 6, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 4, 1, &input},
					"test",
				},
				&PrefixNode{
					lexer.Token{lexer.BANG, "!", 4, 8, &input},
					&BoolNode{
						lexer.Token{lexer.FALSE, "false", 4, 9, &input},
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
		true,
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1, &input},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5, &input},
						"test",
					},
					&FunctionNode{
						lexer.Token{lexer.FUNCTION, "fn", 2, 12, &input},
						[]Node{
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "a", 2, 15, &input},
								"a",
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "b", 2, 18, &input},
								"b",
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "c", 2, 21, &input},
								"c",
							},
						},
						&BlockNode{
							false,
							[]Node{
								&StatementNode{
									lexer.Token{lexer.RETURN, "return", 3, 1, &input},
									&InfixNode{
										lexer.Token{lexer.PLUS, "+", 3, 12, &input},
										&InfixNode{
											lexer.Token{lexer.ASTERISK, "*", 3, 4, &input},
											&IdentifierNode{
												lexer.Token{lexer.IDENT, "a", 3, 10, &input},
												"a",
											},
											&IdentifierNode{
												lexer.Token{lexer.IDENT, "b", 3, 14, &input},
												"b",
											},
										},
										&IdentifierNode{
											lexer.Token{lexer.IDENT, "c", 3, 18, &input},
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
				lexer.Token{lexer.ASSIGN, "=", 5, 6, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 5, 1, &input},
					"test",
				},
				&FunctionNode{
					lexer.Token{lexer.FUNCTION, "fn", 5, 8, &input},
					[]Node{},
					&BlockNode{
						false,
						[]Node{
							&PrefixNode{
								lexer.Token{lexer.BANG, "!", 6, 3, &input},
								&BoolNode{
									lexer.Token{lexer.TRUE, "true", 6, 4, &input},
									true,
								},
							},
						},
					},
				},
			},
			&FunctionNode{
				lexer.Token{lexer.FUNCTION, "fn", 8, 1, &input},
				[]Node{
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "b", 8, 4, &input},
						"b",
					},
				},
				&BlockNode{
					false,
					[]Node{
						&StatementNode{
							lexer.Token{lexer.RETURN, "return", 9, 3, &input},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "b", 9, 10, &input},
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
		true,
		[]Node{
			&StatementNode{
				lexer.Token{lexer.LET, "let", 2, 1, &input},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 2, 5, &input},
						"test",
					},
					&FunctionCallNode{
						lexer.Token{lexer.LPAREN, "(", 2, 18, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "testFn", 2, 12, &input},
							"testFn",
						},
						[]Node{
							&InfixNode{
								lexer.Token{lexer.PLUS, "+", 2, 22, &input},
								&IntNode{
									lexer.Token{lexer.INT, "10", 2, 19, &input},
									10,
								},
								&InfixNode{
									lexer.Token{lexer.ASTERISK, "*", 2, 26, &input},
									&IntNode{
										lexer.Token{lexer.INT, "2", 2, 24, &input},
										2,
									},
									&IntNode{
										lexer.Token{lexer.INT, "6", 2, 28, &input},
										6,
									},
								},
							},
							&IdentifierNode{
								lexer.Token{lexer.IDENT, "test", 2, 31, &input},
								"test",
							},
						},
					},
				},
			},
			&InfixNode{
				lexer.Token{lexer.PLUS, "+", 3, 4, &input},
				&IntNode{
					lexer.Token{lexer.INT, "17", 3, 1, &input},
					17,
				},
				&InfixNode{
					lexer.Token{lexer.ASTERISK, "*", 3, 8, &input},
					&IntNode{
						lexer.Token{lexer.INT, "4", 3, 6, &input},
						4,
					},
					&FunctionCallNode{
						lexer.Token{lexer.LPAREN, "(", 3, 14, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "test", 3, 10, &input},
							"test",
						},
						[]Node{
							&InfixNode{
								lexer.Token{lexer.PLUS, "+", 3, 18, &input},
								&IntNode{
									lexer.Token{lexer.INT, "25", 3, 15, &input},
									25,
								},
								&IntNode{
									lexer.Token{lexer.INT, "2", 3, 20, &input},
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

func TestSlicing(t *testing.T) {
	input := `
test[1];
test[1+6];
test[0:4];
test[0+1:4-2];
`
	expected := BlockNode{
		true,
		[]Node{
			&SliceNode{
				lexer.Token{lexer.LBRACKET, "[", 1, 5, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 1, 1, &input},
					"test1",
				},
				&IntNode{
					lexer.Token{lexer.INT, "1", 1, 6, &input},
					1,
				},
				nil,
			},
			&SliceNode{
				lexer.Token{lexer.LBRACKET, "[", 2, 5, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 2, 1, &input},
					"test",
				},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 2, 7, &input},
					&IntNode{
						lexer.Token{lexer.INT, "1", 2, 6, &input},
						1,
					},
					&IntNode{
						lexer.Token{lexer.INT, "6", 2, 8, &input},
						6,
					},
				},
				nil,
			},
			&SliceNode{
				lexer.Token{lexer.LBRACKET, "[", 3, 5, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 3, 1, &input},
					"test",
				},
				&IntNode{
					lexer.Token{lexer.INT, "0", 3, 6, &input},
					0,
				},
				&IntNode{
					lexer.Token{lexer.INT, "4", 3, 8, &input},
					4,
				},
			},
			&SliceNode{
				lexer.Token{lexer.LBRACKET, "[", 4, 5, &input},
				&IdentifierNode{
					lexer.Token{lexer.IDENT, "test", 4, 1, &input},
					"test",
				},
				&InfixNode{
					lexer.Token{lexer.PLUS, "+", 4, 7, &input},
					&IntNode{
						lexer.Token{lexer.INT, "0", 4, 6, &input},
						0,
					},
					&IntNode{
						lexer.Token{lexer.INT, "1", 4, 8, &input},
						1,
					},
				},
				&InfixNode{
					lexer.Token{lexer.MINUS, "-", 4, 11, &input},
					&IntNode{
						lexer.Token{lexer.INT, "4", 4, 10, &input},
						4,
					},
					&IntNode{
						lexer.Token{lexer.INT, "2", 4, 12, &input},
						2,
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestArrays(t *testing.T) {
	input := `
{};
{1};
{-1, 2 + 4, test};
`
	expected := BlockNode{
		true,
		[]Node{
			&ArrayNode{
				lexer.Token{lexer.LBRACE, "{", 2, 1, &input},
				[]Node{},
			},
			&ArrayNode{
				lexer.Token{lexer.LBRACE, "{", 3, 1, &input},
				[]Node{
					&IntNode{
						lexer.Token{lexer.INT, "1", 3, 2, &input},
						1,
					},
				},
			},
			&ArrayNode{
				lexer.Token{lexer.LBRACE, "{", 4, 1, &input},
				[]Node{
					&PrefixNode{
						lexer.Token{lexer.MINUS, "-", 4, 2, &input},
						&IntNode{
							lexer.Token{lexer.INT, "1", 4, 3, &input},
							1,
						},
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 4, 8, &input},
						&IntNode{
							lexer.Token{lexer.INT, "2", 4, 6, &input},
							2,
						},
						&IntNode{
							lexer.Token{lexer.INT, "4", 4, 10, &input},
							4,
						},
					},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "test", 4, 13, &input},
						"test",
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}

func TestLoops(t *testing.T) {
	input := `
for (let a = 10; a < 23; a = a + 1) {};
for (; a < 23; a = a + 1) {};
for (let a = 10; a < 23;) {};
for (; a < 23;) {};
for (; a < 23;) {
  break;
  continue;
};
`
	expected := BlockNode{
		true,
		[]Node{
			&LoopNode{
				lexer.Token{lexer.FOR, "for", 2, 1, &input},
				&StatementNode{
					lexer.Token{lexer.LET, "let", 2, 6, &input},
					&InfixNode{
						lexer.Token{lexer.ASSIGN, "=", 2, 12, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "a", 2, 10, &input},
							"a",
						},
						&IntNode{
							lexer.Token{lexer.INT, "10", 2, 14, &input},
							10,
						},
					},
				},
				&InfixNode{
					lexer.Token{lexer.LT, "<", 2, 20, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 2, 18, &input},
						"a",
					},
					&IntNode{
						lexer.Token{lexer.INT, "23", 2, 22, &input},
						23,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 2, 28, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 2, 26, &input},
						"a",
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 2, 32, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "a", 2, 30, &input},
							"a",
						},
						&IntNode{
							lexer.Token{lexer.INT, "1", 2, 34, &input},
							1,
						},
					},
				},
				&BlockNode{
					false,
					[]Node{},
				},
			},
			&LoopNode{
				lexer.Token{lexer.FOR, "for", 3, 1, &input},
				nil,
				&InfixNode{
					lexer.Token{lexer.LT, "<", 3, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 3, 8, &input},
						"a",
					},
					&IntNode{
						lexer.Token{lexer.INT, "23", 3, 12, &input},
						23,
					},
				},
				&InfixNode{
					lexer.Token{lexer.ASSIGN, "=", 3, 18, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 3, 16, &input},
						"a",
					},
					&InfixNode{
						lexer.Token{lexer.PLUS, "+", 3, 22, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "a", 3, 20, &input},
							"a",
						},
						&IntNode{
							lexer.Token{lexer.INT, "1", 3, 24, &input},
							1,
						},
					},
				},
				&BlockNode{
					false,
					[]Node{},
				},
			},
			&LoopNode{
				lexer.Token{lexer.FOR, "for", 4, 1, &input},
				&StatementNode{
					lexer.Token{lexer.LET, "let", 4, 6, &input},
					&InfixNode{
						lexer.Token{lexer.ASSIGN, "=", 4, 12, &input},
						&IdentifierNode{
							lexer.Token{lexer.IDENT, "a", 4, 10, &input},
							"a",
						},
						&IntNode{
							lexer.Token{lexer.INT, "10", 4, 14, &input},
							10,
						},
					},
				},
				&InfixNode{
					lexer.Token{lexer.LT, "<", 4, 20, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 4, 18, &input},
						"a",
					},
					&IntNode{
						lexer.Token{lexer.INT, "23", 4, 22, &input},
						23,
					},
				},
				nil,
				&BlockNode{
					false,
					[]Node{},
				},
			},
			&LoopNode{
				lexer.Token{lexer.FOR, "for", 5, 1, &input},
				nil,
				&InfixNode{
					lexer.Token{lexer.LT, "<", 5, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 5, 8, &input},
						"a",
					},
					&IntNode{
						lexer.Token{lexer.INT, "23", 5, 12, &input},
						23,
					},
				},
				nil,
				&BlockNode{
					false,
					[]Node{},
				},
			},
			&LoopNode{
				lexer.Token{lexer.FOR, "for", 6, 1, &input},
				nil,
				&InfixNode{
					lexer.Token{lexer.LT, "<", 6, 10, &input},
					&IdentifierNode{
						lexer.Token{lexer.IDENT, "a", 6, 8, &input},
						"a",
					},
					&IntNode{
						lexer.Token{lexer.INT, "23", 6, 12, &input},
						23,
					},
				},
				nil,
				&BlockNode{
					false,
					[]Node{
						&StatementNode{
							lexer.Token{lexer.BREAK, "break", 7, 3, &input},
							nil,
						},
						&StatementNode{
							lexer.Token{lexer.CONTINUE, "continue", 8, 3, &input},
							nil,
						},
					},
				},
			},
		},
	}

	parseAndCompareAst(t, input, &expected)
}
