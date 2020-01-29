package runtime

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_ValueFormat(t *testing.T) {
	symbols := []string{"abc", "def"}

	testData := []struct {
		input  Value
		result string
	}{
		{
			input:  Bool(false),
			result: "false",
		},
		{
			input:  Int(2147483648),
			result: "2147483648",
		},
		{
			input:  Int(-1),
			result: "-1",
		},
		{
			input:  &Symbol{&symbols[0]},
			result: "abc",
		},
		{
			input: &Pair{
				first: Int(123),
				second: &Pair{
					first: Int(456),
					second: &Pair{
						first:  Int(789),
						second: Nil{},
					},
				},
			},
			result: "(123 456 789)",
		},
		{
			input: &Pair{
				first:  &Pair{first: Bool(true), second: &Symbol{&symbols[1]}},
				second: &Pair{first: Int(456), second: Int(789)},
			},
			result: "((true . def) 456 . 789)",
		},
	}

	for _, test := range testData {
		result := fmt.Sprintf("%s", test.input)
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("\ninput: {%s}\nexpect: '%s'\noutput: '%s'", test.input, test.result, result)
		}
	}
}
