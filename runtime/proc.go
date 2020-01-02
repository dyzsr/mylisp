package runtime

import (
	"mylisp/ast"
)

type ProcValue struct {
	*ast.LambdaExpr
	Name string
	env  *EvalEnv
}

func (e *EvalEnv) evalProc(proc *ProcValue, operands ...ast.Expr) (ast.Expr, error) {
	env := NewEvalEnv(proc.env)
	for i, arg := range proc.Args {
		env.insert(arg.Name, operands[i])
	}
	var result ast.Expr
	for _, expr := range proc.Body {
		var err error
		result, err = env.Eval(expr)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
