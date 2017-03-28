// Unbounded, synchronous buffer based on mutexes and condition variables

package cvbuf

import (
	"../bufi"
	"sync"
)

type Buf struct {
	sb    *bufi.Buf   // Sequential buffer
	mutex *sync.Mutex // Lock to guarantee mutual exclusion
	cvar  *sync.Cond  // Allows blocking when reading
}

func NewBuf() *Buf {
	bp := new(Buf)
	bp.sb = bufi.NewBuf()
	bp.mutex = new(sync.Mutex)
	bp.cvar = sync.NewCond(bp.mutex)
	return bp
}

func (bp *Buf) lock() { bp.mutex.Lock() }

func (bp *Buf) unlock() { bp.mutex.Unlock() }

func (bp *Buf) Insert(val interface{}) {
	bp.lock()
	bp.sb.Insert(val)
	bp.cvar.Signal()
	bp.unlock()
}

// Front should wait until buffer nonempty before returning
func (bp *Buf) Front() interface{} {
	bp.lock()
	for bp.sb.Empty() {
		bp.cvar.Wait()
	}
	v, _ := bp.sb.Front()
	bp.unlock()
	return v
}

// Remove should wait until buffer nonempty before removing element
func (bp *Buf) Remove() interface{} {
	bp.lock()
	for bp.sb.Empty() {
		bp.cvar.Wait()
	}
	v, _ := bp.sb.Remove()
	bp.unlock()
	return v
}

func (bp *Buf) Empty() bool {
	bp.lock()
	rval := bp.sb.Empty()
	bp.unlock()
	return rval
}

func (bp *Buf) Flush() {
	bp.lock()
	bp.sb.Flush()
	bp.unlock()
}
