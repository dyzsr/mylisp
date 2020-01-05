package runtime

import (
	"mylisp/ast"
)

func compare(e1, e2 ast.Expr) bool {
	if e1 == nil {
		return e2 == nil
	}

	switch v1 := e1.(type) {
	case BoolValue:
		if v2, ok := e2.(BoolValue); ok {
			return v1 == v2
		}
	case IntValue:
		if v2, ok := e2.(IntValue); ok {
			return v1 == v2
		}
	}

	return false
}
