package main

import (
	"mylisp/parser"
	"mylisp/printer"
	"mylisp/runtime"
	"mylisp/token"
	"os"
)

func main() {
	lex := token.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
	env := runtime.NewRootEnv()
	prt := printer.NewPrinter()
	eprt := printer.NewErrorPrinter()

	for {
		expr, ok := par.Next()
		if !ok && par.Err() == nil { // EOF
			break
		}
		result, err := env.Eval(expr)
		if err != nil {
			eprt.Print(err)
		} else {
			prt.Print(result)
		}
	}
}
