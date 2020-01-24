package transformer

import (
	"github.com/dyzsr/mylisp/ast"
)

type Transformer struct {
	scope *ast.Scope
}

func NewTransformer() *Transformer {
	return &Transformer{
		scope: ast.NewRootScope(),
	}
}
