package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	i := 0
	var m, u string
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			m = strings.Fields(ln)[0]
			fmt.Println("***** METHOD ", m)
			u = strings.Fields(ln)[1]
			fmt.Println("***** URL ", u)
		}
		if ln == "" {
			break
		}
		i++
	}
	defer conn.Close()

	switch {
	case m == "GET" && u == "/":
		index(conn)
	case m == "GET" && u == "/apply":
		apply(conn)
	case m == "POST" && u == "/apply":
		applyPost(conn)
	default:
		index(conn)
	}
}

func index(conn net.Conn) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8">
				<title>Index</title>
			</head>
			<body>
				<h1>INDEX PAGE</h1>
				<a href="/apply">apply</a>
				<a href="/">index</a>
			</body>
		</html>
	`

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, body)
}

func apply(conn net.Conn) {
	body := `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8">
			<title>Apply</title>
		</head>
		<body>
			<h1>APPLY PAGE</h1>
			<a href="/apply">apply</a>
			<a href="/">index</a>
			<form action="/apply" method="POST">
				<input type="hidden" value="Can u see me ?">
				<button type="submit">CLICK ME</button>
			</form>
		</body>
	</html>
	`

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, body)

}

func applyPost(conn net.Conn) {

	body := `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8">
			<title>Apply</title>
		</head>
		<body>
			<h1>APPLY POST PAGE</h1>
			<a href="/">index</a>
			<a href="/apply">apply</a>
		</body>
	</html>
	`

	io.WriteString(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, body)
}
