package evaluator

import (
	"testing"
)

func evaluateAndCompareResult(t *testing.T, input []string, expected []Object) bool {

	if len(input) != len(expected) {
		t.Errorf("Input and expected sizes differ")
		return false
	}

	for i := 0; i < len(input); i++ {
		obj, err := EvalString(input[i])
		if err != nil {
			t.Errorf("Unable to evaluate program %d: %s", i, err)
			continue
		}

		exp := expected[i]

		if obj.Type() != exp.Type() || obj.Value() != exp.Value() {
			t.Errorf("Wrong result for test %d: %q. Expected %v, got %v", i, input[i], exp, obj)
			continue
		}
	}

	return true
}

func TestLiteralsAndIdentifiers(t *testing.T) {
	input := []string{
		"10;",
		`"zażółć gęślą jaźń";`,
		"true;",
		"false;",
		"!true;",
		"-10;",
	}

	expected := []Object{
		&IntObject{10},
		&StringObject{"zażółć gęślą jaźń"},
		&BoolObject{true},
		&BoolObject{false},
		&BoolObject{false},
		&IntObject{-10},
	}

	evaluateAndCompareResult(t, input, expected)
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

	evaluateAndCompareResult(t, input, expected)
}
