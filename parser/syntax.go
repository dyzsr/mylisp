package parser

import (
	"errors"
	"mylisp/ast"
)

type syntaxFunc func(ast.Expr) (ast.Expr, error)

var (
	syntaxTable map[string]syntaxFunc
)

func init() {
	syntaxTable = map[string]syntaxFunc{
		"define": defineSyntax,
		"lambda": lambdaSyntax,
		"cond":   condSyntax,
	}
}

func parse(inputExpr ast.Expr) (ast.Expr, error) {
	// fmt.Println("Parser next")
	// atomic expression
	expr, ok := inputExpr.(*ast.ListExpr)
	if !ok {
		return inputExpr, nil
	}

	// list expression
	if len(expr.SubExprList) == 0 {
		return nil, errors.New("bad syntax: empty list")
	}
	first := expr.SubExprList[0]
	if ident, ok := first.(*ast.Ident); ok {
		if F := lookup(ident.Name); F != nil {
			// syntax block
			result, err := F(expr)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
	// procedure calls
	for i, sub := range expr.SubExprList {
		sub, err := parse(sub)
		if err != nil {
			return nil, err
		}
		expr.SubExprList[i] = sub
	}
	return expr, nil
}

func lookup(name string) syntaxFunc {
	if syntax, ok := syntaxTable[name]; ok {
		return syntax
	}
	return nil
}

func defineSyntax(expr ast.Expr) (ast.Expr, error) {
	list, ok := expr.(*ast.ListExpr)
	usageErr := errors.New("bad syntax, usage: '(define <id> <value>)'")
	if !ok {
		return nil, usageErr
	}
	subexpr := list.SubExprList
	if len(subexpr) != 3 {
		return nil, usageErr
	}
	ident, ok := subexpr[1].(*ast.Ident)
	if !ok {
		return nil, usageErr
	}

	// parse recursively
	value, usageErr := parse(subexpr[2])
	if usageErr != nil {
		return nil, usageErr
	}
	return &ast.DefineExpr{
		Ident: ident,
		Value: value,
	}, nil
}

func lambdaSyntax(expr ast.Expr) (ast.Expr, error) {
	list, ok := expr.(*ast.ListExpr)
	usageErr := errors.New("bad syntax, usage: '(lambda (<id> ...) <body> ...)'")
	if !ok {
		return nil, usageErr
	}
	subexpr := list.SubExprList
	if len(subexpr) < 3 {
		return nil, usageErr
	}
	argList, ok := subexpr[1].(*ast.ListExpr)
	if !ok {
		return nil, usageErr
	}
	var args []*ast.Ident
	for _, arg := range argList.SubExprList {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			return nil, usageErr
		}
		args = append(args, ident)
	}
	var body []ast.Expr
	for _, sub := range subexpr[2:] {
		sub, err := parse(sub)
		if err != nil {
			return nil, err
		}
		body = append(body, sub)
	}
	return &ast.LambdaExpr{
		Args: args,
		Body: body,
	}, nil
}

func condSyntax(expr ast.Expr) (ast.Expr, error) {
	list, ok := expr.(*ast.ListExpr)
	usageErr := errors.New("bad syntax, usage: '(cond (<condition> <body> ...) ...)'")
	if !ok {
		return nil, usageErr
	}
	subexpr := list.SubExprList
	if len(subexpr) < 1 {
		return nil, usageErr
	}

	var branchList []*ast.BranchExpr
	for _, branch := range subexpr[1:] {
		list, ok := branch.(*ast.ListExpr)
		if !ok {
			return nil, usageErr
		}
		if len(list.SubExprList) < 2 {
			return nil, usageErr
		}

		var elseBranch bool
		var condition ast.Expr
		if ident, ok := list.SubExprList[0].(*ast.Ident); ok {
			if ident.Name == "else" {
				elseBranch = true
			}
		} else {
			var err error
			condition, err = parse(list.SubExprList[0])
			if err != nil {
				return nil, err
			}
		}
		var body []ast.Expr
		for _, sub := range list.SubExprList[1:] {
			sub, err := parse(sub)
			if err != nil {
				return nil, err
			}
			body = append(body, sub)
		}

		branchList = append(branchList, &ast.BranchExpr{
			Else:      elseBranch,
			Condition: condition,
			Body:      body,
		})
	}
	return &ast.CondExpr{
		BranchList: branchList,
	}, nil
}
