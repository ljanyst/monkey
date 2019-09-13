// Code generated by "stringer -type ObjectType object.go"; DO NOT EDIT.

package evaluator

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INT-0]
	_ = x[BOOL-1]
	_ = x[STRING-2]
	_ = x[EXIT-3]
	_ = x[FUNCTION-4]
	_ = x[NIL-5]
	_ = x[RUNE-6]
	_ = x[ARRAY-7]
}

const _ObjectType_name = "INTBOOLSTRINGEXITFUNCTIONNILRUNEARRAY"

var _ObjectType_index = [...]uint8{0, 3, 7, 13, 17, 25, 28, 32, 37}

func (i ObjectType) String() string {
	if i < 0 || i >= ObjectType(len(_ObjectType_index)-1) {
		return "ObjectType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}
