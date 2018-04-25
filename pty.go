package main

import (
	"flag"
	"github.com/kr/pty"
	"io"
	"net"
	"os"
	"os/exec"
	"log"
	"fmt"
)

type wrapf struct {
	rw  io.ReadWriter

	logF *os.File
	timeF *os.File
}

func (w *wrapf) Read(p []byte)(n int, err error){
	n, err = w.rw.Read(p)
	if n > 0 {
		w.logF.Write(p[:n])
	}
	return
}

func (w *wrapf) Write(p []byte)(n int, err error){
	return w.rw.Write(p)
}


func handler(conn net.Conn) {
	defer func(){
		log.Printf("conn: %v close", conn.RemoteAddr())
		conn.Close()
	}()
	c := exec.Command("bash")
	f, err := pty.Start(c)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(){
		f.Close()
		err := c.Wait()
		if err != nil {
			log.Print(err)
		}

	}()
	logF, err := os.Create(fmt.Sprintf("%s.txt", conn.RemoteAddr()))
	if err != nil {
		log.Print(err)
		return
	}
	w := &wrapf{rw: f, logF: logF}
	go func(){
		io.Copy(w, conn)
	}()
	io.Copy(conn, w)
	logF.Close()
}


func main() {

	addr := flag.String("a", "", "addr")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}

	ser, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Printf("conn: %s", conn.RemoteAddr())
		go handler(conn)
	}
}
