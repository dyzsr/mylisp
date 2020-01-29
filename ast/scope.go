package ast

type Scope struct {
	outer  *Scope
	symtab map[*string]interface{}
}

func NewRootScope() *Scope {
	return NewScope(nil)
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		outer:  parent,
		symtab: make(map[*string]interface{}),
	}
}

func (s *Scope) Lookup(name *string) (interface{}, bool) {
	if value, ok := s.symtab[name]; ok {
		return value, true
	}
	if s.outer != nil {
		return s.outer.Lookup(name)
	}
	return nil, false
}

func (s *Scope) Insert(name *string, value interface{}) {
	s.symtab[name] = value
}

func (s *Scope) Assign(name *string, value interface{}) bool {
	if _, ok := s.symtab[name]; ok {
		s.symtab[name] = value
		return true
	}
	if s.outer != nil {
		return s.outer.Assign(name, value)
	}
	return false
}
