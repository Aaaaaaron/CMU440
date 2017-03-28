// Testing code for buffer

package bufb

import (
	"encoding/json"
	"math/rand"
	"testing"
)

// Convert integer to byte array
func i2b(i int) []byte {
	b, _ := json.Marshal(i)
	return b
}

// Convert byte array back to integer
func b2i(b []byte) int {
	var i int
	json.Unmarshal(b, &i)
	return i
}

// How many repetitions
var ntest int = 10

// How many elements per test
var nele int = 50

func TestBuf(t *testing.T) {
	// Run same test ntest times
	for i := 0; i < ntest; i++ {
		bp := NewBuf()
		runtest(t, bp)
		if !bp.Empty() {
			t.Logf("Expected empty buffer")
			t.Fail()
		}
	}
}

func runtest(t *testing.T, bp *Buf) {
	inserted := 0
	removed := 0
	emptycount := 0
	for removed < nele {
		if bp.Empty() {
			emptycount++
		}
		// Choose action: insert or remove
		insert := !(inserted == nele)
		if inserted > removed && rand.Int31n(2) == 0 {
			insert = false
		}
		if insert {
			bp.Insert(i2b(inserted))
			inserted++
		} else {
			b, err := bp.Remove()
			if err != nil {
				t.Logf("Attempt to remove from empty buffer\n")
				t.Fail()
			}
			v := b2i(b)
			if v != removed {
				t.Logf("Removed %d.  Expected %d\n", v, removed)
				t.Fail()
			}
			removed++
		}
	}
}
