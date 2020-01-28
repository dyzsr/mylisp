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
	builtinLambda = BuiltinTransformer{name: "lambda", proc: lambdaSyntax}
	builtinCond   = BuiltinTransformer{name: "cond", proc: condSyntax}
)

func builtinTransformerMap() map[string]Transformer {
	return map[string]Transformer{
		"define": builtinDefine,
		"lambda": builtinLambda,
		"cond":   builtinCond,
	}
}

func defineSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	usageErr := errors.New("bad syntax, usage: '(define <id> <value>)'")
	origList := input.List
	if len(origList) != 3 {
		return nil, usageErr
	}
	// ensure an identifier is given
	ident, ok := origList[1].(*ast.Ident)
	if !ok {
		return nil, usageErr
	}

	value := origList[2]
	// ensure a value is given
	if _, ok := value.(*ast.DefineExpr); ok {
		return nil, usageErr
	}
	return &ast.DefineExpr{
		Ident: ident,
		Value: value,
	}, nil
}

func lambdaSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	usageErr := errors.New("bad syntax, usage: '(lambda (<id> ...) <body> ...)'")
	origList := input.List
	if len(origList) < 3 {
		return nil, usageErr
	}

	// ensure the argument list consists of identifiers
	argList, ok := origList[1].(*ast.ListExpr)
	if !ok {
		return nil, usageErr
	}
	var args []*ast.Ident
	for _, arg := range argList.List {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			return nil, usageErr
		}
		args = append(args, ident)
	}

	return &ast.LambdaExpr{
		Args: args,
		Body: origList[2:],
	}, nil
}

func condSyntax(scope *ast.Scope, input *ast.ListExpr) (ast.Expr, error) {
	usageErr := errors.New("bad syntax, usage: '(cond (<condition> <body> ...) ...)'")
	elseSyntaxErr := errors.New("bad syntax: 'else' clause must be last")
	origList := input.List
	if len(origList) < 1 {
		return nil, usageErr
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
			return nil, usageErr
		}
		if len(list.List) < 2 {
			return nil, usageErr
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
