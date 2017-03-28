// Testing code for channel buffer

package cbuf

import (
	"fmt"
	"math/rand"
	"testing"
)

// How many repetitions
var ntest int = 10

// How many elements per test
var nele int = 15

func insert(cb CBuf, echan chan int, t *testing.T) {
	for inserted := 1; inserted <= nele; inserted++ {
		cb.Insert(inserted)
		echan <- inserted
	}
}

func remove(cb CBuf, echan chan int, t *testing.T) {
	for removed := 1; removed <= nele; removed++ {
		v := cb.Remove().(int)
		if v != removed {
			t.Logf("Removed %d, Expected %d\n", v, removed)
			t.Fail()
		}
		echan <- -removed
	}
}

func TestBuf(t *testing.T) {
	// Run same test ntest times
	for i := 0; i < ntest; i++ {
		size := rand.Intn(5)
		cb := NewBuf(size)
		fmt.Printf("Size %d:", size)
		echan := make(chan int, 2*nele)
		go insert(cb, echan, t)
		go remove(cb, echan, t)
		for v := 0; v < 2*nele; v++ {
			x := <-echan
			if x >= 0 {
				fmt.Printf(" I%d", x)
			} else {
				fmt.Printf(" R%d", -x)
			}
		}
		fmt.Printf("\n")
	}
}
