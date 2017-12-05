package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}

		go serve(conn)
	}
}

func serve(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if ln == "" {
			break
		}
		fmt.Fprintf(conn, "ResponseWriter: %s\n", ln)
	}
	defer conn.Close()

	io.WriteString(conn, "I see you connected.")
}
