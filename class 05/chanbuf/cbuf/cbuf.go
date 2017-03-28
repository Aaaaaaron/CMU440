// Asynchronous buffer based on channels

package cbuf

// Basic data structure
type CBuf chan interface{}

func NewBuf(capacity int) CBuf {
	return make(CBuf, capacity)
}

func (cb CBuf) Insert(val interface{}) {
	cb <- val
}

func (cb CBuf) Remove() interface{} {
	return <-cb
}

/* Can't implement these: 

func (cb CBuf) Front() interface{} {
}

func (cb CBuf) Empty() boolean {
}

func (bp *Buf) Flush() {
}

*/
