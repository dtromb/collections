package avltree

// Comparable is received by types usable in a Tree.  Returns -1 if receiver
// is ordered less than, 1 if ordered greater than, and 0 if equal to the argument.
type Comparable interface {
	CompareTo(o Comparable) int8
}

// Tree is an implementation of an AVL Balanced Binary Tree, with threads.
// O(1) iteration and O(ln(n)) time for other operations.
type Tree struct {
	size uint
	root *avlNode
	head *avlNode
	tail *avlNode
}

// Cursor is an iterator type that provides access to the elements of the tree
// as a list.
type Cursor struct {
	tree     *Tree
	nextNode *avlNode
	end      bool
}

type avlNode struct {
	data              Comparable
	balance           int8
	l, r, p, nxt, prv *avlNode
}

// LookupType is the "direction" of the lookup to perform.
type LookupType int

const (
	// GTE searches for the smallest value >= search argument.
	GTE LookupType = iota
	// LTE searches for the largest value <= search argument.
	LTE
)

// Lookup finds a value in the tree according to the given parameters.
func (t *Tree) Lookup(lt LookupType, data Comparable) (Comparable, bool) {
	n, exact := t.lookupNode(lt, data)
	if n == nil {
		return nil, exact
	}
	return n.data, exact
}

func (t *Tree) lookupNode(lt LookupType, data Comparable) (*avlNode, bool) {
	if t.root == nil {
		return nil, false
	}
	cn := t.root
	for {
		r := data.CompareTo(cn.data)
		if r == 0 {
			return cn, true
		}
		if r < 0 {
			if cn.l == nil {
				if lt == GTE {
					return cn, false
				}
				return cn.prv, false
			}
			cn = cn.l
		} else {
			if cn.r == nil {
				if lt == LTE {
					return cn, false
				}
				return cn.nxt, false
			}
			cn = cn.r
		}
	}
}

// Size returns the number of elements in the tree.
func (t *Tree) Size() uint {
	return t.size
}

func (t *Tree) rebalanceAtNode(cn *avlNode) (newParent *avlNode, heightChange bool) {
	switch cn.balance {
	case -2:
		{
			ln := cn.l
			switch ln.balance {
			case -1:
				{
					ln.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = ln
						} else {
							cn.p.r = ln
						}
					} else {
						t.root = ln
					}
					ln.balance = 0
					lrn := ln.r
					cn.p = ln
					ln.r = cn
					cn.l = lrn
					if lrn != nil {
						lrn.p = cn
					}
					cn.balance = 0
					return ln, true
				}
			case 0:
				{
					ln.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = ln
						} else {
							cn.p.r = ln
						}
					} else {
						t.root = ln
					}
					ln.balance = 1
					cn.l = ln.r
					if cn.l != nil {
						cn.l.p = cn
					}
					cn.balance = -1
					ln.r = cn
					cn.p = ln
					return ln, false
				}
			case 1:
				{
					lrn := ln.r
					lrn.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = lrn
						} else {
							cn.p.r = lrn
						}
					} else {
						t.root = lrn
					}
					ln.r = lrn.l
					if ln.r != nil {
						ln.r.p = ln
					}
					if lrn.balance == 1 {
						ln.balance = -1
					} else {
						ln.balance = 0
					}
					cn.l = lrn.r
					if cn.l != nil {
						cn.l.p = cn
					}
					if lrn.balance == -1 {
						cn.balance = 1
					} else {
						cn.balance = 0
					}
					lrn.l = ln
					ln.p = lrn
					lrn.r = cn
					cn.p = lrn
					lrn.balance = 0
					return lrn, true
				}
			}
		}
	case 2:
		{
			rn := cn.r
			switch rn.balance {
			case -1:
				{
					rln := rn.l
					rln.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = rln
						} else {
							cn.p.r = rln
						}
					} else {
						t.root = rln
					}
					cn.r = rln.l
					if cn.r != nil {
						cn.r.p = cn
					}
					if rln.balance == 1 {
						cn.balance = -1
					} else {
						cn.balance = 0
					}
					rn.l = rln.r
					if rn.l != nil {
						rn.l.p = rn
					}
					if rln.balance == -1 {
						rn.balance = 1
					} else {
						rn.balance = 0
					}
					cn.p = rln
					rln.l = cn
					rn.p = rln
					rln.r = rn
					rln.balance = 0
					return rln, true
				}
			case 0:
				{
					rn.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = rn
						} else {
							cn.p.r = rn
						}
					} else {
						t.root = rn
					}
					rn.balance = -1
					cn.r = rn.l
					if cn.r != nil {
						cn.r.p = cn
					}
					cn.p = rn
					rn.l = cn
					cn.balance = 1
					return rn, false
				}
			case 1:
				{
					rn.p = cn.p
					if cn.p != nil {
						if cn.p.l == cn {
							cn.p.l = rn
						} else {
							cn.p.r = rn
						}
					} else {
						t.root = rn
					}
					rn.balance = 0
					cn.r = rn.l
					if cn.r != nil {
						cn.r.p = cn
					}
					cn.balance = 0
					cn.p = rn
					rn.l = cn
					return rn, true
				}
			}
		}
	}
	return nil, false
}

// Insert adds or replaces a value in the tree.  If there is already an equal value,
// it is replaced and the old value returned.  Otherwise, nil is returned.
func (t *Tree) Insert(data Comparable) Comparable {
	if t.root == nil {
		t.root = &avlNode{
			data: data,
		}
		t.size++
		t.head = t.root
		t.tail = t.root
		return nil
	}
	cn := t.root
	var r int8
	for {
		r = data.CompareTo(cn.data)
		if r == 0 {
			break
		}
		if r < 0 {
			if cn.l == nil {
				break
			}
			cn = cn.l
		} else {
			if cn.r == nil {
				break
			}
			cn = cn.r
		}
	}
	if r == 0 {
		oldData := cn.data
		cn.data = data
		return oldData
	}
	nn := &avlNode{
		data: data,
		p:    cn,
	}
	if r < 0 {
		nn.prv = cn.prv
		nn.nxt = cn
		if cn.prv != nil {
			cn.prv.nxt = nn
		} else {
			t.head = nn
		}
		cn.prv = nn
		cn.l = nn
		cn.balance--
	} else {
		nn.nxt = cn.nxt
		nn.prv = cn
		if cn.nxt != nil {
			cn.nxt.prv = nn
		} else {
			t.tail = nn
		}
		cn.nxt = nn
		cn.r = nn
		cn.balance++
	}
	t.size++
	for {
		switch cn.balance {
		case 0:
			{
				// Depth did not increase, we are done.
				return nil
			}
		case -1:
			fallthrough
		case 1:
			{
				// Depth increased, propagate.
				if cn.p != nil {
					if cn.p.l == cn {
						cn.p.balance--
					} else {
						cn.p.balance++
					}
					cn = cn.p
				} else {
					return nil // At root, done.
				}
			}
		case -2:
			fallthrough
		case 2:
			{
				np, changed := t.rebalanceAtNode(cn)
				if changed {
					// Change offset the insertion height change, we're done.
					return nil
				}
				cn = np
			}
		}
	}
}

// Delete removes an element from the tree.  If the argument is found, the
// canonical value from the tree is returned along with boolean true.  If not,
// the pair (nil,false) is returned.
func (t *Tree) Delete(data Comparable) (Comparable, bool) {
	if t.root == nil {
		return nil, false
	}
	var returnVal Comparable
	cn := t.root
	for {
		r := cn.data.CompareTo(data)
		//fmt.Printf("%s %d %d\n", dumpNode(cn), data, r)
		if r == 0 {
			returnVal = cn.data
			break
		} else if r > 0 {
			if cn.l == nil {
				return nil, false
			}
			cn = cn.l
		} else {
			if cn.r == nil {
				return nil, false
			}
			cn = cn.r
		}
	}

	var bn *avlNode

	// cn is the node to delete.
	if cn.l != nil { // cn's prv is cn.l or under cn.l
		dn := cn.prv
		if dn.p == cn { // dn.r == nil
			//fmt.Println("C0")
			dn.p = cn.p
			if dn.p != nil {
				if dn.p.l == cn {
					dn.p.l = dn
				} else {
					dn.p.r = dn
				}
			} else {
				t.root = dn
			}
			dn.r = cn.r
			if dn.r != nil {
				dn.r.p = dn
			}
			bn = dn
			bn.balance = cn.balance + 1
		} else { // dn.r == nil
			//fmt.Println("C1")
			bn = dn.p // bn.r == dn
			dn.p = cn.p
			dn.balance = cn.balance
			if dn.p != nil {
				if dn.p.l == cn {
					dn.p.l = dn
				} else {
					dn.p.r = dn
				}
			} else {
				t.root = dn
			}
			dn.r = cn.r
			if dn.r != nil {
				dn.r.p = dn
			}
			bn.r = dn.l
			if bn.r != nil {
				bn.r.p = bn
			}
			dn.l = cn.l
			dn.l.p = dn
			bn.balance--
		}
	} else if cn.r != nil {
		dn := cn.nxt
		if dn.p == cn { // dn.l == nil
			//fmt.Println("C2")
			dn.p = cn.p
			if dn.p != nil {
				if dn.p.l == cn {
					dn.p.l = dn
				} else {
					dn.p.r = dn
				}
			} else {
				t.root = dn
			}
			dn.l = cn.l
			if dn.l != nil {
				dn.l.p = dn
			}
			bn = dn
			bn.balance = cn.balance - 1
		} else { // dn.l == nil
			//fmt.Println("C3")
			bn = dn.p // bn.l == dn
			dn.p = cn.p
			dn.balance = cn.balance
			if dn.p != nil {
				if dn.p.l == cn {
					dn.p.l = dn
				} else {
					dn.p.r = dn
				}
			} else {
				t.root = dn
			}
			dn.l = cn.l
			if dn.l != nil {
				dn.l.p = dn
			}
			bn.l = dn.r
			if bn.l != nil {
				bn.l.p = bn
			}
			dn.r = cn.r
			dn.r.p = dn
			bn.balance++
		}
	} else {
		//fmt.Println("C4")
		if cn.p == nil {
			t.root = nil
			t.size--
			return returnVal, true
		}
		bn = cn.p
		if bn.l == cn {
			bn.l = nil
			bn.balance++
		} else {
			bn.r = nil
			bn.balance--
		}
	}

	if cn.prv != nil {
		cn.prv.nxt = cn.nxt
	} else {
		t.head = cn.nxt
	}
	if cn.nxt != nil {
		cn.nxt.prv = cn.prv
	} else {
		t.tail = cn.prv
	}
	cn.nxt = nil
	cn.prv = nil

	t.size--

	for {
		switch bn.balance {
		case 0:
			{ // subtree has lost height, propagate upward
				//fmt.Println("B0")
				if bn.p == nil {
					return returnVal, true
				}
				if bn.p.l == bn {
					bn.p.balance++
				} else {
					bn.p.balance--
				}
				bn = bn.p
			}
		case -1:
			fallthrough
		case 1:
			{
				//fmt.Println("B1")
				// subtree has not lost height, we are done
				return returnVal, true
			}
		case 2:
			fallthrough
		case -2:
			{
				//fmt.Println("B2")
				np, changed := t.rebalanceAtNode(bn)
				if !changed {
					return returnVal, true
				}
				bn = np
			}
		}
	}
}

// GetCursor opens a cursor whose first value will be the value that would have
// been returned by an equivalent Lookup() call.  The iterator is bidirectional
// and can range past the start of the search.  It is not fail-fast - changes
// in the tree will change its behavior.
func (t *Tree) GetCursor(lt LookupType, data Comparable) (*Cursor, bool) {
	n, exact := t.lookupNode(lt, data)
	c := &Cursor{
		tree:     t,
		nextNode: n,
	}
	if c.nextNode == nil {
		if lt == GTE {
			c.end = true
		}
	}
	return c, exact
}

// HasNext checks for the availability of a next data value from the cursor.
func (c *Cursor) HasNext() bool {
	return c.nextNode != nil ||
		(c.nextNode == nil && c.end && c.tree.size > 0)
}

// HasPrev checks for the availability of a previous data value from the cursor.
func (c *Cursor) HasPrev() bool {
	return (c.nextNode != nil && c.nextNode.prv != nil) ||
		(c.nextNode == nil && !c.end && c.tree.size > 0)
}

// Next retrieves the next value from the cursor.
func (c *Cursor) Next() Comparable {
	if c.nextNode == nil {
		if c.end {
			return nil
		}
		c.nextNode = c.tree.head
	}
	if c.nextNode == nil {
		return nil
	}
	val := c.nextNode.data
	c.nextNode = c.nextNode.nxt
	if c.nextNode == nil {
		c.end = true
	}
	return val
}

// Prev retrieves the previous value from the cursor.  (Note that this is exactly
// what it sounds like - if you "switch directions", Prev() will return the previous
// value returned by Next() - **not the one before that** in the order.  And vice
// versa...)
func (c *Cursor) Prev() Comparable {
	if c.nextNode == nil {
		if c.end == true {
			c.nextNode = c.tree.tail
			val := c.nextNode.data
			c.nextNode = c.nextNode.prv
			c.end = false
			return val
		}
		return nil
	}
	c.nextNode = c.nextNode.prv
	if c.nextNode == nil {
		return nil
	}
	return c.nextNode.data
}
