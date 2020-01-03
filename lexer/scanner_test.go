package lexer

import (
	"strings"
	"testing"
)

type testScanner struct {
	input  string
	result string
}

func makeScannerTest(testData []testScanner, f func(string) string) func(*testing.T) {
	return func(t *testing.T) {
		for _, test := range testData {
			result := f(test.input)

			t.Logf("input: '%s', output: '%s'", test.input, result)
			if result != test.result {
				t.Errorf("expect: '%s', output: '%s'", test.result, result)
			}
		}
	}
}

func TestGetPeek(t *testing.T) {
	t.Run("get & peek", makeScannerTest(dataGetPeek, testGetPeek))
	t.Run("get", makeScannerTest(dataGet, testGet))
	t.Run("peek", makeScannerTest(dataGet, testPeek))
}

var (
	dataGetPeek = []testScanner{
		{
			input:  "(define f\n\n(lambda\n (x) x))",
			result: "((ddeeffiinnee  ff\n\n\n\n((llaammbbddaa\n\n  ((xx))  xx))))",
		},
	}

	dataGet = []testScanner{
		{
			input: `(define f
(lambda (x) x))`,
			result: "(define f\n(lambda (x) x))",
		},
	}
)

func testGetPeek(input string) string {
	r := strings.NewReader(input)
	sc := newScanner(r)

	var result []rune
	for sc.notEof() {
		ch, _ := sc.peek()
		ch2, _ := sc.get()
		result = append(result, ch, ch2)
	}
	return string(result)
}

func testPeek(input string) string {
	r := strings.NewReader(input)
	sc := newScanner(r)

	var result []rune
	for sc.notEof() {
		ch, _ := sc.peek()
		sc.get()
		result = append(result, ch)
	}
	return string(result)
}

func testGet(input string) string {
	r := strings.NewReader(input)
	sc := newScanner(r)

	var result []rune
	for sc.notEof() {
		ch, _ := sc.get()
		result = append(result, ch)
	}
	return string(result)
}
