package main

import (
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"bufio"
	"bytes"
	"strconv"
	"time"
	"io"
)

func main(){
	if len(os.Args) <= 2 {
		os.Stderr.WriteString("script_replay timefile scriptfile")
		os.Exit(1)
	}
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(){
		terminal.Restore(fd, oldState)
	}()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	logF, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)

	for s.Scan() {
		seq := bytes.Split(s.Bytes(), []byte{' '})
		if len(seq) < 2 {
			log.Fatal("err timefile format")
		}
		t, err := strconv.ParseFloat(string(seq[0]), 10)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(t * float64(time.Second)))
		i, err := strconv.ParseInt(string(seq[1]), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		io.CopyN(os.Stdout, logF, i)

	}
}
