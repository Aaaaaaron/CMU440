// Sequential testing code for buffer

package dserver

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

// Set up shadow buffer using slice
type sbuf struct {
	b        []int
	capacity int
}

func newsbuf(capacity int) *sbuf {
	return &sbuf{make([]int, 0, capacity), capacity}
}

func (sp *sbuf) insert(val int) {
	n := len(sp.b)
	sp.b = sp.b[0: n+1]
	sp.b[n] = val
}

func (sp *sbuf) front() int {
	if len(sp.b) == 0 {
		return 0
	}
	return sp.b[0]
}

func (sp *sbuf) remove() int {
	if len(sp.b) == 0 {
		return 0
	}
	val := sp.b[0]
	sp.b = sp.b[1:]
	return val
}

func (sp *sbuf) empty() bool {
	return len(sp.b) == 0
}

func (sp *sbuf) flush() {
	sp.b = make([]int, 0, sp.capacity)
}

func (sp *sbuf) length() int {
	return len(sp.b)
}

// How many repetitions
var nstest int = 20

// How many elements per test
var maxcount int = 80

/* These already declared in server code 
const (
	doinsert = iota
	doremove
	doflush
	doempty
	dofront
	)
*/

var (
	emptyops = [...]int{doinsert, doflush, doempty}
	fullops  = [...]int{doremove, doflush, doempty, dofront}
	allops   = [...]int{doinsert, doremove, doflush, dofront, doempty}
	opchars  = [...]string{"+", "-", "x", "e", "?"}
)

// Apply same operation to two different buffer implementations
// Return single character string indicating operation done
func dorandomop(bp *Buf, sp *sbuf, forceremove bool, t *testing.T) string {
	choices := allops[0:]
	if sp.empty() {
		choices = emptyops[0:]
	} else if sp.length() == maxcount {
		choices = fullops[0:]
	}
	choice := choices[rand.Intn(len(choices))]
	if forceremove {
		choice = doremove
	}
	switch choice {
	case doinsert:
		ival := rand.Int()
		bp.Insert(ival)
		sp.insert(ival)
	case doremove:
		vb := bp.Remove().(int)
		vs := sp.remove()
		if vb != vs {
			t.Logf("Removed %d from bp, %d from sp\n", vb, vs)
			t.Fail()
		}
	case doflush:
		bp.Flush()
		sp.flush()
	case doempty:
		eb := bp.Empty()
		es := sp.empty()
		if eb != es {
			t.Logf("Emptiness mismatch: bp %v, sp %v\n", eb, es)
			t.Fail()
		}
	case dofront:
		vb := bp.Front().(int)
		vs := sp.front()
		if vb != vs {
			t.Logf("Front %d in bp, %d in sp\n", vb, vs)
			t.Fail()
		}
	}
	return opchars[choice]
}

func TestSBuf(t *testing.T) {
	fmt.Printf("Running sequential tests\n")
	// Run same test nstest times
	bp := NewBuf()
	sp := newsbuf(maxcount)
	for i := 0; i < nstest; i++ {
		ops := make([]string, maxcount)
		for j := 0; j < maxcount; j++ {
			forceremove := (maxcount - j) <= sp.length()
			ops[j] = dorandomop(bp, sp, forceremove, t)
			//			fmt.Printf("Len %d. Cap %d, Op %s\n", sp.length(), cap(sp.b), ops[j])
		}
		fmt.Println(strings.Join(ops, ""))
	}
}
