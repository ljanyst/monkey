package evaluator

import "fmt"

func builtinLen(params []Object) (Object, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("len() expects exactly one parameter")
	}

	obj := params[0]
	if obj.Type() != STRING && obj.Type() != ARRAY {
		return nil, fmt.Errorf("The parameter should be either a STRING or an ARRAY")
	}

	return &IntObject{objLen(obj)}, nil
}
