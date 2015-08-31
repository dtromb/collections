package set

import c "github.com/dtromb/collections"
import tree "github.com/dtromb/collections/tree"

// Set is a comparable immutable value comprised of other underlying
// comparable values.
type Set interface {
	c.Comparable
	Size() int
	Contains(c c.Comparable) bool
	Union(s Set) Set
	Intersection(s Set) Set
	Difference(s Set) Set
	OpenCursor() c.Cursor
	Ordered() bool
}

type MutableSet interface {
	Set
	Add(cs ...c.Comparable)
	Remove(cs ...c.Comparable)
	Retain(cs ...c.Comparable)
	Clear()
}

type emptySet struct{}

func (es *emptySet) CompareTo(c c.Comparable) int8 {
	s := c.(Set)
	r := s.Size()
	if r == 0 {
		return 0
	}
	return -1
}

func (es *emptySet) Ordered() bool { return true }

func (es *emptySet) Size() int { return 0 }

func (es *emptySet) Contains(c c.Comparable) bool { return false }

func (es *emptySet) Union(s Set) Set { return s }

func (es *emptySet) Intersection(s Set) Set { return es }

func (es *emptySet) Difference(s Set) Set { return es }

func (es *emptySet) OpenCursor() c.Cursor {
	return es
}

func (es *emptySet) HasNext() bool {
	return false
}

func (es *emptySet) HasPrev() bool {
	return false
}

func (es *emptySet) Next() c.Comparable {
	return nil
}
func (es *emptySet) Prev() c.Comparable {
	return nil
}

func Empty() Set {
	return &emptySet{}
}

type singletonSet struct {
	x c.Comparable
}

func (ss *singletonSet) Ordered() bool { return true }

func (ss *singletonSet) Size() int { return 1 }

func (ss *singletonSet) Contains(c c.Comparable) bool { return ss.x.CompareTo(c) == 0 }

func (ss *singletonSet) Union(s Set) Set {
	if sso, ok := s.(*singletonSet); ok {
		return Pair(ss.x, sso.x)
	}
	return s.Union(ss)
}

func (ss *singletonSet) Intersection(s Set) Set {
	if s.Contains(ss.x) {
		return ss
	} else {
		return Empty()
	}
}

func (ss *singletonSet) Difference(s Set) Set {
	if s.Contains(ss.x) {
		return Empty()
	} else {
		return ss
	}
}

func (ss *singletonSet) CompareTo(c c.Comparable) int8 {
	if s, isSet := c.(Set); isSet {
		r := s.Size()
		if r == 0 {
			return 1
		}
		if r > 1 {
			return -1
		}
		return ss.x.CompareTo(s.OpenCursor().Next())
	}
	return 1
}

type ssCursor struct {
	set  *singletonSet
	read bool
}

func (ss *singletonSet) OpenCursor() c.Cursor {
	return &ssCursor{set: ss}
}

func (c *ssCursor) HasNext() bool { return !c.read }
func (c *ssCursor) HasPrev() bool { return c.read }
func (c *ssCursor) Next() c.Comparable {
	if c.read {
		return nil
	}
	return c.set.x
}
func (c *ssCursor) Prev() c.Comparable {
	if !c.read {
		return nil
	}
	return c.set.x
}

func Singleton(c c.Comparable) Set {
	return &singletonSet{x: c}
}

type pairSet struct {
	x, y c.Comparable
}

func (ps *pairSet) Ordered() bool { return true }

func (ps *pairSet) Size() int { return 2 }

func (ps *pairSet) Contains(c c.Comparable) bool {
	return ps.x.CompareTo(c) == 0 ||
		ps.y.CompareTo(c) == 0
}

func (ps *pairSet) Union(s Set) Set {
	t := &tree.Tree{}
	t.Insert(ps.x)
	t.Insert(ps.y)
	c := s.OpenCursor()
	for c.HasNext() {
		t.Insert(c.Next())
	}
	if t.Size() == 2 {
		return ps
	}
	return &treeSet{tree: t}
}

func (ps *pairSet) Intersection(s Set) Set {
	if s.Contains(ps.x) {
		if s.Contains(ps.y) {
			return ps
		}
		return Singleton(ps.x)
	}
	if s.Contains(ps.y) {
		return Singleton(ps.y)
	}
	return Empty()
}

func (ps *pairSet) Difference(s Set) Set {
	if s.Contains(ps.x) {
		if s.Contains(ps.y) {
			return Empty()
		}
		return Singleton(ps.y)
	}
	if s.Contains(ps.y) {
		return Singleton(ps.x)
	}
	return ps
}

func (ps *pairSet) CompareTo(c c.Comparable) int8 {
	if s, isSet := c.(Set); isSet {
		r := s.Size()
		if r < 2 {
			return 1
		}
		if r > 2 {
			return -1
		}
		cc := s.OpenCursor()
		rp := ps.x.CompareTo(cc.Next())
		if rp != 0 {
			return rp
		}
		return ps.y.CompareTo(cc.Next())
	}
	return 1
}

type psCursor struct {
	ps  *pairSet
	pos int
}

func (ps *pairSet) OpenCursor() c.Cursor {
	return &psCursor{ps: ps}
}

func (c *psCursor) HasNext() bool {
	return c.pos < 2
}

func (c *psCursor) HasPrev() bool {
	return c.pos > 0
}

func (c *psCursor) Next() c.Comparable {
	c.pos++
	switch c.pos {
	case 1:
		return c.ps.x
	case 2:
		return c.ps.y
	}
	c.pos = 2
	return nil
}

func (pc *psCursor) Prev() c.Comparable {
	var k c.Comparable
	switch pc.pos {
	case 1:
		k = pc.ps.x
	case 2:
		k = pc.ps.y
	}
	pc.pos--
	return k
}

func Pair(c1, c2 c.Comparable) Set {
	r := c1.CompareTo(c2)
	if r == 0 {
		return Singleton(c1)
	}
	if r < 0 {
		return &pairSet{x: c1, y: c2}
	} else {
		return &pairSet{x: c2, y: c1}
	}
}

type treeSet struct {
	tree *tree.Tree
}

func TreeSet(tree *tree.Tree) MutableSet {
	return &treeSet{tree: tree}
}

func (ts *treeSet) Ordered() bool { return true }

func (ts *treeSet) CompareTo(c c.Comparable) int8 {
	if os, isSet := c.(Set); isSet {
		oz := uint(os.Size())
		tz := ts.tree.Size()
		if tz < oz {
			return -1
		}
		if tz > oz {
			return 1
		}
		if os.Ordered() {
			cc := ts.OpenCursor()
			oc := os.OpenCursor()
			for cc.HasNext() {
				r := cc.Next().CompareTo(oc.Next())
				if r != 0 {
					return r
				}
			}
			return 0
		}
		panic("cannot compare non-ordered sets")
	}
	return 1
}

func (ts *treeSet) Size() int { return int(ts.tree.Size()) }

func (ts *treeSet) Contains(cm c.Comparable) bool {
	_, has := ts.tree.Lookup(c.LTE, cm)
	return has
}

func (ts *treeSet) Union(s Set) Set {
	nt := &tree.Tree{}
	tsc := ts.OpenCursor()
	for tsc.HasNext() {
		nt.Insert(tsc.Next())
	}
	sc := s.OpenCursor()
	for sc.HasNext() {
		nt.Insert(sc.Next())
	}
	return TreeSet(nt)
}

func (ts *treeSet) Intersection(s Set) Set {
	if ts.Size() > s.Size() {
		return s.Intersection(ts)
	}
	nt := &tree.Tree{}
	c := ts.OpenCursor()
	for c.HasNext() {
		k := c.Next()
		if s.Contains(k) {
			nt.Insert(k)
		}
	}
	return TreeSet(nt)
}

func (ts *treeSet) Difference(s Set) Set {
	nt := &tree.Tree{}
	c := ts.OpenCursor()
	for c.HasNext() {
		k := c.Next()
		if !s.Contains(k) {
			nt.Insert(k)
		}
	}
	return TreeSet(nt)
}

func (ts *treeSet) OpenCursor() c.Cursor {
	return ts.tree.First()
}

func (ts *treeSet) Add(cs ...c.Comparable) {
	for _, k := range cs {
		ts.tree.Insert(k)
	}
}

func (ts *treeSet) Remove(cs ...c.Comparable) {
	for _, k := range cs {
		ts.tree.Delete(k)
	}
}

func (ts *treeSet) Retain(cs ...c.Comparable) {
	panic("unimplemented")
}

func (ts *treeSet) Clear() {
	ts.tree = &tree.Tree{}
}
