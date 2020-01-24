package token

import (
	"reflect"
	"strings"
	"testing"
)

type testLexer struct {
	input  string
	result []Token
}

var (
	testData = []testLexer{
		{
			input:  "(define x 1)",
			result: []Token{LPAREN, IDENT, IDENT, INTEGER, RPAREN},
		},
		{
			input: "(define f (lambda (x) (+ x x)))",
			result: []Token{
				LPAREN, IDENT, IDENT, LPAREN, IDENT, LPAREN, IDENT, RPAREN,
				LPAREN, PLUS, IDENT, IDENT, RPAREN, RPAREN, RPAREN,
			},
		},
		{
			input: "(+ x (- y 1) (* a (/ b 2) (% c 2)))",
			result: []Token{
				LPAREN, PLUS, IDENT, LPAREN, MINUS, IDENT, INTEGER, RPAREN,
				LPAREN, ASTER, IDENT, LPAREN, SLASH, IDENT, INTEGER, RPAREN,
				LPAREN, PERCENT, IDENT, INTEGER, RPAREN, RPAREN, RPAREN,
			},
		},
		{
			input: "(&& (= x y) (! (< y 1)) (|| (<= z 9) (> a b c) (>= d e f)))",
			result: []Token{
				LPAREN, AND, LPAREN, EQ, IDENT, IDENT, RPAREN, LPAREN, NOT, LPAREN, LT, IDENT, INTEGER, RPAREN, RPAREN,
				LPAREN, OR, LPAREN, LTE, IDENT, INTEGER, RPAREN, LPAREN, GT, IDENT, IDENT, IDENT, RPAREN,
				LPAREN, GTE, IDENT, IDENT, IDENT, RPAREN, RPAREN, RPAREN,
			},
		},
		{
			input: "`(a b `(d `x `y) `(123 456 `789))",
			result: []Token{
				QUOTE, LPAREN, IDENT, IDENT, QUOTE, LPAREN, IDENT, QUOTE, IDENT, QUOTE, IDENT, RPAREN,
				QUOTE, LPAREN, INTEGER, INTEGER, QUOTE, INTEGER, RPAREN, RPAREN,
			},
		},
	}
)

func TestNextToken(t *testing.T) {
	for _, test := range testData {
		result := testNextToken(test.input)

		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: '%v'\nexpect: %+v\noutput: %+v", test.input, test.result, result)
		}
	}
}

func testNextToken(input string) []Token {
	r := strings.NewReader(input)
	l := NewLexer(r)

	var result []Token
	for t, _ := l.LookupOne(); t != EOF; t, _ = l.LookupOne() {
		tok, _ := l.Next()
		print(t.String(), " ")
		result = append(result, tok)
	}
	println()
	return result
}
