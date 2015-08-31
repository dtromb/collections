package tree

import c "github.com/dtromb/collections"

type Tree interface {
	Size() uint
	Has(data c.Comparable) bool
	Lookup(lt c.LookupType, data c.Comparable) (c.Comparable, bool)
	Insert(data c.Comparable) c.Comparable
	Delete(data c.Comparable) (c.Comparable, bool)
	GetCursor(lt c.LookupType, data c.Comparable) (c.Cursor, bool)
	First() c.Cursor
	Last() c.Cursor
}

type TreeImplementation int

const (
	UNKNOWN TreeImplementation = iota
	AVL_THREAD
)

func NewTree(impl ...TreeImplementation) Tree {
	if len(impl) == 0 {
		return &AvlTree{}
	}
	if len(impl) > 1 {
		panic("NewTree() may take at most one argument")
	}
	switch impl[0] {
	case AVL_THREAD:
		return &AvlTree{}
	}
	panic("Unknown tree implementation requested")
}
