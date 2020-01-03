package lexer

import (
	"mylisp/token"
	"reflect"
	"strings"
	"testing"
)

type testLexer struct {
	input  string
	result []token.Token
}

var (
	testData = []testLexer{
		{
			input:  "(define x 1)",
			result: []token.Token{token.DEFINE, token.IDENT, token.INTEGER},
		},
		{
			input: "(+ x (- y 1) (* a (/ b 2)))",
			result: []token.Token{
				token.PLUS, token.IDENT, token.LPAREN, token.IDENT, token.INTEGER, token.RPAREN,
				token.LPAREN, token.ASTER, token.IDENT, token.LPAREN, token.SLASH, token.IDENT,
				token.INTEGER, token.RPAREN, token.RPAREN, token.RPAREN,
			},
		},
	}
)

func TestNextToken(t *testing.T) {
	for _, test := range testData {
		result := testNextToken(test.input)

		t.Logf("input: '%v', output: %+v", test.input, result)
		if reflect.DeepEqual(result, test.result) {
			t.Errorf("expect: %+v, output: %+v", test.result, result)
		}
	}
}

func testNextToken(input string) []token.Token {
	r := strings.NewReader(input)
	l := NewLexer(r)

	var result []token.Token
	for tok := l.NextToken(); tok != token.EOF; tok = l.NextToken() {
		result = append(result, tok)
	}
	return result
}
