package combinatorics

import (
	"testing"
	"fmt"
)

func TestPermutation(t *testing.T) {
	p := FirstPermutation(5)
	i := 1
	fmt.Printf("%d: %s | %s\n", i, p.String(), p.CodeString())
	for p.HasNext() {
		p = p.Next()
		i++
		fmt.Printf("%d: %s | %s\n", i, p.String(), p.CodeString())
	}
}

func TestSelection(t *testing.T) {
	s := FirstSelection(6)
	i := 1
	fmt.Printf("%d: %s\n", i, s.String())
	for s.HasNext() {
		s = s.Next()
		i++
		fmt.Printf("%d: %s\n", i, s.String())
	}
}