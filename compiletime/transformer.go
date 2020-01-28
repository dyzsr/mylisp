package compiletime

import "github.com/dyzsr/mylisp/ast"

type Transformer interface {
	Transform(*ast.Scope, *ast.ListExpr) (ast.Expr, error)
}
