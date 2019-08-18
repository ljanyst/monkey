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
		obj, err := EvalString(input[i], c)
		if err != nil {
			t.Errorf("[test %d] Unable to evaluate program: %s", i, err)
			status = false
			continue
		}

		exp := expected[i]

		if obj.Type() != exp.Type() || obj.Value() != exp.Value() {
			t.Errorf("[test %d] Wrong result: %q. Expected %v, got %v", i, input[i], exp, obj)
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

			if obj.Type() != v.Type() || obj.Value() != v.Value() {
				t.Errorf("[test %d] Wrong object in variable %q. Expected %v, got %v", i, k, v, obj)
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
	}

	expected := []Object{
		&IntObject{10},
		&StringObject{"zażółć gęślą jaźń"},
		&BoolObject{true},
		&BoolObject{false},
		&BoolObject{false},
		&IntObject{-10},
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
