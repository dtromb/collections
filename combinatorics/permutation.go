package combinatorics

import (
	"fmt"
	c "github.com/dtromb/collections"
)

type Permutation struct {
	n int
	lcode []int
	perm []int
}

func GetPermutation(pi ...int) *Permutation {
	p := &Permutation{n:len(pi)}
	p.perm = make([]int,p.n)
	copy(p.perm, pi)
	if !p.Valid() {
		panic("invalid permutation")
	}
	p.mklcode()
	return p
}

func (p *Permutation) String() string {
	var buf []byte
	sep := ","
	if p.perm == nil {
		p.mkperm()
	} 
	buf = append(buf, "{"...)
	for i, k := range p.perm {
		if i == p.n-1 {
			sep = ""
		}
		buf = append(buf, fmt.Sprintf("%d%s",k,sep)...)
	}
	buf = append(buf,"}"...)
	return string(buf)
}

func (p *Permutation) CodeString() string {
	var buf []byte
	sep := ","
	if p.lcode == nil {
		p.mklcode()
	} 
	buf = append(buf, "{"...)
	for i, k := range p.lcode {
		if i == p.n-2 {
			sep = ""
		}
		buf = append(buf, fmt.Sprintf("%d%s",k,sep)...)
	}
	buf = append(buf,"}"...)
	return string(buf)
}

func (p *Permutation) Valid() bool {
	if p.lcode == nil {
		p.mklcode()
	}
	for i, k := range p.lcode {
		if k < 0 || k >= p.n-i {
			return false
		}
	}
	return true
}

func (p *Permutation) Order() int {
	return p.n
}

func (p *Permutation) Code(k int) int {
	if p.lcode == nil {
		p.mklcode()
	}
	if k < 0 || k >= p.n-1 {
		return -1
	}
	return p.lcode[k]
}

func (p *Permutation) Index(k int) int {
	if p.perm == nil {
		p.mkperm()
	}
	if k < 0 || k >= p.n {
		return -1
	}
	return p.perm[k]
}

func (p *Permutation) mkperm() {
	p.perm = make([]int, p.n)
	copy(p.perm[0:p.n-1], p.lcode)
	for i := p.n-2; i >= 0; i-- {
		for j := i+1; j < p.n; j++ {
			if p.perm[i] <= p.perm[j] {
				p.perm[j]++
			}
		}
	}
}

func (p *Permutation) mklcode() {
	p.lcode = make([]int, p.n-1)
	copy(p.lcode, p.perm[0:p.n-1])
	for i := 0; i < p.n-2; i++ {
		for j := i+1; j < p.n-1; j++ {
			if p.lcode[i] <= p.lcode[j] {
				p.lcode[j]--
			}
		}
	}
}

func (p *Permutation) HasNext() bool {
	if p.lcode == nil {
		p.mklcode()
	}
	for i, k := range p.lcode {
		if k < p.n-i-1 {
			return true
		}
	}
	return false
}

func (p *Permutation) HasPrev() bool {
	if p.lcode == nil {
		p.mklcode()
	}
	for _, k := range p.lcode {
		if k > 0 {
			return true
		}
	}
	return false
}

func (p *Permutation) Next() *Permutation {
	np := &Permutation{n:p.n,lcode:make([]int,p.n-1)}
	copy(np.lcode,p.lcode)
	for i := np.n-2; i >= 0; i-- {
		if np.lcode[i] < np.n-i-1 {
			np.lcode[i]++
			i++
			for i < np.n-1 {
				np.lcode[i] = 0
				i++
			}
			return np
		}
	}
	return nil
}

func (p *Permutation) Prev() *Permutation {
	np := &Permutation{n:p.n,lcode:make([]int,p.n-1)}
	copy(np.lcode,p.lcode)
	for i := np.n-2; i >= 0; i-- {
		if np.lcode[i] > 0 {
			np.lcode[i]--
			i++
		}
		for i < np.n-1 {
			np.lcode[i] = np.n-i-1
			i++
		}
		return np
	}
	return nil
}

func FirstPermutation(n int) *Permutation {
	return &Permutation{n:n,lcode:make([]int,n-1)}
}

func LastPermutation(n int) *Permutation {
	p := &Permutation{n:n,lcode:make([]int,n-1)}
	for i := 0; i < n-1; i++ {
		p.lcode[i] = p.n-i-1
	}
	return p
}

func (pa *Permutation) CompareTo(c c.Comparable) int8 {
	pb := c.(*Permutation)
	if pa.n < pb.n {
		return -1
	}
	if pa.n > pb.n {
		return 1
	}
	if pa.lcode == nil {
		pa.mklcode()
	}
	if pb.lcode == nil {
		pb.mklcode()
	}
	for i := 0; i < pa.n-1; i++ {
		a := pa.lcode[i]
		b := pb.lcode[i]
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
	}
	return 0
}

