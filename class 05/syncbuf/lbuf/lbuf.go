// Unbounded, synchronous buffer based on locking.
// One of several attempts to implement buffer using only mutex's.
// Uses mutex as binary semaphore to get producer/consumer synchronization
// WARNING: This code can deadlock!

package lbuf

import (
	"../bufi"
	"sync"
)

type Buf struct {
	sb     *bufi.Buf   // Sequential buffer
	mutex  *sync.Mutex // Lock to guarantee mutual exclusion
	readok *sync.Mutex // Locked when buffer empty
}

func NewBuf() *Buf {
	bp := new(Buf)
	bp.sb = bufi.NewBuf()
	bp.mutex = new(sync.Mutex)
	bp.readok = new(sync.Mutex)
	bp.waitonempty() // Lock when empty
	return bp
}

func (bp *Buf) lock() {
	bp.mutex.Lock()
}

func (bp *Buf) unlock() {
	bp.mutex.Unlock()
}

func (bp *Buf) waitonempty() {
	bp.readok.Lock()
}

func (bp *Buf) setnonempty() {
	bp.readok.Unlock()
}

func (bp *Buf) Insert(val interface{}) {
	bp.lock()
	wasempty := bp.sb.Empty()
	bp.sb.Insert(val)
	if wasempty {
		bp.setnonempty()
	}
	bp.unlock()
}

func (bp *Buf) Front() interface{} {
	for {
		bp.waitonempty()
		// Possible for buffer to get flushed right here
		bp.lock()
		if bp.sb.Empty() {
			// Try again
			bp.unlock()
		} else {
			break
		}
	}
	v, _ := bp.sb.Front()
	bp.setnonempty() // Need to restore nonempty status
	bp.unlock()
	return v
}

// Remove should wait until buffer nonempty before removing element
func (bp *Buf) Remove() interface{} {
	for {
		bp.waitonempty()
		// Possible for buffer to get flushed right here
		bp.lock()
		if bp.sb.Empty() {
			// Try again
			bp.unlock()
			// Need to let Flush get past waitonempty
			bp.setnonempty()
		} else {
			break
		}
	}
	v, _ := bp.sb.Remove()
	if !bp.sb.Empty() {
		bp.setnonempty() // Mark buffer as nonempty
	}
	bp.unlock()
	return v
}

func (bp *Buf) Empty() bool {
	bp.lock()
	e := bp.sb.Empty()
	bp.unlock()
	return e
}

func (bp *Buf) Flush() {
	bp.lock()
	wasempty := bp.sb.Empty()
	bp.sb.Flush()
	if !wasempty {
		// Clear nonempty status
		// Unfortunately, can deadlock here, since
		// Insert can't execute until lock released
		bp.waitonempty()
	}
	bp.unlock()

}
