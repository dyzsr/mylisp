package runtime

import (
	"errors"
	"fmt"

	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/env"
)

type EvalEnv struct {
	*env.Scope
}

func NewRootEvalEnv() *EvalEnv {
	scope := env.NewRootScope()
	for k, v := range defaultSymbols {
		scope.Insert(k, v)
	}
	return &EvalEnv{
		Scope: scope,
	}
}

func NewEvalEnv(parent *EvalEnv) *EvalEnv {
	return &EvalEnv{
		Scope: env.NewScope(parent.Scope),
	}
}

func (e *EvalEnv) Lookup(name string) (Value, bool) {
	if value, ok := e.Scope.Lookup(name); ok {
		return value.(Value), true
	}
	return nil, false
}

func (e *EvalEnv) Insert(name string, value Value) {
	e.Scope.Insert(name, value)
}

func (e *EvalEnv) Eval(inputExpr ast.Expr) (Value, error) {
	if inputExpr == nil {
		return nil, nil
	}

	switch expr := inputExpr.(type) {
	case Value:
		return expr, nil

	case *ast.BoolLit:
		return Bool(expr.Value), nil
	case *ast.IntLit:
		return Int(expr.Value), nil
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

func (e *EvalEnv) evalIdent(ident *ast.Ident) (Value, error) {
	value, ok := e.Lookup(ident.Name)
	if !ok {
		return nil, fmt.Errorf("identifier '%s' is not bound for any value", ident.Name)
	}
	return value, nil
}

func (e *EvalEnv) evalListExpr(listExpr *ast.ListExpr) (Value, error) {
	if len(listExpr.SubExprList) == 0 {
		return nil, errors.New("invalid expression list: at least one sub expression is required")
	}
	valueList := make([]Value, len(listExpr.SubExprList))
	for i, subExpr := range listExpr.SubExprList {
		var err error
		valueList[i], err = e.Eval(subExpr)
		if err != nil {
			return nil, err
		}
	}

	operator := valueList[0]
	operands := valueList[1:]

	switch op := operator.(type) {
	case *BuiltinProc:
		return e.evalBuiltinProc(op.Name, operands...)
	case *Proc:
		return e.evalProc(op, operands...)
	}
	return nil, errors.New("invalid operator: either an identifier, basic procedure, or lambda expression is expected")
}

func (e *EvalEnv) evalProc(proc *Proc, operands ...Value) (Value, error) {
	if len(proc.Args) != len(operands) {
		return nil, errors.New("operands and arguments numbers did not match")
	}
	env := NewEvalEnv(proc.env)
	for i, arg := range proc.Args {
		env.Insert(arg.Name, operands[i])
	}
	var result Value
	for _, expr := range proc.Body {
		var err error
		result, err = env.Eval(expr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (e *EvalEnv) evalDefineExpr(defineExpr *ast.DefineExpr) (Value, error) {
	value, err := e.Eval(defineExpr.Value)
	if err != nil {
		return nil, err
	}
	e.Insert(defineExpr.Ident.Name, value)
	if proc, ok := value.(*Proc); ok && len(proc.Name) == 0 {
		proc.Name = defineExpr.Ident.Name
	}
	return nil, nil
}

func (e *EvalEnv) evalLambdaExpr(lambdaExpr *ast.LambdaExpr) (Value, error) {
	return &Proc{
		LambdaExpr: lambdaExpr,
		env:        e,
	}, nil
}

func (e *EvalEnv) evalCondExpr(condExpr *ast.CondExpr) (Value, error) {
	for _, branch := range condExpr.BranchList {
		if !branch.Else {
			condValue, err := e.Eval(branch.Condition)
			if err != nil {
				return nil, err
			}
			boolean, ok := condValue.(Bool)
			if !ok {
				return nil, errors.New("branch condition should be a boolean expression")
			}
			if !boolean {
				continue
			}
		}

		var result Value
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
