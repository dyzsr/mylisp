package runtime

import (
	"errors"
	"fmt"
	"mylisp/ast"
)

func (e *EvalEnv) Eval(inputExpr ast.Expr) (ast.Expr, error) {
	if inputExpr == nil {
		return nil, nil
	}

	switch expr := inputExpr.(type) {
	case *BoolValue:
		return expr, nil
	case *IntValue:
		return expr, nil
	case *BuiltinProc:
		return expr, nil
	case *ProcValue:
		return expr, nil

	case *ast.Ident:
		return e.evalIdent(expr)
	case *ast.ListExpr:
		return e.evalListExpr(expr)
	case *ast.DefineExpr:
		return e.evalDefineExpr(expr)
	case *ast.LambdaExpr:
		return e.evalLambdaExpr(expr)
	case *ast.CondExpr:
		return e.evalCondExpr(expr)
	}
	return nil, errors.New("eval error")
}

func (e *EvalEnv) evalIdent(ident *ast.Ident) (ast.Expr, error) {
	value, ok := e.lookup(ident.Name)
	if !ok {
		return nil, fmt.Errorf("identifier '%s' is not bound for any value", ident.Name)
	}
	return value, nil
}

func (e *EvalEnv) evalListExpr(listExpr *ast.ListExpr) (ast.Expr, error) {
	subExprList := listExpr.SubExprList
	if len(subExprList) == 0 {
		return nil, errors.New("invalid expression list: at least one sub expression is required")
	}
	for i, subExpr := range subExprList {
		var err error
		subExprList[i], err = e.Eval(subExpr)
		if err != nil {
			return nil, err
		}
	}

	operator := subExprList[0]
	operands := subExprList[1:]

	switch op := operator.(type) {
	case *BuiltinProc:
		return e.evalBuiltinProc(op.Name, operands...)
	case *ProcValue:
		return e.evalProc(op, operands...)
	}
	return nil, errors.New("invalid operator: either an identifier, basic procedure, or lambda expression is expected")
}

func (e *EvalEnv) evalDefineExpr(defineExpr *ast.DefineExpr) (ast.Expr, error) {
	value, err := e.Eval(defineExpr.Value)
	if err != nil {
		return nil, err
	}
	e.insert(defineExpr.Ident.Name, value)
	if proc, ok := value.(*ProcValue); ok && len(proc.Name) == 0 {
		proc.Name = defineExpr.Ident.Name
	}
	return nil, nil
}

func (e *EvalEnv) evalLambdaExpr(lambdaExpr *ast.LambdaExpr) (ast.Expr, error) {
	return &ProcValue{
		LambdaExpr: lambdaExpr,
		env:        e,
	}, nil
}

func (e *EvalEnv) evalCondExpr(condExpr *ast.CondExpr) (ast.Expr, error) {
	for _, branch := range condExpr.BranchList {
		if !branch.Else {
			condValue, err := e.Eval(branch.Condition)
			if err != nil {
				return nil, err
			}
			boolean, ok := condValue.(*BoolValue)
			if !ok {
				return nil, errors.New("branch condition should be a boolean expression")
			}
			if !boolean.Value {
				continue
			}
		}

		var result ast.Expr
		env := NewEvalEnv(e)
		for _, expr := range branch.Body {
			var err error
			result, err = env.Eval(expr)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, nil
}
