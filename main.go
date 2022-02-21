package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {

	listner, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listner.Accept()

		if err != nil {
			log.Fatalln(err)
		}

		go Handle(conn)
	}

}

func Handle(conn net.Conn) {
	defer conn.Close()
	url, method := RequestHandler(conn)

	//url multiplexer
	switch url {
	case "/":
		ResponseHome(conn, method)
	case "/about":
		ResponseAbout(conn, method)
	default:
		ResponseNotFound(conn, method)
	}

}

func RequestHandler(conn net.Conn) (string, string) {
	var url, method string
	firstLine := true
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		if firstLine == true {
			method = strings.Fields(ln)[0]
			url = strings.Fields(ln)[1]
			log.Println("REQUESTED METHOD >>", method)
			log.Println("REQUESTED URL >>", url)
			firstLine = false
		}
		if ln == "" {
			log.Println("---header end---")
			break
		}
	}
	return url, method
}

func ResponseHome(conn net.Conn, method string) {

	f, err := os.ReadFile("templates/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(len(f))

	//HEADER according to rfc 7230 IETF///////////////////////////////////
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "DATE: %s\r\n", time.Now().Format(time.RFC1123))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "Accept-Ranges: bytes\r\n")
	fmt.Fprintf(conn, "Content-length: %d\r\n", len(f))
	fmt.Fprint(conn, "\r\n") // this is important according to rfc 7230 IETF
	/////////////////////////////////////////
	//BODY//
	fmt.Fprint(conn, string(f))

}

func ResponseAbout(conn net.Conn, method string) {

	f, err := os.ReadFile("templates/about.html")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(len(f))

	//HEADER according to rfc 7230 IETF///////////////////////////////////
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "DATE: %s\r\n", time.Now().Format(time.RFC1123))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "Accept-Ranges: bytes\r\n")
	fmt.Fprintf(conn, "Content-length: %d\r\n", len(f))
	fmt.Fprint(conn, "\r\n") // this is important according to rfc 7230 IETF
	/////////////////////////////////////////
	//BODY//
	fmt.Fprint(conn, string(f))

}

func ResponseNotFound(conn net.Conn, method string) {

	f, err := os.ReadFile("templates/notfound.html")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(len(f))
	fmt.Println(method)

	//HEADER according to rfc 7230 IETF///////////////////////////////////
	fmt.Fprint(conn, "HTTP/1.1 404 NOT FOUND\r\n")
	fmt.Fprintf(conn, "DATE: %s\r\n", time.Now().Format(time.RFC1123))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "Accept-Ranges: bytes\r\n")
	fmt.Fprintf(conn, "Content-length: %d\r\n", len(f))
	fmt.Fprint(conn, "\r\n") // this is important according to rfc 7230 IETF
	/////////////////////////////////////////
	//BODY//
	fmt.Fprint(conn, string(f))

}
