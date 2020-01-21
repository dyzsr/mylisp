package main

import (
	"os"

	"github.com/dyzsr/mylisp/parser"
	"github.com/dyzsr/mylisp/repl"
	"github.com/dyzsr/mylisp/runtime"
	"github.com/dyzsr/mylisp/token"
)

func main() {
	println("start")
	lex := token.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
	env := runtime.NewRootEnv()
	prt := repl.NewPrinter()
	eprt := repl.NewErrorPrinter()

	for {
		println("parse")
		expr, ok := par.Next()
		if !ok && par.Err() == nil { // EOF
			break
		}
		println("eval")
		result, err := env.Eval(expr)
		println("print")
		if err != nil {
			eprt.Print(err)
		} else {
			prt.Print(result)
		}
	}
}
