package compiletime

import (
	"fmt"

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
	if ident, ok := input.(*ast.Ident); ok {
		if _, ok := scope.Lookup(ident.Name); ok {
			return nil, fmt.Errorf("%s: bad syntax", *ident.Name)
		}
	}
	// skip non-list expression:
	// either atomic or has already been transformed
	origList, ok := input.(*ast.ListExpr)
	if !ok {
		return input, nil
	}
	if len(origList.List) == 0 {
		return input, nil
	}

	// apply transformers
	var intermediate ast.Expr = input
	first := origList.List[0]
	if ident, ok := first.(*ast.Ident); ok {
		if value, ok := scope.Lookup(ident.Name); ok {
			if transformer, ok := value.(Transformer); ok {
				var err error
				intermediate, err = transformer.Transform(scope, origList)
				if err != nil {
					return nil, err
				}
			} else {
				panic("invalid tranformer type")
			}
		}
	}

	// continue to transform intermediate result
	switch expr := intermediate.(type) {

	case *ast.ListExpr:
		var list []ast.Expr
		for _, item := range expr.List {
			expr, err := transform(scope, item)
			if err != nil {
				return nil, err
			}
			list = append(list, expr)
		}
		return &ast.ListExpr{List: list}, nil

	case *ast.DefineExpr:
		value, err := transform(scope, expr.Value)
		if err != nil {
			return nil, err
		}
		return &ast.DefineExpr{Ident: expr.Ident, Value: value}, nil

	case *ast.SetExpr:
		value, err := transform(scope, expr.Value)
		if err != nil {
			return nil, err
		}
		return &ast.SetExpr{Ident: expr.Ident, Value: value}, nil

	case *ast.LambdaExpr:
		var body []ast.Expr
		for i := range expr.Body {
			result, err := transform(scope, expr.Body[i])
			if err != nil {
				return nil, err
			}
			body = append(body, result)
		}
		return &ast.LambdaExpr{Args: expr.Args, Body: body}, nil

	case *ast.CondExpr:
		var branchList []*ast.BranchExpr
		for _, branch := range expr.List {
			var condition ast.Expr
			if !branch.Else {
				var err error
				condition, err = transform(scope, branch.Condition)
				if err != nil {
					return nil, err
				}
			}

			var body []ast.Expr
			for i := range branch.Body {
				result, err := transform(scope, branch.Body[i])
				if err != nil {
					return nil, err
				}
				body = append(body, result)
			}

			branchList = append(branchList, &ast.BranchExpr{
				Else:      branch.Else,
				Condition: condition,
				Body:      body,
			})
		}
		return &ast.CondExpr{List: branchList}, nil

	case *ast.Quote:
		return intermediate, nil

	default:
		return intermediate, nil
	}
	// transform sub-expressions
}
