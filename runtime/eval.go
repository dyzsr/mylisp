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
	for k, v := range builtinVariables() {
		scope.Insert(ast.SymbolMap(k), v)
	}
	return &Runtime{
		scope: scope,
	}
}

func (r *Runtime) Eval(input ast.Expr) (Value, error) {
	return eval(r.scope, input)
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
	case *ast.Quote:
		return evalQuote(scope, expr)
	case *ast.Ident:
		return evalIdent(scope, expr)
	case *ast.ListExpr:
		return evalListExpr(scope, expr)
	case *ast.DefineExpr:
		return evalDefineExpr(scope, expr)
	case *ast.SetExpr:
		return evalSetExpr(scope, expr)
	case *ast.LambdaExpr:
		return evalLambdaExpr(scope, expr)
	case *ast.CondExpr:
		return evalCondExpr(scope, expr)
	}
	return nil, errors.New("error: cannot eval input")
}

func evalQuote(scope *ast.Scope, quote *ast.Quote) (Value, error) {
	switch expr := quote.Expr.(type) {
	case *ast.BoolLit:
		return Bool(expr.Value), nil
	case *ast.IntLit:
		return Int(expr.Value), nil
	case *ast.Ident:
		return Symbol{symbolMap(*expr.Name)}, nil
	case *ast.ListExpr:
		var args []Value
		for _, expr := range expr.List {
			value, err := evalQuote(scope, &ast.Quote{Expr: expr})
			if err != nil {
				return nil, err
			}
			args = append(args, value)
		}
		return _list(args...)
	}
	return nil, errors.New("quote: bad value")
}

func evalIdent(scope *ast.Scope, ident *ast.Ident) (Value, error) {
	value, ok := scope.Lookup(ident.Name)
	if !ok {
		return nil, fmt.Errorf("%s: undefined", *ident.Name)
	}
	return value.(Value), nil
}

func evalListExpr(scope *ast.Scope, listExpr *ast.ListExpr) (Value, error) {
	if len(listExpr.List) == 0 {
		return nil, errors.New("missing procedure expression")
	}
	valueList := make([]Value, len(listExpr.List))
	for i, subExpr := range listExpr.List {
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
	return nil, errors.New("not a procedure")
}

func evalBuiltinProc(op *BuiltinProc, operands ...Value) (value Value, err error) {
	return op.proc(operands...)
}

func evalProc(proc *Proc, operands ...Value) (Value, error) {
	if len(proc.Args) != len(operands) {
		return nil, arityMismatchErr
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
	if proc, ok := value.(*Proc); ok && proc.name == nil {
		proc.name = defineExpr.Ident.Name
	}
	return nil, nil
}

func evalSetExpr(scope *ast.Scope, setExpr *ast.SetExpr) (Value, error) {
	value, err := eval(scope, setExpr.Value)
	if err != nil {
		return nil, err
	}
	if ok := scope.Assign(setExpr.Ident.Name, value); !ok {
		return nil, fmt.Errorf("%s: undefined", *setExpr.Ident.Name)
	}
	if proc, ok := value.(*Proc); ok && proc.name == nil {
		proc.name = setExpr.Ident.Name
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
	for _, branch := range condExpr.List {
		if !branch.Else {
			condValue, err := eval(scope, branch.Condition)
			if err != nil {
				return nil, err
			}
			boolean, ok := condValue.(Bool)
			if !ok {
				return nil, errors.New("cond: error condition type")
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
