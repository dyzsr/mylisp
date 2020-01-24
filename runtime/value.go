package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dyzsr/mylisp/ast"
)

type Value interface {
	Type() *Type
}

type (
	Nil struct{}

	Bool bool

	Int int64

	Symbol struct {
		*string
	}

	Pair struct {
		values [2]Value
	}

	BuiltinProc struct {
		name string
		proc func(...Value) (Value, error)
	}

	Proc struct {
		name  string
		outer *ast.Scope
		*ast.LambdaExpr
	}
)

func (Nil) Type() *Type            { return nil }
func (Bool) Type() *Type           { return nil }
func (Int) Type() *Type            { return nil }
func (v *Symbol) Type() *Type      { return nil }
func (v *Pair) Type() *Type        { return nil }
func (v *BuiltinProc) Type() *Type { return nil }
func (v *Proc) Type() *Type        { return nil }

func (Nil) String() string {
	return "nil"
}

func (v Bool) String() string {
	return strconv.FormatBool(bool(v))
}

func (v Int) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func (v *Symbol) String() string {
	return *v.string
}

func (v *Pair) String() string {
	var substr []string
	p := v
	for {
		substr = append(substr, fmt.Sprintf("%s", p.values[0]))
		if _, ok := p.values[1].(Nil); ok {
			break
		}
		u, ok := p.values[1].(*Pair)
		if !ok {
			substr = append(substr, ".", fmt.Sprintf("%s", p.values[1]))
			break
		}
		p = u
	}
	return "(" + strings.Join(substr, " ") + ")"
}

func (v *BuiltinProc) String() string {
	return "<built-in " + v.name + ">"
}

func (v *Proc) String() string {
	if v.name == "" {
		return "<procedure>:proc"
	}
	return "<procedure " + v.name + ">"
}
