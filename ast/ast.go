package ast

type Expr interface {
	Expr()
}

type (
	BoolLit struct {
		Value bool
	}

	NumLit struct {
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

func (e *BoolLit) Expr()    {}
func (e *NumLit) Expr()     {}
func (e *Ident) Expr()      {}
func (e *ListExpr) Expr()   {}
func (e *DefineExpr) Expr() {}
func (e *LambdaExpr) Expr() {}
func (e *CondExpr) Expr()   {}
func (e *BranchExpr) Expr() {}
