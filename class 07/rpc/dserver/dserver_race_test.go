// Testing code for concurrent buffer.
// This test generates long sequence of insertions, removals, and flushes
// in attempt to hit race condition

package dserver

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// How many repetitions
var nrun int = 5

// Used to detect when all removers are done
var rdone []bool

// Concurrent testing

func rinsert(bp *Buf, nele int, dchan chan int) {
	for inserted := 1; inserted <= nele; inserted++ {
		bp.Insert(inserted)
	}
	dchan <- 1
}

func rremove(bp *Buf, nele int, dchan chan int, done *bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panicked while removing: %v\n", err)
			*done = true
			dchan <- 1
		}
	}()
	for removed := 1; removed <= nele; removed++ {
		bp.Remove()
	}
	*done = true
	dchan <- 1
}

func rflush(bp *Buf, nele int, dchan chan int) {
	for flushed := 1; flushed <= nele; flushed++ {
		bp.Flush()
	}
	dchan <- 1
}

func alldone() bool {
	done := true
	for _, v := range rdone {
		done = done && v
	}
	return done
}

// Run test with designated number of elements, flushers and removers
func runtest(nele, flushercount, removercount int) {
	runtime.GOMAXPROCS(2 + flushercount + removercount)
	fmt.Printf("Testing with %d flushers & %d removers\n", flushercount, removercount)
	for r := 0; r < nrun; r++ {
		dchan := make(chan int, 2)
		rchan := make(chan int, 1)
		bp := NewBuf()
		rdone = make([]bool, removercount)
		var i, nfele, nrele int
		if flushercount > 0 {
			nfele = nele / flushercount
		} else {
			nfele = 0
		}
		if removercount > 0 {
			nrele = nele / removercount
		} else {
			nrele = 0
		}

		go rinsert(bp, nele, dchan)

		for i = 0; i < flushercount; i++ {
			go rflush(bp, nfele, dchan)
		}

		for i = 0; i < removercount; i++ {
			go rremove(bp, nrele, rchan, &rdone[i])
		}

		// Synchronize completion of insertion & flushing
		for i := 0; i < 1+flushercount; i++ {
			<-dchan
		}

		// Let things run a while
		time.Sleep(3 * 1e9)
		var xinsert int
		for xinsert = 1; !alldone(); xinsert++ {
			bp.Insert(0)
		}
		for i = 0; i < removercount; i++ {
			<-rchan
		}
		fmt.Printf("Completed %d insertions, %dx%d flushes, %dx%d removals, %d extra insertions\n",
			nele, nfele, flushercount, nrele, removercount, xinsert)
	}
	runtime.GOMAXPROCS(1)
}

func TestRace0(t *testing.T) {
	runtest(1000, 0, 1)
}

func TestRace1(t *testing.T) {
	runtest(1000, 1, 1)
}

func TestRace2(t *testing.T) {
	runtest(1000, 0, 2)
}

func TestRace3(t *testing.T) {
	runtest(1000, 1, 2)
}

func TestRace4(t *testing.T) {
	runtest(1000, 2, 2)
}
