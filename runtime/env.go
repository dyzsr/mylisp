package runtime

import (
	"mylisp/ast"
)

type EvalEnv struct {
	parent *EvalEnv
	symtab map[string]ast.Expr
}

func NewRootEnv() *EvalEnv {
	return NewEvalEnv(nil)
}

func NewEvalEnv(parent *EvalEnv) *EvalEnv {
	return &EvalEnv{
		parent: parent,
		symtab: initSymtab(),
	}
}

func (e *EvalEnv) lookup(name string) (ast.Expr, bool) {
	if value, ok := e.symtab[name]; ok {
		return value, true
	}
	if e.parent != nil {
		return e.parent.lookup(name)
	}
	return nil, false
}

func (e *EvalEnv) insert(name string, value ast.Expr) {
	e.symtab[name] = value
}
