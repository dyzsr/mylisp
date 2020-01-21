package ast

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

	Ident struct {
		Name string
	}

	ListExpr struct {
		SubExprList []Expr
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
		BranchList []*BranchExpr
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
