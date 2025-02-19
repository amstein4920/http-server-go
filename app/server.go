package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Println("Unable to read connection input")
		return
	}

	conn.Write(writeResponse(req))
}

func writeResponse(req *http.Request) []byte {

	path := req.URL.Path

	if path == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	}
	if strings.HasPrefix(path, "/echo/") {
		content := path[6:]
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%d \r\n\r\n%s", len(content), content))
	}
	if strings.HasPrefix(path, "/user-agent") {
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%d \r\n\r\n%s", len(req.UserAgent()), req.UserAgent()))
	}
	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}
