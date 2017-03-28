// Testing code for buffer

package bufi

import (
	"testing"
	"json"
	"rand"
	"fmt"
        )

// How many repetitions
var ntest int = 10
// How many elements per test
var nele int = 50


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

func btest(t *testing.T, bp *Buf) {
	inserted := 0
	removed := 0
	emptycount := 0
	fmt.Printf("Byte array data: ")
	for removed < nele {
		if bp.Empty() {	emptycount ++ }
		// Choose action: insert or remove
		insert := !(inserted == nele)
		if inserted > removed && rand.Int31n(2) == 0 {
				insert = false
		}
		if insert {
			bp.Insert(i2b(inserted))
			inserted ++
		} else {
			x := bp.Remove()  // Type = interface{}
			b := x.([]byte)   // Type = []byte
			v := b2i(b)
			if v != removed {
				t.Logf("Removed %d.  Expected %d\n", v, removed)
				t.Fail()
			}
			removed ++
		}
	}
	fmt.Printf("Empty buffer %d/%d times\n", emptycount, nele)
}

func itest(t *testing.T, bp *Buf) {
	inserted := 0
	removed := 0
	emptycount := 0
	fmt.Printf("Integer data: ")
	for (removed < nele) {
		if bp.Empty() {	emptycount ++ }
		// Choose action: insert or remove
		insert := !(inserted == nele)
		if inserted > removed && rand.Int31n(2) == 0 {
				insert = false
		}
		if insert {
			bp.Insert(inserted)
			inserted ++
		} else {
			x := bp.Remove()  // Type = interface{}
			v := x.(int)      // Type = int
			if v != removed {
				t.Logf("Removed %d.  Expected %d\n", v, removed)
				t.Fail()
			}
			removed ++
		}
	}
	fmt.Printf("Empty buffer %d/%d times\n", emptycount, nele)
}

func mtest(t *testing.T, bp *Buf) {
	inserted := 0
	removed := 0
	emptycount := 0
	fmt.Printf("Mixed data: ")
	for (removed < nele) {
		if bp.Empty() {	emptycount ++ }
		// Choose action: insert or remove
		insert := !(inserted == nele)
		if inserted > removed && rand.Int31n(2) == 0 {
				insert = false
		}
		if insert {
			if rand.Int31n(2) == 0 {
				// Insert as integer
				bp.Insert(inserted)
			} else {
				// Insert as byte array
				bp.Insert(i2b(inserted))
			}
			inserted ++
		} else {
			x := bp.Remove()  // Type = interface{}
			var iv int
			switch v := x.(type) {
			case int:
				iv = v
			case []byte:
				iv = b2i(v)
			default:
				t.Logf("Invalid data\n")
				t.Fail()
			}
			if iv != removed {
				t.Logf("Removed %d.  Expected %d\n", iv, removed)
				t.Fail()
			}
			removed ++
		}
	}
	fmt.Printf("Empty buffer %d/%d times\n", emptycount, nele)
}


type TestFun func(*testing.T, *Buf)

func testBuf(t *testing.T, f TestFun) {
	// Run same test ntest times
	for i := 0; i < ntest; i++ {
		bp := NewBuf()
		f(t, bp)
		if !bp.Empty() {
			t.Logf("Expected empty buffer")
			t.Fail()
		}
	}
}

func Testb(t *testing.T) {
	testBuf(t, btest)
}

func Testi(t *testing.T) {
	testBuf(t, itest)
}

func Testm(t *testing.T) {
	testBuf(t, mtest)
}
