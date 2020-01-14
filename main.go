package main

import (
	"github.com/dyzsr/mylisp/parser"
	"github.com/dyzsr/mylisp/repl"
	"github.com/dyzsr/mylisp/runtime"
	"github.com/dyzsr/mylisp/token"
	"os"
)

func main() {
	lex := token.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
	env := runtime.NewRootEvalEnv()
	prt := repl.NewPrinter()
	eprt := repl.NewErrorPrinter()

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
