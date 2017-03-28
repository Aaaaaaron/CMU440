// Use separate channels for read requests vs. other operations

package sserver

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
	// Buffer has two request channels
	opc   chan *request // Nonblocking operations
	readc chan *request // Operations that block for empty buffer
}

////////////////////////////////////////////////////////////////////
// Server implementation

func (bp *Buf) runServer() {
	// Create actual buffer
	sb := bufi.NewBuf()
	for {
		var r *request
		if sb.Empty() {
			r = <-bp.opc
		} else {
			select {
			case r1 := <-bp.opc:
				r = r1
			case r2 := <-bp.readc:
				r = r2
			}
		}
		switch r.op {
		case doinsert:
			sb.Insert(r.val)
			r.replyc <- nil
		case doremove:
			v := sb.Remove()
			r.replyc <- v
		case doflush:
			sb.Flush()
			r.replyc <- nil
		case doempty:
			e := sb.Empty()
			// Can send Boolean along channel
			r.replyc <- e
		case dofront:
			v := sb.Front()
			r.replyc <- v
		}
	}
}

func NewBuf() *Buf {
	bp := &Buf{make(chan *request), make(chan *request)}
	go bp.runServer()
	return bp
}

func (bp *Buf) doop(op int, val interface{}) interface{} {
	r := &request{op, val, make(chan interface{})}
	bp.opc <- r
	v := <-r.replyc
	return v
}

func (bp *Buf) doread(op int, val interface{}) interface{} {
	r := &request{op, val, make(chan interface{})}
	bp.readc <- r
	v := <-r.replyc
	return v
}

// Exported interface

func (bp *Buf) Insert(val interface{}) {
	bp.doop(doinsert, val)
}

func (bp *Buf) Front() interface{} {
	return bp.doread(dofront, nil)
}

func (bp *Buf) Remove() interface{} {
	return bp.doread(doremove, nil)
}

func (bp *Buf) Empty() bool {
	v := bp.doop(doempty, nil)
	e := v.(bool)
	return e
}

func (bp *Buf) Flush() {
	bp.doop(doflush, nil)
}
