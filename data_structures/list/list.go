package mario

type List interface {
	// Print return format string of a List
	Print() string

	// Find return first index matched the key, and '-1' if not matched
	Find(key int) int

	// Insert insert node at tail of List, and return List itself
	// return value is useless as usual, except List interface is nil
	Insert(key int) List

	// Delete delete node with target key, return List itself with delete result
	// return value is useful as usual, specially del head node:
	// when you del the first node(head node), List itself should be changed but not expect handle it as a special case
	// so, after thinking, description it here and left this bug in code.
	Delete(key int) (List, bool)

	// IsEmpty judge if a List is nil or only has a default node
	IsEmpty() bool

	// IsNil return if a List is nil
	IsNil() bool
}
