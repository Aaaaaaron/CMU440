// Unbounded, synchronous buffer based on heterogenous locking

package abuf

import (
	"../bufi"
)

type Buf struct {
	sb       *bufi.Buf // Sequential buffer
	ackchan  chan int  // Signals completion of operation
	readchan chan int  // Allows blocking when reading
	opchan   chan int  // For nonblocking operations
}

func NewBuf() *Buf {
	bp := new(Buf)
	bp.sb = bufi.NewBuf()
	bp.ackchan = make(chan int)
	bp.readchan = make(chan int)
	bp.opchan = make(chan int)
	go bp.director()
	return bp
}

// Go routine to respond to requests
func (bp *Buf) director() {
	for {
		if bp.sb.Empty() {
			// Enable only nonblocking operations
			bp.opchan <- 1
		} else {
			// Enable reads and other operations
			select {
			case bp.readchan <- 1:
			case bp.opchan <- 1:
			}
		}
		<-bp.ackchan // Wait until operations completed
	}
}

func (bp *Buf) startop() { <-bp.opchan }

func (bp *Buf) startread() { <-bp.readchan }

func (bp *Buf) finish() { bp.ackchan <- 1 }

func (bp *Buf) Insert(val interface{}) {
	bp.startop()
	bp.sb.Insert(val)
	bp.finish()
}

func (bp *Buf) Front() interface{} {
	bp.startread()
	v, _ := bp.sb.Front()
	bp.finish()
	return v
}

func (bp *Buf) Remove() interface{} {
	bp.startread()
	v, _ := bp.sb.Remove()
	bp.finish()
	return v
}

func (bp *Buf) Empty() bool {
	bp.startop()
	rval := bp.sb.Empty()
	bp.finish()
	return rval
}

func (bp *Buf) Flush() {
	bp.startop()
	bp.sb.Flush()
	bp.finish()
}
