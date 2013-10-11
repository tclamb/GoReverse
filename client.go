package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var addr = flag.String("addr", ":63837", "remote address to connect to")

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	defer conn.Close()
	if err != nil {
		// can't connect -> kill program
		log.Fatal("error connecting to server: ", err)
	}

	go func() {
		// bind input to TCP connection
		_, err := io.Copy(conn, os.Stdin)
		if err != nil {
			// some error in output or connection
			log.Fatal("error copying input to server: ", err)
		}
		// close our end when done so the server knows to close its end
		conn.(*net.TCPConn).CloseWrite()
	}()

	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		// error while receiving from server
		log.Fatal("error receiving from server: ", err)
	}
}
