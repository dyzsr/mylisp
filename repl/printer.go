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
		// println("nil")
		return
	}
	fmt.Printf("$ %s\n", v)
}
