package main

import (
	"os"

	"github.com/dyzsr/mylisp/compiletime"
	"github.com/dyzsr/mylisp/parser"
	"github.com/dyzsr/mylisp/repl"
	"github.com/dyzsr/mylisp/runtime"
	"github.com/dyzsr/mylisp/token"
)

func main() {
	// println("start")

	lex := token.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
	ct := compiletime.NewCompileTime()
	rt := runtime.NewRuntime()
	prt := repl.NewPrinter()
	eprt := repl.NewErrorPrinter()

	for {
		// println("parse")
		expr, ok := par.Next()
		if !ok {
			err := par.Err()
			if err == nil { // EOF
				break
			}
			eprt.Print(err)
			continue
		}

		// println("compile time")
		expr, err := ct.Eval(expr)
		if err != nil {
			eprt.Print(err)
			continue
		}

		// println("runtime")
		result, err := rt.Eval(expr)

		// println("print")
		if err != nil {
			eprt.Print(err)
		} else {
			prt.Print(result)
		}
	}
}
