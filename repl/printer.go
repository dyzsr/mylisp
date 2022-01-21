package repl

import (
	"fmt"

	"github.com/dyzsr/mylisp/runtime"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(v runtime.Value) {
	if v == nil {
		return
	} else if _, ok := v.(runtime.Nil); ok {
		return
	}
	fmt.Printf("$ %s\n", v)
}
