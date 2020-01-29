package compiletime

import (
	"errors"

	"github.com/dyzsr/mylisp/ast"
)

type BuiltinTransformer struct {
	name string
	proc func(*ast.Scope, *ast.ListExpr) (ast.Expr, error)
}

func (t BuiltinTransformer) Transform(scope *ast.Scope, list *ast.ListExpr) (ast.Expr, error) {
	return t.proc(scope, list)
}

var (
	builtinDefine = BuiltinTransformer{name: "define", proc: defineSyntax}
	builtinSet    = BuiltinTransformer{name: "set!", proc: setSyntax}
	builtinLambda = BuiltinTransformer{name: "lambda", proc: lambdaSyntax}
	builtinCond   = BuiltinTransformer{name: "cond", proc: condSyntax}
	builtinQuote  = BuiltinTransformer{name: "quote", proc: quoteSyntax}
)

func builtinTransformerMap() map[string]Transformer {
	return map[string]Transformer{
		"define": builtinDefine,
		"set!":   builtinSet,
		"lambda": builtinLambda,
		"cond":   builtinCond,
		"quote":  builtinQuote,
	}
}

func defineSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	badSyntaxErr := errors.New("define: bad syntax")
	origList := input.List
	if len(origList) != 3 {
		return nil, badSyntaxErr
	}
	// ensure an identifier is given
	ident, ok := origList[1].(*ast.Ident)
	if !ok {
		return nil, badSyntaxErr
	}

	return &ast.DefineExpr{
		Ident: ident,
		Value: origList[2],
	}, nil
}

func setSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	badSyntaxErr := errors.New("set!: bad syntax")
	origList := input.List
	if len(origList) != 3 {
		return nil, badSyntaxErr
	}
	// ensure an identifier is given
	ident, ok := origList[1].(*ast.Ident)
	if !ok {
		return nil, badSyntaxErr
	}

	return &ast.SetExpr{
		Ident: ident,
		Value: origList[2],
	}, nil
}

func lambdaSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	badSyntaxErr := errors.New("lambda: bad syntax")
	origList := input.List
	if len(origList) < 3 {
		return nil, badSyntaxErr
	}

	// ensure the argument list consists of identifiers
	argList, ok := origList[1].(*ast.ListExpr)
	if !ok {
		return nil, badSyntaxErr
	}
	var args []*ast.Ident
	for _, arg := range argList.List {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			return nil, badSyntaxErr
		}
		args = append(args, ident)
	}

	return &ast.LambdaExpr{
		Args: args,
		Body: origList[2:],
	}, nil
}

func condSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	badSyntaxErr := errors.New("cond: bad syntax")
	elseSyntaxErr := errors.New("cond: bad syntax: 'else' clause must be last")
	origList := input.List
	if len(origList) < 1 {
		return nil, badSyntaxErr
	}

	var branchList []*ast.BranchExpr
	var last bool
	for _, branch := range origList[1:] {
		if last {
			return nil, elseSyntaxErr
		}
		// ensure each branch is a list with at least 2 sub-expressions
		list, ok := branch.(*ast.ListExpr)
		if !ok {
			return nil, badSyntaxErr
		}
		if len(list.List) < 2 {
			return nil, badSyntaxErr
		}

		var elseBranch bool
		condition := list.List[0]
		if ident, ok := list.List[0].(*ast.Ident); ok {
			if *ident.Name == "else" {
				last, elseBranch = true, true
				condition = nil
			}
		}

		branchList = append(branchList, &ast.BranchExpr{
			Else:      elseBranch,
			Condition: condition,
			Body:      list.List[1:],
		})
	}
	return &ast.CondExpr{
		List: branchList,
	}, nil
}

func quoteSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	badSyntaxErr := errors.New("quote: bad syntax")
	origList := input.List
	if len(origList) != 2 {
		return nil, badSyntaxErr
	}
	return &ast.Quote{Expr: origList[1]}, nil
}
