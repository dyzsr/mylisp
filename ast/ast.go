package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface {
	Pos() *Pos
	End() *Pos
}

type (
	BoolLit struct {
		Value bool
	}

	IntLit struct {
		Value int64
	}

	Quote struct {
		Expr Expr
	}

	Ident struct {
		Name *string
	}

	ListExpr struct {
		List []Expr
	}

	DefineExpr struct {
		Ident *Ident
		Value Expr
	}

	LambdaExpr struct {
		Args []*Ident
		Body []Expr
	}

	CondExpr struct {
		List []*BranchExpr
	}

	BranchExpr struct {
		Else      bool
		Condition Expr
		Body      []Expr
	}
)

func (e *BoolLit) Pos() *Pos    { return nil }
func (e *IntLit) Pos() *Pos     { return nil }
func (e *Ident) Pos() *Pos      { return nil }
func (e *ListExpr) Pos() *Pos   { return nil }
func (e *DefineExpr) Pos() *Pos { return nil }
func (e *LambdaExpr) Pos() *Pos { return nil }
func (e *CondExpr) Pos() *Pos   { return nil }
func (e *BranchExpr) Pos() *Pos { return nil }

func (e *BoolLit) End() *Pos    { return nil }
func (e *IntLit) End() *Pos     { return nil }
func (e *Ident) End() *Pos      { return nil }
func (e *ListExpr) End() *Pos   { return nil }
func (e *DefineExpr) End() *Pos { return nil }
func (e *LambdaExpr) End() *Pos { return nil }
func (e *CondExpr) End() *Pos   { return nil }
func (e *BranchExpr) End() *Pos { return nil }

func NewIdent(name string) *Ident {
	return &Ident{
		Name: SymbolMap(name),
	}
}

func (e *BoolLit) String() string {
	return strconv.FormatBool(e.Value)
}

func (e *IntLit) String() string {
	return strconv.FormatInt(e.Value, 10)
}

func (e *Ident) String() string {
	return *e.Name
}

func (e *ListExpr) String() string {
	var substr []string
	for _, expr := range e.List {
		substr = append(substr, fmt.Sprintf("%s", expr))
	}
	return "[" + strings.Join(substr, " ") + "]"
}

func (e *DefineExpr) String() string {
	return fmt.Sprintf("(define %s %s)", e.Ident, e.Value)
}

func (e *LambdaExpr) String() string {
	return fmt.Sprintf("(lambda %s %s)", e.Args, e.Body)
}

func (e *CondExpr) String() string {
	var substr []string
	for i := range e.List {
		substr = append(substr, fmt.Sprintf("%s", e.List[i]))
	}
	return "(cond " + strings.Join(substr, " ") + ")"
}

func (e *BranchExpr) String() string {
	var substr []string
	if e.Else {
		substr = append(substr, "else")
	} else {
		substr = append(substr, fmt.Sprintf("%s", e.Condition))
	}
	for i := range e.Body {
		substr = append(substr, fmt.Sprintf("%s", e.Body[i]))
	}
	return "(" + strings.Join(substr, " ") + ")"
}
