package collections

// Comparable is received by types usable in a Tree.  Returns -1 if receiver
// is ordered less than, 1 if ordered greater than, and 0 if equal to the argument.
type Comparable interface {
	CompareTo(o Comparable) int8
}
