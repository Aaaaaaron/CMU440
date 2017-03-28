These directories & files provide examples of Go code
implementing different buffers.  Each has its own test code

syncbuf:
  Using classical synchronization mechanisms
  bufi:
    Sequential buffer
  lbuf:
    Using two binary semaphores.  This code is prone to deadlock!
  cvbuf:
    Using condition variables

chanbuf:
  Using Go channels

  bufi:
    Sequential buffer
  cbuf:
    Bounded buffer using single channel.  Limited functionality
  abuf:
    Use channels as rendezvous
    Separate channels for blocking vs. nonblocking ops

srvbuf:
  Series of implementations where sequential buffer hidden behind server
  Very Go-like.  Here are the implementations:

  bufi:
    Sequential buffer
  bserver:
    First attempt.  Gets stuck when it encounters blocking operation
  sserver:
    Uses separate channels for reads vs. other requests
  dserver:
    Server manages queue of operations that have been deferred.
