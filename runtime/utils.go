package runtime

import (
	"mylisp/ast"
)

func compare(e1, e2 ast.Expr) bool {
	if e1 == nil {
		return e2 == nil
	}

	switch v1 := e1.(type) {
	case *ast.IntLit:
		if v2, ok := e2.(*ast.IntLit); ok {
			return v1.Value == v2.Value
		}
		return false
	}

	return false
}
