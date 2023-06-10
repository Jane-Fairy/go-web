package socket

import (
	"fmt"
	"log"
	"net"
	"os"
)

func init() {
	SocketDemo()
}

func SocketDemo() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	arg1 := os.Args[1]
	log.Fatalf("system args one : %s", arg1)
	addr := net.ParseIP(arg1)
	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("The address is ", addr.String())
	}
	os.Exit(0)
}
