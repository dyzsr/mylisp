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
			input: &Pair{[2]Value{
				Int(123),
				&Pair{[2]Value{
					Int(456),
					&Pair{[2]Value{
						Int(789),
						Nil{},
					}},
				}},
			}},
			result: "(123 456 789)",
		},
		{
			input: &Pair{[2]Value{
				&Pair{[2]Value{Bool(true), &Symbol{&symbols[1]}}},
				&Pair{[2]Value{
					Int(456),
					Int(789),
				}},
			}},
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
