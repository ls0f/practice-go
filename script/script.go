package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"log"
	"os/exec"
	"github.com/kr/pty"
	"io"
	"flag"
	"time"
	"fmt"
)

var (
	startTime time.Time
)

type wrapf struct {
	rw  io.ReadWriter
	logF *os.File
	timeF *os.File
}

func (w *wrapf) Read(p []byte)(n int, err error){
	n, err = w.rw.Read(p)
	if n > 0 {
		now := time.Now()
		w.logF.Write(p[:n])
		io.WriteString(w.timeF, fmt.Sprintf("%f %d\n", now.Sub(startTime).Seconds(), n))
		startTime = now
	}
	return
}

func (w *wrapf) Write(p []byte)(n int, err error){
	return w.rw.Write(p)
}


func main() {

	f1 := flag.String("s", "scriptfile", "script file")
	f2 := flag.String("t", "scripttime", "script time")
	if *f1 == "" || *f2 == "" {
		flag.Usage()
		os.Exit(1)
	}
	flag.Parse()
	logF, err := os.Create(*f1)
	if err != nil {
		log.Fatal(err)
	}
	defer logF.Close()

	timeF, err := os.Create(*f2)
	if err != nil {
		log.Fatal(err)
	}
	defer timeF.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(){
		terminal.Restore(fd, oldState)
	}()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	startTime = time.Now()

	c := exec.Command(shell)
	f, err := pty.Start(c)
	if err != nil {
		log.Println(err)
		return
	}

	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGWINCH)
	//go func() {
	//	for range ch {
	//		if err := pty.InheritSize(os.Stdin, f); err != nil {
	//			log.Printf("error resizing pty: %s", err)
	//		}
	//	}
	//}()

	defer func(){
		c.Wait()
		f.Close()
	}()
	wf := &wrapf{rw: f, logF: logF, timeF: timeF}
	go io.Copy(wf, os.Stdin)
	io.Copy(os.Stdout, wf)
}


