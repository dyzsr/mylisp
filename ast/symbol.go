package ast

var (
	SymbolMap = NewSymbolMap()
)

func NewSymbolMap() func(string) *string {
	symbolMap := make(map[string]*string)
	return func(name string) *string {
		if value, ok := symbolMap[name]; ok {
			return value
		}
		value := &name
		symbolMap[name] = value
		return value
	}
}
