// Uses second buffer to hold deferred requests

package dserver

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

// Which operations require waiting when buffer is empty
var deferOnEmpty = map[int]bool{doremove: true, dofront: true }

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
	// Buffer to hold data
	sb := bufi.NewBuf()
	// Buffer to hold deferred requests
	db := bufi.NewBuf()
	for {
		var r *request
		if !sb.Empty() && !db.Empty() {
			// Revisit deferred operation
			b, _ := db.Remove()
			r = b.(*request)
		} else {
			r = <-bp.requestc
			if sb.Empty() && deferOnEmpty[r.op] {
				// Must defer this operation
				db.Insert(r)
				continue
			}
		}
		switch r.op {
		case doinsert:
			sb.Insert(r.val)
			r.replyc <- nil
		case doremove:
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
