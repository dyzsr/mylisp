package env

import (
	"github.com/dyzsr/mylisp/ast"
)

type Scope struct {
	parent *Scope
	symtab map[string]ast.Expr
}

func NewRootScope() *Scope {
	return NewScope(nil)
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent: parent,
		symtab: make(map[string]ast.Expr),
	}
}

func (s *Scope) Lookup(name string) (ast.Expr, bool) {
	if expr, ok := s.symtab[name]; ok {
		return expr, true
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return nil, false
}

func (s *Scope) Insert(name string, expr ast.Expr) {
	s.symtab[name] = expr
}
