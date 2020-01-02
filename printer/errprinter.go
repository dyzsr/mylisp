package printer

import "fmt"

import "os"

type ErrorPrinter struct {
}

func NewErrorPrinter() *ErrorPrinter {
	return &ErrorPrinter{}
}

func (p *ErrorPrinter) Print(err error) {
	fmt.Fprintf(os.Stderr, "eval error: %s", err)
}
