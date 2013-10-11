package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
)

var addr = flag.String("addr", ":63837", "local address to listen for connections on")

func handleConnection(conn net.Conn) {
	// make sure we cleanup if anything happens
	defer conn.Close()

	// bufio.Reader leaves delimiter in slice,
	// bufio.Scanner does not
	r := bufio.NewReader(conn)
	for {
		b, err := r.ReadSlice('\n')
		if err == io.EOF {
			// all ended well
			return
		} else if err != nil {
			// failed at reading, kill goroutine
			log.Println("error reading:", err)
			return
		}

		// reverse the slice except for the '\n' at the end
		for i, n := 0, len(b)-2; i < (n+1)/2; i++ {
			b[i], b[n-i] = b[n-i], b[i]
		}

		_, err = conn.Write(b)
		if err != nil {
			// failed at writing, kill goroutine
			log.Println("error writing:", err)
			return
		}
	}
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		// can't listen on port -> kill program
		log.Fatalln("error initializing listener:", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// failed accepting connection -> wait for next
			log.Println("error establishing connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
