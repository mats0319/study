package heap

type Heap interface {
	Push(int)
	Pop() (int, error)
}
