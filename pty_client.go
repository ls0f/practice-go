package main

import (
	"flag"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"log"
	"net"
)

func main(){

	addr := flag.String("a", "", "addr")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.
	go io.Copy(conn, os.Stdin)
	io.Copy(os.Stdout, conn)
}
