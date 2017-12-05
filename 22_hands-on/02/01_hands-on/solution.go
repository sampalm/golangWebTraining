package main

import (
	"io"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "I see you connected.")
}
