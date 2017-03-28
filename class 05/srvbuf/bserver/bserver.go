// Unbounded buffer, where underlying values are byte arrays
// WARNING: This version doesn't work

package bserver

import (
	"../bufi"
)

////////////////////////////////////////////////////////////////////
// Encapsulate requests to buffer server

const (
	doinsert = iota
	doremove
	doflush
	doempty
	dofront
)

type request struct {
	op     int              // What operation is requested
	val    interface{}      // Optional value for operation
	replyc chan interface{} // Channel to which to send response
}

type Buf struct {
	requestc chan *request // Request channel for buffer
}

////////////////////////////////////////////////////////////////////
// Server implementation

func (bp *Buf) runServer() {
	// Create actual buffer
	sb := bufi.NewBuf()
	for {
		r := <-bp.requestc
		switch r.op {
		case doinsert:
			sb.Insert(r.val)
			r.replyc <- nil
		case doremove:
			// Should block if buffer empty!
			v, _ := sb.Remove()
			r.replyc <- v
		case doflush:
			sb.Flush()
			r.replyc <- nil
		case doempty:
			e := sb.Empty()
			// Can send Boolean along channel
			r.replyc <- e
		case dofront:
			// Should block if buffer empty!
			v, _ := sb.Front()
			r.replyc <- v
		}
	}
}

func NewBuf() *Buf {
	bp := &Buf{make(chan *request)}
	go bp.runServer()
	return bp
}

func (bp *Buf) dorequest(op int, val interface{}) interface{} {
	r := &request{op, val, make(chan interface{})}
	bp.requestc <- r
	v := <-r.replyc
	return v
}

// Exported interface

func (bp *Buf) Insert(val interface{}) {
	bp.dorequest(doinsert, val)
}

func (bp *Buf) Front() interface{} {
	return bp.dorequest(dofront, nil)
}

func (bp *Buf) Remove() interface{} {
	return bp.dorequest(doremove, nil)
}

func (bp *Buf) Empty() bool {
	v := bp.dorequest(doempty, nil)
	e := v.(bool)
	return e
}

func (bp *Buf) Flush() {
	bp.dorequest(doflush, nil)
}
