package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main(){
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}
	defer li.Close()

	conn, err := li.Accept()
	if err != nil {
		log.Println(err)
	}

	io.WriteString(conn, "\n Hello from TPC Server \n")
	fmt.Fprintln(conn, "Testing connection...")
	fmt.Fprintf(conn, "%v", "1.. 2.. 3..")

	conn.Close()
}