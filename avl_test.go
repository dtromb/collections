package avltree

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
)

func treeCounts(n *avlNode) (leftDepth int, rightDepth int, size int) {
	var lld, lrd, ls, rld, rrd, rs, ld, rd, s int
	if n == nil {
		return 0, 0, 0
	}
	if n.l != nil {
		lld, lrd, ls = treeCounts(n.l)
	}
	if n.r != nil {
		rld, rrd, rs = treeCounts(n.r)
	}
	if lld > lrd {
		ld = lld + 1
	} else {
		ld = lrd + 1
	}
	if rld > rrd {
		rd = rld + 1
	} else {
		rd = rrd + 1
	}
	s = ls + rs + 1
	return ld, rd, s
}

func checkStructuralIntegrity(t *Tree, n *avlNode) {
	if n.l != nil {
		if n.l.p != n {
			panic("left child's parent is not self")
		}
	}
	if n.r != nil {
		if n.r.p != n {
			fmt.Printf("%s%s\n", dumpNode(n), dumpNode(n.r))
			panic("right child's parent is not self")
		}
	}
	if n.p != nil {
		if !(n.p.l == n && n.p.r != n) &&
			!(n.p.r == n && n.p.l != n) {
			if n.p.r == n.p.l {
				panic("node is double-linked from parent")
			}
			panic("node is not linked by parent")
		}
	} else {
		if n != t.root {
			panic("node without parent is not tree root")
		}
	}
	if n.prv != nil {
		if n.prv.nxt != n {
			panic("node is not linked as predecessor's successor'")
		}
	} else {
		if n.l != nil {
			panic("node without predecessor has left child")
		}
		cn := n.p
		for cn != nil {
			if cn.p != nil && cn.p.r == cn {
				panic("node without predecessor has inorder predecessor")
			}
			cn = cn.p
		}
	}
	if n.nxt != nil {
		if n.nxt.prv != n {
			panic("node is not linked as sucessor's predecessor")
		}
	} else {
		if n.r != nil {
			panic("node without successor has right child")
		}
		cn := n.p
		for cn != nil {
			if cn.p != nil && cn.p.l == cn {
				panic("node without successor has inorder sucessor")
			}
			cn = cn.p
		}
	}
}

func inorderSucc(n *avlNode) *avlNode {
	if n.r != nil {
		cn := n.r
		for cn.l != nil {
			cn = cn.l
		}
		return cn
	}
	cn := n
	for cn.p != nil {
		if cn.p.l == cn {
			return cn.p
		}
		cn = cn.p
	}
	return nil
}

func inorderPred(n *avlNode) *avlNode {
	if n.l != nil {
		cn := n.l
		for cn.r != nil {
			cn = cn.r
		}
		return cn
	}
	cn := n
	for cn.p != nil {
		if cn.p.r == cn {
			return cn.p
		}
		cn = cn.p
	}
	return nil
}

func checkLinks(n *avlNode) {
	if inorderPred(n) != n.prv {
		fmt.Printf("%d %s\n", n.data.(ComparableInt), strnode(inorderPred(n)))
		panic("incorrect predecessor link")
	}
	if inorderSucc(n) != n.nxt {
		panic("incorrect successor link")
	}
}

func checkOrder(n *avlNode) {
	if n.l != nil {
		if n.data.CompareTo(n.l.data) <= 0 {
			panic("left child order violation")
		}
	}
	if n.prv != nil {
		if n.data.CompareTo(n.prv.data) <= 0 {
			panic("predecessor order violation")
		}
	}
	if n.r != nil {
		if n.data.CompareTo(n.r.data) >= 0 {
			panic("right child order violation")
		}
	}
	if n.nxt != nil {
		if n.data.CompareTo(n.nxt.data) >= 0 {
			panic("successor order violation")
		}
	}
}

type ComparableInt int

func (ci ComparableInt) CompareTo(o Comparable) int8 {
	co := o.(ComparableInt)
	if ci < co {
		return -1
	}
	if ci > co {
		return 1
	}
	return 0
}

func checkBalance(n *avlNode) {
	ld, rd, _ := treeCounts(n)
	balance := rd - ld
	if n.balance != int8(balance) {
		panic("incorrect balance record")
	}
	if n.balance < -1 || n.balance > 1 {
		panic("balance out of range")
	}
}

func checkNode(t *Tree, n *avlNode) {
	checkStructuralIntegrity(t, n)
	checkOrder(n)
	checkLinks(n)
	checkBalance(n)
}

func checkNodesRecursive(t *Tree, n *avlNode) {
	if n == nil {
		return
	}
	checkNode(t, n)
	checkNodesRecursive(t, n.l)
	checkNodesRecursive(t, n.r)
}

func checkTree(t *Tree) {
	_, _, size := treeCounts(t.root)
	if size != int(t.Size()) {
		panic("incorrect size in tree")
	}
	checkNodesRecursive(t, t.root)
}

func strnode(n *avlNode) string {
	if n == nil {
		return "nil"
	}
	return strconv.Itoa(int(n.data.(ComparableInt)))
}

func dumpNode(n *avlNode) string {
	if n == nil {
		return fmt.Sprintf("nil")
	}
	return fmt.Sprintf("%s: par=%s left=%s right=%s prv=%s nxt=%s bal=%d\n",
		strnode(n), strnode(n.p), strnode(n.l), strnode(n.r),
		strnode(n.prv), strnode(n.nxt), n.balance)
}

func dumpNodeRecursive(n *avlNode) {
	fmt.Println(dumpNode(n))
	if n.l != nil {
		dumpNodeRecursive(n.l)
	}
	if n.r != nil {
		dumpNodeRecursive(n.r)
	}
}

func dumpTree(t *Tree) {
	dumpNodeRecursive(t.root)
}

func TestAvl(t *testing.T) {
	tree := &Tree{}
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			fmt.Println(string(buf))
			dumpTree(tree)
		}
	}()

	// Inserts
	fmt.Println("Instrumented test /w incremental constraint verification...")
	var N = 10000
	for i := 0; i < N; i++ {
		x := ComparableInt(i)
		old := tree.Insert(x)
		//fmt.Printf("Inserted %d\n", x)
		if old != nil {
			panic("unexpected old value")
		}
		checkTree(tree)
	}
	fmt.Printf("Size = %d\n", tree.Size())
	//dumpTree(tree)
	// Deletions
	var dcount int
	var deleted = make(map[ComparableInt]ComparableInt)
	for i := 0; i < N*3/5; i++ {
		//dumpTree(tree)
		x := ComparableInt(int(rand.Int31()) % N)
		//fmt.Printf("DELETE %d\n",x)
		old, found := tree.Delete(x)
		if _, has := deleted[x]; has {
			if found {
				panic("deleted node found in tree")
			}
		} else {
			dcount++
			deleted[x] = x
			if old != x {
				//fmt.Printf("%d %d\n", old, x)
				panic("wrong node returned after deletion")
			}
		}
		//dumpTree(tree)
		checkTree(tree)
	}
	fmt.Printf("Deleted %d, Size = %d\n", dcount, tree.Size())
	if int(tree.Size()) != N-dcount {
		panic("incorrect size after deletions")
	}
}

func TestAvl2(t *testing.T) {
	tree := &Tree{}
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
			buf := make([]byte, 1024)
			runtime.Stack(buf, false)
			fmt.Println(string(buf))
			dumpTree(tree)
		}
	}()

	// Inserts
	fmt.Println("Operational mode tests...")
	var N = 5000000
	for i := 0; i < N; i++ {
		x := ComparableInt(i)
		old := tree.Insert(x)
		if old != nil {
			panic("unexpected old value")
		}
	}
	fmt.Printf("Size = %d\n", tree.Size())
	// Deletions
	var dcount int
	var deleted = make(map[ComparableInt]ComparableInt)
	for i := 0; i < N*3/5; i++ {
		x := ComparableInt(int(rand.Int31()) % N)
		old, found := tree.Delete(x)
		if _, has := deleted[x]; has {
			if found {
				panic("deleted node found in tree")
			}
		} else {
			dcount++
			deleted[x] = x
			if old != x {
				panic("wrong node returned after deletion")
			}
		}
	}
	fmt.Printf("Deleted %d, Size = %d\n", dcount, tree.Size())
	if int(tree.Size()) != N-dcount {
		panic("incorrect size after deletions")
	}
}
