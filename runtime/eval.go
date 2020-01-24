package runtime

import (
	"errors"
	"fmt"

	"github.com/dyzsr/mylisp/ast"
)

type Runtime struct {
	scope *ast.Scope
}

func NewRuntime() *Runtime {
	scope := ast.NewRootScope()
	for k, v := range builtinProcMap() {
		scope.Insert(k, v)
	}
	return &Runtime{
		scope: scope,
	}
}

func (e *Runtime) Lookup(name string) (Value, bool) {
	if value, ok := e.scope.Lookup(name); ok {
		return value.(Value), true
	}
	return nil, false
}

func (e *Runtime) Insert(name string, value Value) {
	e.scope.Insert(name, value)
}

func (e *Runtime) Eval(input ast.Expr) (Value, error) {
	return eval(e.scope, input)
}

func eval(scope *ast.Scope, input ast.Expr) (Value, error) {
	if input == nil {
		return nil, nil
	}

	switch expr := input.(type) {
	case *ast.BoolLit:
		return Bool(expr.Value), nil
	case *ast.IntLit:
		return Int(expr.Value), nil
	case *ast.Ident:
		return evalIdent(scope, expr)
	case *ast.ListExpr:
		return evalListExpr(scope, expr)
	case *ast.DefineExpr:
		return evalDefineExpr(scope, expr)
	case *ast.LambdaExpr:
		return evalLambdaExpr(scope, expr)
	case *ast.CondExpr:
		return evalCondExpr(scope, expr)
	}
	return nil, errors.New("eval error: unknown node type")
}

func evalIdent(scope *ast.Scope, ident *ast.Ident) (Value, error) {
	value, ok := scope.Lookup(ident.Name)
	if !ok {
		return nil, fmt.Errorf("identifier '%s' is not bound for any value", ident.Name)
	}
	return value.(Value), nil
}

func evalListExpr(scope *ast.Scope, listExpr *ast.ListExpr) (Value, error) {
	if len(listExpr.SubExprList) == 0 {
		return nil, errors.New("invalid expression list: at least one sub expression is required")
	}
	valueList := make([]Value, len(listExpr.SubExprList))
	for i, subExpr := range listExpr.SubExprList {
		var err error
		valueList[i], err = eval(scope, subExpr)
		if err != nil {
			return nil, err
		}
	}

	operator := valueList[0]
	operands := valueList[1:]

	switch op := operator.(type) {
	case *BuiltinProc:
		return evalBuiltinProc(op, operands...)
	case *Proc:
		return evalProc(op, operands...)
	}
	return nil, errors.New("invalid operator: either an identifier, basic procedure, or lambda expression is expected")
}

func evalBuiltinProc(op *BuiltinProc, operands ...Value) (value Value, err error) {
	return op.proc(operands...)
}

func evalProc(proc *Proc, operands ...Value) (Value, error) {
	if len(proc.Args) != len(operands) {
		return nil, errors.New("operands and arguments numbers did not match")
	}
	scope := ast.NewScope(proc.outer)
	for i, arg := range proc.Args {
		scope.Insert(arg.Name, operands[i])
	}
	var result Value
	for _, expr := range proc.Body {
		var err error
		result, err = eval(scope, expr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalDefineExpr(scope *ast.Scope, defineExpr *ast.DefineExpr) (Value, error) {
	value, err := eval(scope, defineExpr.Value)
	if err != nil {
		return nil, err
	}
	scope.Insert(defineExpr.Ident.Name, value)
	if proc, ok := value.(*Proc); ok && len(proc.name) == 0 {
		proc.name = defineExpr.Ident.Name
	}
	return nil, nil
}

func evalLambdaExpr(scope *ast.Scope, lambdaExpr *ast.LambdaExpr) (Value, error) {
	return &Proc{
		LambdaExpr: lambdaExpr,
		outer:      scope,
	}, nil
}

func evalCondExpr(scope *ast.Scope, condExpr *ast.CondExpr) (Value, error) {
	for _, branch := range condExpr.BranchList {
		if !branch.Else {
			condValue, err := eval(scope, branch.Condition)
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
		inner := ast.NewScope(scope)
		for _, expr := range branch.Body {
			var err error
			result, err = eval(inner, expr)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, nil
}
