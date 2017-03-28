// Testing code for channel buffer

package cvbuf

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

// How many repetitions
var ntest int = 10

// How many elements per test
var nele int = 15

func insert(bp *Buf, echan chan int, t *testing.T) {
	for inserted := 1; inserted <= nele; inserted++ {
		bp.Insert(inserted)
		echan <- inserted
	}
}

func remove(bp *Buf, echan chan int, t *testing.T) {
	for removed := 1; removed <= nele; removed++ {
		v := bp.Remove().(int)
		if v != removed {
			fmt.Printf("Removed %d, Expected %d\n", v, removed)
			t.Fail()
		}
		echan <- -removed
	}
}

func TestBuf1(t *testing.T) {
	// Run same test ntest times
	fmt.Printf("Running with one thread\n")
	for i := 0; i < ntest; i++ {
		bp := NewBuf()
		echan := make(chan int, 2*nele)
		go insert(bp, echan, t)
		go remove(bp, echan, t)
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

func TestBuf2(t *testing.T) {
	runtime.GOMAXPROCS(2)
	fmt.Printf("Running with two threads\n")
	// Run same test ntest times
	for i := 0; i < ntest; i++ {
		bp := NewBuf()
		echan := make(chan int, 2*nele)
		go insert(bp, echan, t)
		go remove(bp, echan, t)
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
	runtime.GOMAXPROCS(1)
}

func binsert(bp *Buf, echan chan int, bchan chan int, t *testing.T) {
	for inserted := 1; inserted <= nele; inserted++ {
		bchan <- 1
		bp.Insert(inserted)
		echan <- inserted
	}
}

func bremove(bp *Buf, echan chan int, bchan chan int, t *testing.T) {
	for removed := 1; removed <= nele; removed++ {
		<-bchan
		v := bp.Remove().(int)
		if v != removed {
			fmt.Printf("Removed %d, Expected %d\n", v, removed)
			t.Fail()
		}
		echan <- -removed
	}
}

func TestBBuf(t *testing.T) {
	fmt.Printf("Limiting number of buffer elements.  Two threads\n")
	runtime.GOMAXPROCS(2)
	// Run same test ntest times
	for i := 0; i < ntest; i++ {
		size := rand.Intn(5)
		bchan := make(chan int, size)
		bp := NewBuf()
		fmt.Printf("Size %d:", size)
		echan := make(chan int, 2*nele)
		go binsert(bp, echan, bchan, t)
		go bremove(bp, echan, bchan, t)
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
	runtime.GOMAXPROCS(1)
}
