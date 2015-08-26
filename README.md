# avltree

### AVL Balanced Binary Tree for Go

--
    import "github.com/dtromb/avltree"


## Usage

#### type Comparable

```go
type Comparable interface {
	CompareTo(o Comparable) int8
}
```

Comparable is received by types usable in a Tree. Returns -1 if receiver is
ordered less than, 1 if ordered greater than, and 0 if equal to the argument.

#### type LookupType

```go
type LookupType int
```

LookupType is the "direction" of the lookup to perform.

```go
const (
	// GTE searches for the smallest value >= search argument.
	GTE LookupType = iota
	// LTE searches for the largest value <= search argument.
	LTE
)
```

#### type Tree

```go
type Tree struct {
}
```

Tree is an implementation of an AVL Balanced Binary Tree, with threads. O(1)
iteration and O(ln(n)) time for other operations.

#### func (*Tree) Delete

```go
func (t *Tree) Delete(data Comparable) (Comparable, bool)
```
Delete removes an element from the tree. If the argument is found, the canonical
value from the tree is returned along with boolean true. If not, the pair
(nil,false) is returned.

#### func (*Tree) Insert

```go
func (t *Tree) Insert(data Comparable) Comparable
```
Insert adds or replaces a value in the tree. If there is already an equal value,
it is replaced and the old value returned. Otherwise, nil is returned.

#### func (*Tree) Lookup

```go
func (t *Tree) Lookup(lt LookupType, data Comparable) (Comparable, bool)
```
Lookup finds a value in the tree according to the given parameters.

#### func (*Tree) Size

```go
func (t *Tree) Size() uint
```
Size returns the number of elements in the tree.
