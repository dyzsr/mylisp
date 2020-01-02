package main

import (
	"mylisp/parser"
	"mylisp/printer"
	"mylisp/runtime"
)

func main() {
	par := parser.NewParser()
	env := runtime.NewRootEnv()
	prt := printer.NewPrinter()
	eprt := printer.NewErrorPrinter()

	for {
		expr := par.NextExpr()
		result, err := env.Eval(expr)
		if err != nil {
			eprt.Print(err)
		} else {
			prt.Print(result)
		}
	}
}
