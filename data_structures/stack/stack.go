package mario

type Stack interface {
	// Push add an elem into stack
	Push(elems ...ElemType)

	// Pop remove the top elem of stack and return it
	Pop() (ElemType, bool)

	// Top return the index of top elem in stack, for invalid stack, return -1
	Top() int

	// IsEmpty check if a stack is empty
	IsEmpty() bool

	// Empty empty a stack
	Empty()
}

type ElemType = int

const ZeroValue ElemType = 0
