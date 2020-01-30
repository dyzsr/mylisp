package runtime

type callstack struct {
	frames []stackframe
}

func newCallstack() *callstack {
	return &callstack{
		frames: make([]stackframe, 0),
	}
}

func (s *callstack) empty() bool {
	return len(s.frames) == 0
}

func (s *callstack) push(proc *Proc, params []Value) {
	s.frames = append(s.frames, stackframe{proc: proc, params: params})
}

func (s *callstack) modify(proc *Proc, params []Value) {
	if s.empty() {
		panic("callstack is empty")
	}
	s.frames[len(s.frames)-1] = stackframe{
		proc:     proc,
		params:   params,
		modified: true,
		last:     false,
	}
}

func (s *callstack) unmodify() {
	if s.empty() {
		panic("callstack is empty")
	}
	s.frames[len(s.frames)-1].modified = false
}

func (s *callstack) modified() bool {
	if s.empty() {
		panic("callstack is empty")
	}
	return s.frames[len(s.frames)-1].modified
}

func (s *callstack) last() bool {
	if s.empty() {
		panic("callstack is empty")
	}
	return s.frames[len(s.frames)-1].last
}

func (s *callstack) setLast(v bool) {
	if s.empty() {
		panic("callstack is empty")
	}
	s.frames[len(s.frames)-1].last = v
}

func (s *callstack) top() stackframe {
	if s.empty() {
		panic("callstack is empty")
	}
	return s.frames[len(s.frames)-1]
}

func (s *callstack) pop() {
	if s.empty() {
		panic("callstack is empty")
	}
	s.frames = s.frames[0 : len(s.frames)-1]
}

type stackframe struct {
	proc     *Proc
	params   []Value
	modified bool
	last     bool
}
