// Unbounded buffer, where underlying values are arbitrary values

package bufi

import (
	"errors"
)

// Linked list element
type BufEle struct {
	val  interface{}
	next *BufEle
}

type Buf struct {
	head *BufEle // Oldest element
	tail *BufEle // Most recently inserted element
}

func NewBuf() *Buf {
	return new(Buf)
}

func (bp *Buf) Insert(val interface{}) {
	ele := &BufEle{val: val}
	if bp.head == nil {
		// Inserting into empty list
		bp.head = ele
	} else {
		bp.tail.next = ele
	}
	bp.tail = ele
}

func (bp *Buf) Front() (interface{}, error) {
	if bp.head == nil {
		return nil, errors.New("Empty Buffer")
	}
	return bp.head.val, nil
}

func (bp *Buf) Remove() (interface{}, error) {
	e := bp.head
	if e == nil {
		err := errors.New("Empty Buffer")
		return nil, err
	}
	bp.head = e.next
	// List becoming empty 
	if e == bp.tail {
		bp.tail = nil
	}
	return e.val, nil
}

func (bp *Buf) Empty() bool {
	return bp.head == nil
}

func (bp *Buf) Flush() {
	bp.head = nil
	bp.tail = nil
}
