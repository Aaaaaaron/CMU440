These directories & files provide examples of Go code.
"buf" code implements different buffers.  Each has its own test code

For use within one thread / goroutine

bufb:
  Linked list, where data elements are of type []byte
bufi:
  Linked list, where data elements are generic, using interface types

proxy:
  A UDP proxy.
  Demonstrates use of net library and "map" to implement dictionary