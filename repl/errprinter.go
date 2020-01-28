package repl

import (
	"fmt"
	"os"
)

type ErrorPrinter struct{}

func NewErrorPrinter() *ErrorPrinter {
	return &ErrorPrinter{}
}

func (p *ErrorPrinter) Print(err error) {
	fmt.Fprintf(os.Stderr, "! %s\n", err)
}
