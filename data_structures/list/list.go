package list

type List interface {
	PrintList()

	// Find return first index matched the key, and '-1' if not matched
	Find(key int) int

	// Insert insert node at tail of list, add in situ, return list itself
	// it will not use return value as usual, except List interface is nil
	// because method receiver's modify is ineffective outside the method in go
	Insert(key int) List

	// Delete delete node with target key, return if delete successfully
	// if list is nil original or after del, return an empty list with false flag
	//
	Delete(key int) (List, bool)

	// IsEmpty judge if a list is nil or only has a default node
	IsEmpty() bool

	// IsNil return if a list is nil
	IsNil() bool
}
