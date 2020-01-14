package runtime

import (
	"github.com/dyzsr/mylisp/ast"
)

type Value interface {
	ast.Expr
	Type() Type
}

type (
	Bool bool

	Int int64

	Symbol struct {
	}

	BuiltinProc struct {
		Name string
	}

	Proc struct {
		*ast.LambdaExpr
		Name string
		env  *EvalEnv
	}
)

func (e Bool) Expr()         {}
func (e Int) Expr()          {}
func (e *Symbol) Expr()      {}
func (e *BuiltinProc) Expr() {}
func (e *Proc) Expr()        {}

func (e Bool) Type() Type         { return BOOLEAN }
func (e Int) Type() Type          { return INTEGER }
func (e *Symbol) Type() Type      { return SYMBOL }
func (e *BuiltinProc) Type() Type { return BUILTIN_PROC }
func (e *Proc) Type() Type        { return PROC }

type Type int

const (
	BOOLEAN = iota
	INTEGER
	SYMBOL
	BUILTIN_PROC
	PROC
)
