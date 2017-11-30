package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
)

func main(){
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn){
	scanner := bufio.NewScanner(conn)
	for scanner.Scan(){
		ln := scanner.Text()
		fmt.Println(ln)
		fmt.Fprintf(conn, "Enter some text: %s\n",ln)
	}
	defer conn.Close()
}