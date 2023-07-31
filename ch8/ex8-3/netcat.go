package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8008")
	if err != nil {
		log.Print(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			return
		}
		log.Println("done")
		done <- struct{}{}
	}()
	mustCopy(conn, os.Stdin)
	err = conn.CloseWrite()
	if err != nil {
		return
	}
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}
