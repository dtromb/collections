package collections

// Cursor is an iterator type that provides access to an ordered list of
// Comparable receivers.
type Cursor interface {
	HasNext() bool
	HasPrev() bool
	Next() Comparable
	Prev() Comparable
}
