package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"log"
	"time"
	"os/user"
)

const (
	delay = 300 * time.Millisecond
)

type input struct {
	line string
	i int
}

var hostname string
var cuser *user.User

func init() {
	if name, err :=  os.Hostname(); err == nil {
		hostname = name
	}

	if u, err := user.Current(); err == nil {
		cuser = u
	}
}


func (it *input) Read(p []byte) (n int, err error){
	if it.i >= len(it.line) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		p = []byte{it.line[it.i]}
	} else {
		p[0] = it.line[it.i]
	}
	it.i += 1
	time.Sleep(delay)
	return 1, nil
}

func write(question, answer string) {
	io.WriteString(os.Stdout, fmt.Sprintf("%s@%s:", cuser.Username, hostname))
	io.Copy(os.Stdin, &input{line : question})
	io.WriteString(os.Stdout, answer)

}

func main() {

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(){
		terminal.Restore(fd, oldState)
	}()

	write("who am i\n", "geek\n")
	write("really?\n", "of course\n")
}

