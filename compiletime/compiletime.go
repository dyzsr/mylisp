package compiletime

import (
	"github.com/dyzsr/mylisp/ast"
)

type CompileTime struct {
	scope *ast.Scope
}

func NewCompileTime() *CompileTime {
	scope := ast.NewRootScope()
	for k, v := range builtinTransformerMap() {
		scope.Insert(ast.SymbolMap(k), v)
	}
	return &CompileTime{
		scope: scope,
	}
}

func (c *CompileTime) Eval(input ast.Expr) (ast.Expr, error) {
	return transform(c.scope, input)
}

func transform(scope *ast.Scope, input ast.Expr) (ast.Expr, error) {
	// skip non-list expression:
	// either atomic or has already been transformed
	origList, ok := input.(*ast.ListExpr)
	if !ok {
		return input, nil
	}
	if len(origList.List) == 0 {
		return input, nil
	}

	// transform sub-expressions first
	var list []ast.Expr
	for i := range origList.List {
		expr, err := transform(scope, origList.List[i])
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
	}

	result := &ast.ListExpr{List: list}

	// apply transformers
	first := origList.List[0]
	ident, ok := first.(*ast.Ident)
	if !ok { // the leading expression is not an identifier
		return result, nil
	}
	value, ok := scope.Lookup(ident.Name)
	if !ok {
		return result, nil
	}
	transformer, ok := value.(Transformer)
	if !ok {
		panic("invalid tranformer type")
	}
	return transformer.Transform(scope, result)
}
