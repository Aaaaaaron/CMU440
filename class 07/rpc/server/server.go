package main

import (
	"../rpclib"
	"flag"
	"fmt"
	"os"
)

func main() {
	var ihelp *bool = flag.Bool("h", false, "Print help information")
	var iport *int = flag.Int("p", 6666, "Port number")
	var iverb *int = flag.Int("v", 1, "Verbosity (0-2)")
	flag.Parse()
	if *ihelp {
		flag.Usage()
		os.Exit(0)
	}
	rpclib.SetVerbosity(*iverb)
	var port int = *iport
	if flag.NArg() > 0 {
		nread, _ := fmt.Sscanf(flag.Arg(0), "%d", &port)
		if nread != 1 {
			flag.Usage()
			os.Exit(0)
		}
	}
	rpclib.Serve(port)
}
