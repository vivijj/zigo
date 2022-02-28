package main

import "fmt"

type Hasher[T any] interface {
	HashElement([]T) T
}

type Ah struct {
	a int
	b int
}

func (a Ah) HashElement(si []int) int {
	return len(si)
}

func main() {
	a := make(map[int]int)
	fmt.Println(a[1])
	// a := big.Int{}
	// a.SetString("123456789", 10)
	// b := a
	// fmt.Printf("a p : %p", &a)
	// fmt.Printf("b p: %p", &b)
	// fmt.Println("a is ", &a, "b is", &b)

}
