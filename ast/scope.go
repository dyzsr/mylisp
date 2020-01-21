package ast

type Scope struct {
	parent *Scope
	symtab map[string]interface{}
}

func NewRootScope() *Scope {
	return NewScope(nil)
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent: parent,
		symtab: make(map[string]interface{}),
	}
}

func (s *Scope) Lookup(name string) (interface{}, bool) {
	if value, ok := s.symtab[name]; ok {
		return value, true
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return nil, false
}

func (s *Scope) Insert(name string, value interface{}) {
	s.symtab[name] = value
}
