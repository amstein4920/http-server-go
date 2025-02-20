package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
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
	notFoundResponse := "HTTP/1.1 404 Not Found\r\n\r\n"

	path := req.URL.Path
	dPath := os.Args[slices.Index(os.Args, "--directory")+1]

	if path == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	}
	if strings.HasPrefix(path, "/echo/") {
		content := path[6:]
		return []byte(createResponseString("text/plain", len(content), content))
	}
	if strings.HasPrefix(path, "/user-agent") {
		return []byte(createResponseString("text/plain", len(req.UserAgent()), req.UserAgent()))
	}
	if strings.HasPrefix(path, "/files/") {
		filePath := dPath + path[6:]
		info, err := os.Stat(filePath)
		if err != nil {
			return []byte(notFoundResponse)
		}
		size := info.Size()

		f, err := os.Open(filePath)
		if err != nil {
			return []byte(notFoundResponse)
		}
		content, err := io.ReadAll(f)
		if err != nil {
			return []byte(notFoundResponse)
		}

		return []byte(createResponseString("application/octet-stream", int(size), string(content)))
	}
	return []byte(notFoundResponse)
}

func createResponseString(contentType string, length int, content string) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length:%d \r\n\r\n%s", contentType, length, content)
}
