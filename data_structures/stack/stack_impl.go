package mario

type stackImpl struct {
	// make sure elements will never be nil,
	// all methods depend on this suppose, for that, we shadow func 'NilStack()'
	elements []ElemType

	// top can be used as index directly
	top int
}

var _ Stack = (*stackImpl)(nil)

func NewStack() Stack {
	return MakeStack()
}

func MakeStack(elems ...ElemType) Stack {
	return &stackImpl{
		elements: append(make([]ElemType, 0), elems...),
		top:      len(elems) - 1,
	}
}

func (s *stackImpl) Push(elems ...ElemType) {
	s.elements = append(s.elements, elems...)
	s.top = len(s.elements) - 1

	return
}

func (s *stackImpl) Pop() (ElemType, bool) {
	if s.top == -1 {
		return ZeroValue, false
	}

	value := s.elements[s.top]

	s.elements = s.elements[:len(s.elements)-1] // s[:0] is not nil
	s.top--

	return value, true
}

func (s *stackImpl) Top() int {
	return s.top
}

func (s *stackImpl) IsEmpty() bool {
	return s.top == -1
}

func (s *stackImpl) Empty() {
	*s = stackImpl{
		elements: make([]ElemType, 0),
		top:      -1,
	}

	return
}
