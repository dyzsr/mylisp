package main

import (
	"mylisp/lexer"
	"mylisp/parser"
	"mylisp/printer"
	"mylisp/runtime"
	"os"
)

func main() {
	lex := lexer.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
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
