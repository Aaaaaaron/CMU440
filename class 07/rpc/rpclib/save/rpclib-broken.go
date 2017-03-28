// Implementation of a network-based FIFO buffer using RPC

package rpclib

import (
	"../bufi"
	"../dserver"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Error reporting
// Log result if verbosity level high enough
var verbosity int = 1

func SetVerbosity(verb int) {
	verbosity = verb
}

func Vlogf(level int, format string, v ...interface{}) {
	if level <= verbosity {
		log.Printf(format, v...)
	}
}

// Handle errors
func Checkreport(level int, err error) bool {
	if err == nil {
		return false
	}
	Vlogf(level, "Error: %s", err.Error())
	return true
}

func Checkfatal(err error) {
	if err == nil {
		return
	}
	log.Fatalf("Fatal: %s", err.Error())
}

type iVal interface{}

func nullVal() iVal {
	return iVal(nil)
}

func trueval() iVal {
	return iVal(true)
}

func falseval() iVal {
	return iVal(false)
}

func truth(v *iVal) bool {
	return (*v).(bool)
}

type Islice []interface{}

var islice Islice
var sbuf *bufi.Buf

////////////////////////////////////////////////////////////////////
// Server implementation
type SrvBuf struct {
	abuf *dserver.Buf
}

func NewSrvBuf() *SrvBuf {
	gob.Register(islice)
	gob.Register(sbuf)
	return &SrvBuf{dserver.NewBuf()}
}

// Server functions

// arg = value to insert, no reply
func (srv *SrvBuf) Insert(arg *iVal, reply *iVal) error {
	srv.abuf.Insert(*arg)
	*reply = nullVal()
	Vlogf(2, "Inserted %v\n", *arg)
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	return nil
}

// No argument, reply = front value
func (srv *SrvBuf) Front(arg *iVal, reply *iVal) error {
	v := srv.abuf.Front()
	*reply = iVal(v)
	Vlogf(2, "Front value %v\n", v)
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	return nil
}

// No argument, reply = front value
func (srv *SrvBuf) Remove(arg *iVal, reply *iVal) error {
	v := srv.abuf.Remove()
	*reply = iVal(v)
	Vlogf(2, "Removed %v\n", v)
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	return nil
}

// No argument, reply = boolean
func (srv *SrvBuf) Empty(arg *iVal, reply *iVal) error {
	v := falseval()
	if srv.abuf.Empty() {
		v = trueval()
	}
	*reply = iVal(v)
	Vlogf(2, "Empty? %v\n", v)
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	return nil
}

// No argument, no reply
func (srv *SrvBuf) Flush(arg *iVal, reply *iVal) error {
	srv.abuf.Flush()
	*reply = nullVal()
	Vlogf(2, "Flushed\n")
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	return nil
}

// No argument, reply = slice
func (srv *SrvBuf) Contents(arg *iVal, reply *iVal) error {
	c := Islice(srv.abuf.Contents())
	*reply = iVal(c)
	Vlogf(2, "Generated contents: %v\n", c)
	return nil
}

// No argument, reply = *bufi.Buf
func (srv *SrvBuf) List(arg *iVal, reply *iVal) error {
	b := srv.abuf.List()
	*reply = iVal(b)
	Vlogf(2, "Generated list: %v\n", b.Contents())
	return nil
}

func Serve(port int) {
	srv := NewSrvBuf()
	rpc.Register(srv)
	rpc.HandleHTTP()
	addr := fmt.Sprintf(":%d", port)
	l, e := net.Listen("tcp", addr)
	Checkfatal(e)
	Vlogf(1, "Running server on port %d\n", port)
	Vlogf(3, "Buffer: %s\n", srv.abuf.String())
	http.Serve(l, nil)
}

//////////////////////////////////////////////////////
// Client using all synchronous calls

type SClient struct {
	client *rpc.Client
}

func (cli *SClient) Call(serviceMethod string, args interface{},
	reply interface{}) error {
	return cli.client.Call(serviceMethod, args, reply)
}

func NewSClient(host string, port int) *SClient {
	gob.Register(islice)
	gob.Register(sbuf)
	hostport := fmt.Sprintf("%s:%d", host, port)
	client, e := rpc.DialHTTP("tcp", hostport)
	Checkfatal(e)
	Vlogf(1, "Connected to %s\n", hostport)
	return &SClient{client}
}

// Client functions

func (cli *SClient) Insert(val interface{}) {
	v := iVal(val)
	var rv iVal
	e := cli.Call("SrvBuf.Insert", &v, &rv)
	Vlogf(2, "Inserted %v\n", val)
	if Checkreport(1, e) {
		fmt.Printf("Insert failure\n")
	}
}

func (cli *SClient) Front() interface{} {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.Front", &av, &rv)
	if Checkreport(1, e) {
		fmt.Printf("Front failure: %s\n", e.Error())
		return nullVal()
	}
	Vlogf(2, "Front value %v\n", rv)
	return rv
}

func (cli *SClient) Remove() interface{} {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.Remove", &av, &rv)
	if Checkreport(1, e) {
		fmt.Printf("Remove failure: %s\n", e.Error())
		return nullVal()
	}
	Vlogf(2, "Removed %v\n", rv)
	return rv
}

func (cli *SClient) Empty() bool {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.Empty", &av, &rv)
	if Checkreport(1, e) {
		fmt.Printf("Empty failure: %s\n", e)
		return false
	}
	Vlogf(2, "Empty? %s\n", rv.(bool))
	return truth(&rv)
}

func (cli *SClient) Flush() {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.Flush", &av, &rv)
	Vlogf(2, "Flushed\n")
	if Checkreport(1, e) {
		fmt.Printf("Flush failure: %s\n", e)
	}
}

func (cli *SClient) Contents() Islice {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.Contents", &av, &rv)
	Vlogf(2, "Contents %v\n", rv.(Islice))
	if Checkreport(1, e) {
		fmt.Printf("Contents failure: %s\n", e.Error())
	}
	return rv.(Islice)
}

func (cli *SClient) List() *bufi.Buf {
	av := nullVal()
	var rv iVal
	e := cli.Call("SrvBuf.List", &av, &rv)
	b := rv.(*bufi.Buf)
	Vlogf(2, "List %v\n", b.Contents())
	if Checkreport(1, e) {
		fmt.Printf("Contents failure: %s\n", e.Error())
	}
	return b
}
