package runtime

import (
	"errors"
	"fmt"

	"github.com/dyzsr/mylisp/ast"
)

type Runtime struct {
	scope       *ast.Scope // variable binding scope
	stack       *callstack // procedure calls stack
	lastInScope bool       // is the last expression in a scope
	enableTCOpt bool       // enable tail call optimization
}

func NewRuntime() *Runtime {
	scope := ast.NewRootScope()
	for k, v := range builtinVariables() {
		scope.Insert(ast.SymbolMap(k), v)
	}
	return &Runtime{
		scope:       scope,
		stack:       newCallstack(),
		enableTCOpt: true,
	}
}

func (r *Runtime) Eval(input ast.Expr) (Value, error) {
	return r.eval(r.scope, input)
}

func (r *Runtime) eval(scope *ast.Scope, input ast.Expr) (Value, error) {
	if input == nil {
		return nil, nil
	}

	switch expr := input.(type) {
	case *ast.BoolLit:
		return Bool(expr.Value), nil
	case *ast.IntLit:
		return Int(expr.Value), nil
	case *ast.Quote:
		return r.evalQuote(scope, expr)
	case *ast.Ident:
		return r.evalIdent(scope, expr)
	case *ast.ListExpr:
		return r.evalListExpr(scope, expr)
	case *ast.DefineExpr:
		return r.evalDefineExpr(scope, expr)
	case *ast.SetExpr:
		return r.evalSetExpr(scope, expr)
	case *ast.LambdaExpr:
		return r.evalLambdaExpr(scope, expr)
	case *ast.CondExpr:
		return r.evalCondExpr(scope, expr)
	}
	return nil, errors.New("error: cannot eval input")
}

func (r *Runtime) evalQuote(scope *ast.Scope, quote *ast.Quote) (Value, error) {
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
			value, err := r.evalQuote(scope, &ast.Quote{Expr: expr})
			if err != nil {
				return nil, err
			}
			args = append(args, value)
		}
		return _list(args...)
	}
	return nil, errors.New("quote: bad value")
}

func (r *Runtime) evalIdent(scope *ast.Scope, ident *ast.Ident) (Value, error) {
	value, ok := scope.Lookup(ident.Name)
	if !ok {
		return nil, fmt.Errorf("%s: undefined", *ident.Name)
	}
	return value.(Value), nil
}

func (r *Runtime) evalListExpr(scope *ast.Scope, listExpr *ast.ListExpr) (Value, error) {
	if len(listExpr.List) == 0 {
		return nil, errors.New("missing procedure expression")
	}

	var tailcall bool
	if r.enableTCOpt && !r.stack.empty() && r.stack.last() && r.lastInScope {
		// mark tailcall when requirements are satisfied
		tailcall = true
		// turn off the flag since a list expression (procedure call) is not a scope
		r.lastInScope = false
	}

	valueList := make([]Value, len(listExpr.List))
	for i, expr := range listExpr.List {
		var err error
		valueList[i], err = r.eval(scope, expr)
		if err != nil {
			return nil, err
		}
	}

	operator := valueList[0]
	operands := valueList[1:]

	switch op := operator.(type) {
	case *BuiltinProc:
		return r.evalBuiltinProc(op, operands...)
	case *Proc:
		if tailcall {
			r.stack.modify(op, operands)
			return nil, nil
		}
		r.stack.push(op, operands)
		result, err := r.evalProc(op, operands...)
		r.stack.pop()
		return result, err
	}
	return nil, errors.New("not a procedure")
}

func (r *Runtime) evalDefineExpr(scope *ast.Scope, defineExpr *ast.DefineExpr) (Value, error) {
	value, err := r.eval(scope, defineExpr.Value)
	if err != nil {
		return nil, err
	}
	scope.Insert(defineExpr.Ident.Name, value)
	if proc, ok := value.(*Proc); ok && proc.name == nil {
		proc.name = defineExpr.Ident.Name
	}
	return Nil{}, nil
}

func (r *Runtime) evalSetExpr(scope *ast.Scope, setExpr *ast.SetExpr) (Value, error) {
	value, err := r.eval(scope, setExpr.Value)
	if err != nil {
		return nil, err
	}
	if ok := scope.Assign(setExpr.Ident.Name, value); !ok {
		return nil, fmt.Errorf("%s: undefined", *setExpr.Ident.Name)
	}
	if proc, ok := value.(*Proc); ok && proc.name == nil {
		proc.name = setExpr.Ident.Name
	}
	return Nil{}, nil
}

func (r *Runtime) evalLambdaExpr(scope *ast.Scope, lambdaExpr *ast.LambdaExpr) (Value, error) {
	return &Proc{
		LambdaExpr: lambdaExpr,
		outer:      scope,
	}, nil
}

func (r *Runtime) evalCondExpr(scope *ast.Scope, condExpr *ast.CondExpr) (Value, error) {
	for _, branch := range condExpr.List {
		if !branch.Else {
			condValue, err := r.eval(scope, branch.Condition)
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

		inner := ast.NewScope(scope)
		return r.evalScope(inner, false, branch.Body)
	}
	return Nil{}, nil
}

func (r *Runtime) evalBuiltinProc(op *BuiltinProc, operands ...Value) (value Value, err error) {
	return op.proc(operands...)
}

func (r *Runtime) evalProc(proc *Proc, operands ...Value) (Value, error) {
	// println("callstack depth:", len(r.stack.frames))
	for {
		if r.stack.modified() {
			topFrame := r.stack.top()
			proc, operands = topFrame.proc, topFrame.params
			r.stack.unmodify()
		}

		if len(proc.Args) != len(operands) {
			return nil, errArityMismatch
		}
		scope := ast.NewScope(proc.outer)
		for i, arg := range proc.Args {
			scope.Insert(arg.Name, operands[i])
		}

		result, err := r.evalScope(scope, true, proc.Body)
		if r.stack.modified() {
			continue
		}
		return result, err
	}
}

func (r *Runtime) evalScope(scope *ast.Scope, isProc bool, list []ast.Expr) (Value, error) {
	// turn off the flag when entering a new scope
	r.lastInScope = false

	var result Value
	for i, expr := range list {
		if i == len(list)-1 {
			// turn on the flag if it is last in the scope
			r.lastInScope = true
			if isProc {
				r.stack.setLast(true)
			}
		}

		var err error
		result, err = r.eval(scope, expr)

		// restore the flag after being used
		r.lastInScope = false

		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
