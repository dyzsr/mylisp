package runtime

import (
	"github.com/dyzsr/mylisp/ast"
)

type Value interface {
	Type() *Type
}

type (
	Bool bool

	Int int64

	Symbol struct {
		Value *string
	}

	BuiltinProc struct {
		Name string
	}

	Proc struct {
		*ast.LambdaExpr
		Name  string
		outer *ast.Scope
	}
)

func (e Bool) Type() *Type         { return nil }
func (e Int) Type() *Type          { return nil }
func (e *Symbol) Type() *Type      { return nil }
func (e *BuiltinProc) Type() *Type { return nil }
func (e *Proc) Type() *Type        { return nil }
