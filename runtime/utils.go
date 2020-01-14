package runtime

import (
	"github.com/dyzsr/mylisp/ast"
)

func compare(e1, e2 ast.Expr) bool {
	if e1 == nil {
		return e2 == nil
	}

	switch v1 := e1.(type) {
	case Bool:
		if v2, ok := e2.(Bool); ok {
			return v1 == v2
		}
	case Int:
		if v2, ok := e2.(Int); ok {
			return v1 == v2
		}
	}

	return false
}
