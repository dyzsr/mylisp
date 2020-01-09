package printer

import (
	"fmt"
	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/runtime"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(e ast.Expr) {
	if e == nil {
		return
	}

	switch v := e.(type) {
	case runtime.BoolValue:
		fmt.Printf("$ %v\n", v)
	case runtime.IntValue:
		fmt.Printf("$ %v\n", v)
	case *runtime.BuiltinProc:
		fmt.Printf("$ <procedure %s>\n", v.Name)
	case *runtime.ProcValue:
		if len(v.Name) > 0 {
			fmt.Printf("$ <procedure %s at %p>\n", v.Name, v.LambdaExpr)
		} else {
			fmt.Printf("$ <procedure at %p>\n", v.LambdaExpr)
		}
	}
}
