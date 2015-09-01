package combinatorics

import (
	"fmt"
	"math/big"
)

type Selection struct {
	n uint
	x *big.Int
}

func FirstSelection(n uint) *Selection {
	return &Selection{n:n,x:big.NewInt(0)}
}

func LastSelection(n uint) *Selection {
	one := big.NewInt(1)
	z := big.NewInt(0)
	z.Sub(z.Lsh(one,n),one)
	return &Selection{n:n,x:z}
}

func (s *Selection) HasNext() bool {
	z := big.NewInt(1)
	return z.Add(s.x,z).BitLen() <= int(s.n)
}

func (s *Selection) HasPrev() bool {
	return s.x.BitLen() > 0
}

func (s *Selection) Next() *Selection {
	x := *s.x
	ns := &Selection{n:s.n,x:&x}
	ns.x.Add(ns.x,big.NewInt(1))
	return ns
}

func (s *Selection) Prev() *Selection {
	x := *s.x
	ns := &Selection{n:s.n,x:&x}
	ns.x.Sub(ns.x,big.NewInt(1))
	return ns
}


func (s *Selection) Test(idx int) bool {
	return s.x.Bit(idx) > 0
} 

func (s *Selection) String() string {
	var buf []byte 
	sep := ","
	buf = append(buf, "{"...)
	for i := 0; i < int(s.n); i++ {
		if i == int(s.n)-1 {
			sep = ""
		}
		var v int
		if s.Test(i) {
			v = 1
		} else {
			v = 0
		}
		buf = append(buf,fmt.Sprintf("%d%s",v,sep)...)
	}
	buf = append(buf,"}"...)
	return string(buf)
}