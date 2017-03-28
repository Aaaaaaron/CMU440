// demonstration of using channels to implement mutex
package mutex

type Mutex struct {
	mc chan int
}

// Create an unlocked mutex
func NewMutex() *Mutex {
	m := &Mutex{make(chan int, 1)}
	// Buffer initially locked
	m.Unlock()
	return m
}

func (m *Mutex) Lock() {
	// Remove token from buffer
	<-m.mc
}

func (m *Mutex) Unlock() {
	// Insert token into buffer
	m.mc <- 1
}
