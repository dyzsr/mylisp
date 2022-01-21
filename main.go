package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dyzsr/mylisp/compiletime"
	"github.com/dyzsr/mylisp/parser"
	"github.com/dyzsr/mylisp/repl"
	"github.com/dyzsr/mylisp/runtime"
	"github.com/dyzsr/mylisp/token"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for {
			sig := <-sigs
			println("\n~", sig.String())
		}
	}()

	lex := token.NewLexer(os.Stdin)
	par := parser.NewParser(lex)
	ct := compiletime.NewCompileTime()
	rt := runtime.NewRuntime()
	prt := repl.NewPrinter()
	eprt := repl.NewErrorPrinter()

	for {
		expr, ok := par.Next()
		if !ok {
			err := par.Err()
			if err == nil { // EOF
				break
			}
			eprt.Print(err)
			continue
		}

		expr, err := ct.Eval(expr)
		if err != nil {
			eprt.Print(err)
			continue
		}

		result, err := rt.Eval(expr)

		if err != nil {
			eprt.Print(err)
		} else {
			prt.Print(result)
		}
	}
}
