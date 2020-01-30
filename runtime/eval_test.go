package runtime

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dyzsr/mylisp/ast"
	"github.com/dyzsr/mylisp/compiletime"
	"github.com/dyzsr/mylisp/parser"
	"github.com/dyzsr/mylisp/token"
)

type testStruct struct {
	str    string
	input  ast.Expr
	result Value
}

func makeTest(testData []testStruct) func(*testing.T) {
	return func(t *testing.T) {
		ct := compiletime.NewCompileTime()
		rt := NewRuntime()
		for _, test := range testData {
			if len(test.str) > 0 {
				r := strings.NewReader(test.str)
				p := parser.NewParser(token.NewLexer(r))

				expr, _ := p.Next()
				test.input, _ = ct.Eval(expr)
			}

			result, err := rt.Eval(test.input)
			if err != nil {
				t.Error(err)
				break
			}
			if !reflect.DeepEqual(result, test.result) {
				t.Errorf("\ninput: '%s'\nexpect: '%s'\noutput: '%s'", test.input, test.result, result)
			}
			t.Log(test.input)
		}
		if !rt.stack.empty() {
			t.Error("callstack is not empty")
		}
	}
}

func Test_Eval(t *testing.T) {
	t.Run("DefineExpr", makeTest(testDefineExpr))
	t.Run("SetExpr", makeTest(testSetExpr))
	t.Run("BuiltinProc", makeTest(testListExpr))
	t.Run("LambdaExpr", makeTest(testLambdaExpr))
	t.Run("CondExpr", makeTest(testCondExpr))
	t.Run("QuoteExpr", makeTest(testQuote))
	t.Run("TailCall", makeTest(testTailCall))
	t.Run("ChurchNumeral", makeTest(testChurchNumeral))
	t.Run("MessagePassing", makeTest(testMessagePassing))
	t.Run("Fibonacci", makeTest(testFibonacci))
}

var (
	testDefineExpr = []testStruct{
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: 32768},
			},
			result: nil,
		},
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("y"),
				Value: &ast.IntLit{Value: 123},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(32768),
		},
		{
			input:  ast.NewIdent("y"),
			result: Int(123),
		},
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: -32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(-32768),
		},
	}

	testSetExpr = []testStruct{
		{
			input: &ast.DefineExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: 32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(32768),
		},
		{
			input: &ast.SetExpr{
				Ident: ast.NewIdent("x"),
				Value: &ast.IntLit{Value: -32768},
			},
			result: nil,
		},
		{
			input:  ast.NewIdent("x"),
			result: Int(-32768),
		},
	}

	testListExpr = []testStruct{
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("+"),
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: 32768},
				}},
			result: Int(65535),
		},
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("+"),
					&ast.IntLit{Value: 32767},
					&ast.IntLit{Value: -32768},
				}},
			result: Int(-1),
		},
	}

	testLambdaExpr = []testStruct{
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					&ast.LambdaExpr{
						Args: []*ast.Ident{ast.NewIdent("x")},
						Body: []ast.Expr{
							ast.NewIdent("x"),
						}},
					&ast.IntLit{Value: 123},
				}},
			result: Int(123),
		},
		{
			input: &ast.ListExpr{
				List: []ast.Expr{
					&ast.ListExpr{
						List: []ast.Expr{
							&ast.LambdaExpr{
								Args: []*ast.Ident{ast.NewIdent("x")},
								Body: []ast.Expr{
									&ast.LambdaExpr{
										Args: []*ast.Ident{ast.NewIdent("y")},
										Body: []ast.Expr{
											&ast.ListExpr{
												List: []ast.Expr{
													ast.NewIdent("+"),
													ast.NewIdent("y"),
													ast.NewIdent("x"),
												}}}}}},
							&ast.IntLit{Value: 32767}}},
					&ast.IntLit{Value: 32768},
				}},
			result: Int(65535),
		},
	}

	testCondExpr = []testStruct{
		{
			input: &ast.CondExpr{
				List: []*ast.BranchExpr{
					{
						Condition: &ast.BoolLit{Value: false},
						Body: []ast.Expr{
							&ast.IntLit{Value: 123},
						},
					},
					{
						Else: true,
						Body: []ast.Expr{
							&ast.IntLit{Value: 654},
						},
					},
				},
			},
			result: Int(654),
		},
		{
			input: &ast.CondExpr{
				List: []*ast.BranchExpr{
					{
						Condition: &ast.BoolLit{Value: false},
						Body: []ast.Expr{
							&ast.IntLit{Value: 123},
						},
					},
					{
						Condition: &ast.BoolLit{Value: true},
						Body: []ast.Expr{
							&ast.IntLit{Value: 456},
						},
					},
				},
			},
			result: Int(456),
		},
	}

	testQuote = []testStruct{
		{
			input:  &ast.Quote{Expr: ast.NewIdent("abc")},
			result: Symbol{symbolMap("abc")},
		},
		{
			input: &ast.Quote{Expr: &ast.ListExpr{
				List: []ast.Expr{
					ast.NewIdent("abc"),
					&ast.ListExpr{
						List: []ast.Expr{
							ast.NewIdent("quote"),
							&ast.IntLit{Value: 567},
						},
					},
				},
			}},
			result: &Pair{
				first: Symbol{symbolMap("abc")},
				second: &Pair{
					first: &Pair{
						first:  Symbol{symbolMap("quote")},
						second: &Pair{first: Int(567), second: Nil{}},
					},
					second: Nil{},
				},
			},
		},
	}

	testTailCall = []testStruct{
		{
			str: `
				(define fac
				 (lambda (x)
				  (define iter
				   (lambda (i res)
				    (cond ((= i 0) res)
				          (else (iter (- i 1) (* i res))))))
				  (iter x 1)))
				 `,
		},
		{str: "(fac 5)", result: Int(120)},
		{str: "(fac 15)", result: Int(1307674368000)},

		{
			str: `
				(define sum
				 (lambda (n)
				  (define iter
				   (lambda (i res)
				    (cond ((= i 0) res)
				          (else (iter (- i 1) (+ res i))))))
				  (iter n 0)))
			`,
		},
		{str: "(sum 10)", result: Int(55)},
		{str: "(sum 10000)", result: Int(50005000)},
	}

	testChurchNumeral = []testStruct{
		{str: "(define n0 (lambda (f) (lambda (x) x)))"},
		{str: "(define n1 (lambda (f) (lambda (x) (f x))))"},
		{str: "(define show (lambda (n) ((n (lambda (x) (+ x 1))) 0)))"},
		{str: "(define add (lambda (a b) (lambda (f) (lambda (x) ((a f) ((b f) x))))))"},
		{str: "(define mul (lambda (a b) (lambda (f) (lambda (x) ((a (b f)) x)))))"},
		{str: "(define n2 (add n1 n1))"},
		{str: "(define n3 (add n1 n2))"},
		{str: "(define n4 (add n2 n2))"},
		{str: "(define n5 (add n2 n3))"},
		{str: "(define n8 (add n3 n5))"},
		{str: "(define n13 (add n5 n8))"},
		{str: "(define n65 (mul n5 n13))"},
		{str: "(define n32 (mul n4 n8))"},
		{str: "(define n64 (mul n8 n8))"},
		{str: "(define n1024 (mul n32 n32))"},
		{str: "(show n64)", result: Int(64)},
		{str: "(show n1024)", result: Int(1024)},
		{str: "(show n65)", result: Int(65)},
	}

	testMessagePassing = []testStruct{
		{str: `(define NewProfile
			(lambda ()
			 (define id 0)
			 (define name 'name)
			 (define setId (lambda (x) (set! id x)))
			 (define setName (lambda (x) (set! name x)))
			 (lambda (msg)
			  (cond ((eq? msg 'Id) id)
					((eq? msg 'SetId) setId)
					((eq? msg 'Name) name)
					((eq? msg 'SetName) setName)))))`},
		{str: "(define p (NewProfile))"},
		{str: "(p 'Id)", result: Int(0)},
		{str: "((p 'SetId) 1)"},
		{str: "(p 'Id)", result: Int(1)},
		{str: "(p 'Name)", result: Symbol{symbolMap("name")}},
		{str: "((p 'SetName) 'dy)"},
		{str: "(p 'Name)", result: Symbol{symbolMap("dy")}},
	}

	testFibonacci = []testStruct{
		{str: `(define fib (lambda (n)
				(cond ((< n 2) 1)
					(else (+ (fib (- n 1)) (fib (- n 2)))))))`},
		{str: "(fib 10)", result: Int(89)},
		{str: "(fib 15)", result: Int(987)},
		{str: "(fib 20)", result: Int(10946)},
	}
)
