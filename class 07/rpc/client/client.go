package main

import (
	"../rpclib"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var ihelp *bool = flag.Bool("h", false, "Print help information")
	var iport *int = flag.Int("p", 6666, "Port number")
	var ihost *string = flag.String("H", "localhost", "Host address")
	var iverb *int = flag.Int("v", 1, "Verbosity (0-2)")
	flag.Parse()
	if *ihelp {
		flag.Usage()
		os.Exit(0)
	}
	rpclib.SetVerbosity(*iverb)
	var port int = *iport
	var host string = *ihost
	if flag.NArg() > 0 {
		ok := true
		fields := strings.Split(flag.Arg(0), ":")
		ok = ok && len(fields) == 2
		if ok {
			host = fields[0]
			n, err := fmt.Sscanf(fields[1], "%d", &port)
			ok = ok && n == 1 && err == nil
		}
		if !ok {
			flag.Usage()
			os.Exit(0)
		}
	}
	cli := rpclib.NewSClient(host, port)
	for {
		fmt.Printf("CMD:")
		var cmd, arg string
		n, _ := fmt.Scanln(&cmd, &arg)
		if n < 1 {
			fmt.Printf("Exiting\n")
			os.Exit(0)
		}
		switch (cmd) {
		case "i":
			if n < 2 {
				fmt.Printf("Need argument\n")
			} else {
				var iarg int
				n, _ = fmt.Sscanf(arg, "%d", &iarg)
				if n > 0 {
					cli.Insert(arg)
				} else {
					cli.Insert(arg)
				}
			}
		case "?":
			v := cli.Front()
			fmt.Printf("Front value %v\n", v)
		case "r":
			v := cli.Remove()
			fmt.Printf("Removed value %v\n", v)
		case "e":
			e := cli.Empty()
			if e {
				fmt.Printf("empty\n")
			} else {
				fmt.Printf("not empty\n")
			}
		case "f":
			cli.Flush()
			fmt.Printf("Flushed\n")
		case "c":
			v := cli.Contents()
			fmt.Printf("Contents: %v\n", v)
		case "l":
			b := cli.List()
			fmt.Printf("List: %v\n", b.Contents())
		case "q":
			fmt.Printf("Quitting\n")
			os.Exit(0)
		default:
			fmt.Printf("Unrecognized command '%s'\n", cmd)
		}
	}
	fmt.Printf("Done\n")
}
