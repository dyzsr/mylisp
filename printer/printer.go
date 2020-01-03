package printer

import (
	"fmt"
	"mylisp/ast"
	"mylisp/runtime"
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
	case *runtime.BoolValue:
		fmt.Printf("$ %v\n", v.Value)
	case *runtime.IntValue:
		fmt.Printf("$ %v\n", v.Value)
	case *runtime.BuiltinProc:
		fmt.Printf("$ <procedure %s>", v.Name)
	case *runtime.ProcValue:
		fmt.Printf("$ <procedure %s>", v.Name)
	}
}
