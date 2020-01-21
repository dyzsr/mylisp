package runtime

type Kind int

const (
	BOOLEAN = iota
	INTEGER
	SYMBOL
	BUILTIN_PROC
	PROC
)

type Type struct {
	kind Kind
}

func (t *Type) Kind() Kind {
	return t.kind
}
