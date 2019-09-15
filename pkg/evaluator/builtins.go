package evaluator

import (
	"fmt"
	"strings"
)

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

func builtinPrint(params []Object) (Object, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("print() expects at least one parameter")
	}

	if params[0].Type() != STRING {
		return nil, fmt.Errorf("The first parameter needs to be a format string")
	}

	fmtStr := string(params[0].(*StringObject).Value)
	if strings.Index(fmtStr, "%") != -1 {
		return nil, fmt.Errorf("The format string cannot contain '%%' characters")
	}

	count := strings.Count(fmtStr, "#")
	if count != len(params)-1 {
		return nil, fmt.Errorf("Number of items in format string does not match number of parameters")
	}

	fmtStr = strings.ReplaceAll(fmtStr, "#", "%s")
	var fmtParams []interface{}

	for i := 0; i < count; i++ {
		fmtParams = append(fmtParams, params[i+1].Inspect())
	}

	fmt.Printf(fmtStr, fmtParams...)
	fmt.Printf("\n")

	return &IntObject{int64(count)}, nil
}
