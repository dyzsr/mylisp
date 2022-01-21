package runtime

type Kind int

const (
	NIL = iota
	BOOLEAN
	INTEGER
	SYMBOL
	PAIR
	BUILTIN_PROC
	PROC
)

type Type struct {
	kind Kind
}

func (t Type) Kind() Kind {
	return t.kind
}

var (
	TypeNil         = Type{kind: NIL}
	TypeBool        = Type{kind: BOOLEAN}
	TypeInt         = Type{kind: INTEGER}
	TypeSymbol      = Type{kind: SYMBOL}
	TypePair        = Type{kind: PAIR}
	TypeBuiltinProc = Type{kind: BUILTIN_PROC}
	TypeProc        = Type{kind: PROC}
)
