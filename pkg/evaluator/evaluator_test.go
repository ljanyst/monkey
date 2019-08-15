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
			t.Errorf("Wrong result. Expected %v, got %v", exp, obj)
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
