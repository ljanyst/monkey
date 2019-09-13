package evaluator

import (
	"testing"
)

func evaluateAndCompareResult(t *testing.T, input []string, expected []Object,
	sideEffects []map[string]Object) bool {

	if len(input) != len(expected) {
		t.Errorf("Input and expected sizes differ")
		return false
	}

	status := true

	for i := 0; i < len(input); i++ {
		c := NewContext()
		obj, err := EvalString(input[i], c, "input")
		if err != nil {
			t.Errorf("[test %d] Unable to evaluate program: %s", i, err)
			status = false
			continue
		}

		exp := expected[i]

		if obj.Type() != exp.Type() || obj.Inspect() != exp.Inspect() {
			t.Errorf("[test %d] Wrong result: %q. Expected %v, got %v", i, input[i],
				exp.Inspect(), obj.Inspect())
			status = false
			continue
		}

		if len(sideEffects) == 0 {
			continue
		}

		se := sideEffects[i]
		for k, v := range se {
			obj, err := c.Resolve(k)
			if err != nil {
				t.Errorf("[test %d] Expected to find variable %q but found none", i, v)
				status = false
				continue
			}

			if obj.Type() != v.Type() || obj.Inspect() != v.Inspect() {
				t.Errorf("[test %d] Wrong object in variable %q. Expected %v, got %v", i, k,
					v.Inspect(), obj.Inspect())
				status = false
				continue
			}
		}
	}

	return status
}

func TestLiteralsAndIdentifiers(t *testing.T) {
	input := []string{
		"10;",
		`"zażółć gęślą jaźń";`,
		"true;",
		"false;",
		"!true;",
		"-10;",
		"nil;",
		"'ć';",
	}

	expected := []Object{
		&IntObject{10},
		&StringObject{"zażółć gęślą jaźń"},
		&BoolObject{true},
		&BoolObject{false},
		&BoolObject{false},
		&IntObject{-10},
		&NilObject{},
		&RuneObject{'ć'},
	}

	evaluateAndCompareResult(t, input, expected, []map[string]Object{})
}

func TestInfixPriority(t *testing.T) {
	input := []string{
		"10 + 2;",
		"3 * 20;",
		"10 + 2 * 6;",
		"12 * 7 + 12;",
		"12 * 7 + 12 * 8;",
		"2 + 4 * 5 * 6 * 7;",
		"-12 * 7 + 12 * -8;",
		"-12 * 7 == 12 + -8;",
		"-12 * (7 + 12) * -8;",
		"-(12 + 4);",
	}

	expected := []Object{
		&IntObject{12},
		&IntObject{60},
		&IntObject{22},
		&IntObject{96},
		&IntObject{180},
		&IntObject{842},
		&IntObject{-180},
		&BoolObject{false},
		&IntObject{1824},
		&IntObject{-16},
	}

	evaluateAndCompareResult(t, input, expected, []map[string]Object{})
}

func TestSimpleAssignments(t *testing.T) {
	input := []string{
		"let test1 = 12;",
		"let test2 = 12 * 7 + 12 * 8;",
		"let test3 = -12 * 7 == 12 + -8;",
		"let test4 = 12; test4 = 22; test4 - 1;",
		"let test5 = 18; return 22 + test5;",
	}

	expected := []Object{
		&IntObject{12},
		&IntObject{180},
		&BoolObject{false},
		&IntObject{21},
		&IntObject{40},
	}

	sideEffects := []map[string]Object{
		map[string]Object{"test1": &IntObject{12}},
		map[string]Object{"test2": &IntObject{180}},
		map[string]Object{"test3": &BoolObject{false}},
		map[string]Object{"test4": &IntObject{22}},
		map[string]Object{"test5": &IntObject{18}},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestIfElse(t *testing.T) {
	input := []string{
		"let test1 = 12; if (test1 == 12) { test1 = 3; };",
		"let test2 = 1; if (test2 > 2) { 12 * 4; let test2 = 2; } else { test2 = 3; };",
		"let test3 = 2; if (test3 == 2) { let test3 = 12; test3; };",
	}

	expected := []Object{
		&IntObject{3},
		&IntObject{3},
		&IntObject{12},
	}

	sideEffects := []map[string]Object{
		map[string]Object{"test1": &IntObject{3}},
		map[string]Object{"test2": &IntObject{3}},
		map[string]Object{"test3": &IntObject{2}},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestFunctionDefs(t *testing.T) {
	input0 := `
let test = fn(a, b, c) {
  return a * b + c;
};
`
	input1 := `
let test = 12;
test = fn() {
  !true;
};
`
	input2 := `
fn(b) {
  return b;
};
`

	input := []string{input0, input1, input2}

	expected := []Object{
		&FunctionObject{[]string{"a", "b", "c"}, nil, nil},
		&FunctionObject{[]string{}, nil, nil},
		&FunctionObject{[]string{"b"}, nil, nil},
	}

	sideEffects := []map[string]Object{
		map[string]Object{"test": expected[0]},
		map[string]Object{"test": expected[1]},
		map[string]Object{},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestFunctionCall(t *testing.T) {
	input := []string{`
let adder = fn(x) {
  return fn(y) { return x + y; };
};

let multiplier = fn(x) {
  return fn(y) { return x * y; };
};

let compositor = fn(f1, f2) {
  return fn(x) { return f1(f2(x)); };
};

let func = compositor(adder(5), multiplier(2));

let result = func(3);

result;
`, `
let func = fn(x) {
  if (x > 5) {
    return x;
  } else {
    2;
  };
  return 3;
};

let funcOuter = fn(x) {
  let y = func(x);
  if (y > 5) {
     return x;
  };
  return 1;
};

let x1 = funcOuter(6);
let x2 = funcOuter(1);
`,
	}

	expected := []Object{
		&IntObject{11},
		&IntObject{1},
	}

	sideEffects := []map[string]Object{
		map[string]Object{
			"adder":      &FunctionObject{[]string{"x"}, nil, nil},
			"multiplier": &FunctionObject{[]string{"x"}, nil, nil},
			"compositor": &FunctionObject{[]string{"f1", "f2"}, nil, nil},
			"func":       &FunctionObject{[]string{"x"}, nil, nil},
			"result":     &IntObject{11},
		},
		map[string]Object{
			"x1": &IntObject{6},
			"x2": &IntObject{1},
		},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestStringSlicingAndConcat(t *testing.T) {
	input := []string{`
let test1 = "zażółć";
let test2 = "gęślą";
let test3 = "jaźń";
let test4 = test1 + " " + test2 + " " + test3;
`,
		`
let test1 = "zażółć gęślą jaźń";
let test2 = test1[7:12];
let test3 = "zażółć gęślą jaźń"[7:12];
`,
		`
let test1 = "zażółć";
let test2 = test1[5];
`,
	}

	expected := []Object{
		&StringObject{"zażółć gęślą jaźń"},
		&StringObject{"gęślą"},
		&RuneObject{'ć'},
	}

	sideEffects := []map[string]Object{
		map[string]Object{
			"test1": &StringObject{"zażółć"},
			"test2": &StringObject{"gęślą"},
			"test3": &StringObject{"jaźń"},
			"test4": &StringObject{"zażółć gęślą jaźń"},
		},
		map[string]Object{
			"test1": &StringObject{"zażółć gęślą jaźń"},
			"test2": &StringObject{"gęślą"},
			"test3": &StringObject{"gęślą"},
		},
		map[string]Object{
			"test1": &StringObject{"zażółć"},
			"test2": &RuneObject{'ć'},
		},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestArraySlicingAndConcat(t *testing.T) {
	input := []string{`
let test1 = {"zażółć"};
let test2 = {"gęślą"};
let test3 = {"jaźń"};
let test4 = test1 + test2 + test3;
`,
		`
let test1 = {"zażółć", 12, 3 + 2, 'ł'};
let test2 = test1[2:4];
let test3 = {"zażółć", 12, 3 + 2, 'ł'}[2:4];
`,
		`
let test1 = {"zażółć", 12, 3 + 2, 'ł'};
let test2 = test1[2];
`,
	}

	expected := []Object{
		&ArrayObject{
			[]Object{
				&StringObject{"zażółć"},
				&StringObject{"gęślą"},
				&StringObject{"jaźń"},
			},
		},
		&ArrayObject{
			[]Object{
				&IntObject{5},
				&RuneObject{'ł'},
			},
		},
		&IntObject{5},
	}

	sideEffects := []map[string]Object{
		map[string]Object{
			"test1": &ArrayObject{[]Object{&StringObject{"zażółć"}}},
			"test2": &ArrayObject{[]Object{&StringObject{"gęślą"}}},
			"test3": &ArrayObject{[]Object{&StringObject{"jaźń"}}},
			"test4": &ArrayObject{
				[]Object{
					&StringObject{"zażółć"},
					&StringObject{"gęślą"},
					&StringObject{"jaźń"},
				},
			},
		},
		map[string]Object{
			"test1": &ArrayObject{
				[]Object{
					&StringObject{"zażółć"},
					&IntObject{12},
					&IntObject{5},
					&RuneObject{'ł'},
				},
			},
			"test2": &ArrayObject{
				[]Object{
					&IntObject{5},
					&RuneObject{'ł'},
				},
			},
			"test3": &ArrayObject{
				[]Object{
					&IntObject{5},
					&RuneObject{'ł'},
				},
			},
		},
		map[string]Object{
			"test1": &ArrayObject{
				[]Object{
					&StringObject{"zażółć"},
					&IntObject{12},
					&IntObject{5},
					&RuneObject{'ł'},
				},
			},
			"test2": &IntObject{5},
		},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestLoops(t *testing.T) {
	input := []string{`
let test = {};
for (let i = 0; i < 5; i = i + 1) {
  test = test + {i};
};
`, `
let test1 = {};
let test2 = "gęślą";
let i = 0;
for (; i < 5;) {
  test1 = test1 + {test2[i]};
  i = i + 1;
};
`, `
let test = {};
for (let i = 0; i < 5; i = i + 1) {
  if (i == 2) {
    return 71;
  };
  test = test + {i};
};
`, `
let test = {};
for (let i = 0; i < 5; i = i + 1) {
  if (i == 2) {
    break;
  };
  test = test + {i};
};
`, `
let test = {};
for (let i = 0; i < 5; i = i + 1) {
  if (i == 2) {
    continue;
  };
  test = test + {i};
};
`,
	}

	expected := []Object{
		&ArrayObject{
			[]Object{
				&IntObject{0},
				&IntObject{1},
				&IntObject{2},
				&IntObject{3},
				&IntObject{4},
			},
		},
		&IntObject{5},
		&IntObject{71},
		&NilObject{},
		&ArrayObject{
			[]Object{
				&IntObject{0},
				&IntObject{1},
				&IntObject{3},
				&IntObject{4},
			},
		},
	}

	sideEffects := []map[string]Object{
		map[string]Object{
			"test": &ArrayObject{
				[]Object{
					&IntObject{0},
					&IntObject{1},
					&IntObject{2},
					&IntObject{3},
					&IntObject{4},
				},
			},
		},
		map[string]Object{
			"test1": &ArrayObject{
				[]Object{
					&RuneObject{'g'},
					&RuneObject{'ę'},
					&RuneObject{'ś'},
					&RuneObject{'l'},
					&RuneObject{'ą'},
				},
			},
			"test2": &StringObject{"gęślą"},
			"i":     &IntObject{5},
		},
		map[string]Object{
			"test": &ArrayObject{
				[]Object{
					&IntObject{0},
					&IntObject{1},
				},
			},
		},
		map[string]Object{
			"test": &ArrayObject{
				[]Object{
					&IntObject{0},
					&IntObject{1},
				},
			},
		},
		map[string]Object{
			"test": &ArrayObject{
				[]Object{
					&IntObject{0},
					&IntObject{1},
					&IntObject{3},
					&IntObject{4},
				},
			},
		},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}

func TestLogic(t *testing.T) {
	input := []string{`
let test = {};
for (let i = 0; i < 6 || i > 6; i = i + 1) {
  if (i >= 2 && i <= 3) {
    continue;
  };
  test = test + {i};
};
`,
	}

	expected := []Object{
		&ArrayObject{
			[]Object{
				&IntObject{0},
				&IntObject{1},
				&IntObject{4},
				&IntObject{5},
			},
		},
	}

	sideEffects := []map[string]Object{
		map[string]Object{
			"test": &ArrayObject{
				[]Object{
					&IntObject{0},
					&IntObject{1},
					&IntObject{4},
					&IntObject{5},
				},
			},
		},
	}

	evaluateAndCompareResult(t, input, expected, sideEffects)
}
