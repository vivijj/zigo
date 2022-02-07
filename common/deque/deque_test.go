package deque

import (
	"fmt"
	"testing"
)

func TestDeque(t *testing.T) {
	de := New[int]()

	fmt.Printf("%+v ", de)

	for i := 0; i < 100; i++ {
		de.PushBack(i)
	}
	fmt.Printf("%+v ", de)
	for j := 0; j < 100; j++ {
		t := de.PopFront()
		fmt.Println("cur int is: ", t)
	}
}
