package repl

import (
	"fmt"

	"github.com/dyzsr/mylisp/runtime"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(e runtime.Value) {
	if e == nil {
		return
	}

	switch v := e.(type) {
	case runtime.Bool:
		fmt.Printf("$ %v\n", v)
	case runtime.Int:
		fmt.Printf("$ %v\n", v)
	case *runtime.BuiltinProc:
		fmt.Printf("$ <procedure %s>\n", v.Name)
	case *runtime.Proc:
		if len(v.Name) > 0 {
			fmt.Printf("$ <procedure %s at %p>\n", v.Name, v.LambdaExpr)
		} else {
			fmt.Printf("$ <procedure at %p>\n", v.LambdaExpr)
		}
	}
}
