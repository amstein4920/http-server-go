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

const (
	notFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
)

type ResponseItems struct {
	contentType     string
	contentLength   int
	content         string
	contentEncoding string
}

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
	dPath := os.Args[slices.Index(os.Args, "--directory")+1]

	r := ResponseItems{}

	if req.Header.Get("Accept-Encoding") == strings.ToLower("gzip") {
		r.contentEncoding = "\r\nContent-Encoding: gzip"
	}

	if path == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	}
	if strings.HasPrefix(path, "/echo/") {
		r.contentType = "text/plain"
		r.content = path[6:]
		r.contentLength = len(r.content)
		return []byte(createResponseString(r))
	}
	if strings.HasPrefix(path, "/user-agent") {
		r.contentType = "text/plain"
		r.content = req.UserAgent()
		r.contentLength = len(req.UserAgent())
		return []byte(createResponseString(r))
	}
	if strings.HasPrefix(path, "/files/") {
		fmt.Println(req.Method)
		if req.Method == "GET" {
			return filesGet(path, dPath, r)
		}
		if req.Method == "POST" {
			return filesPost(req, dPath)
		}
	}
	return []byte(notFoundResponse)
}

func filesGet(path string, dPath string, r ResponseItems) []byte {
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

	r.contentType = "application/octet-stream"
	r.content = string(content)
	r.contentLength = int(size)

	return []byte(createResponseString(r))
}

func filesPost(req *http.Request, dPath string) []byte {
	filePath := dPath + req.URL.Path[6:]
	f, err := os.Create(filePath)
	if err != nil {
		return []byte(notFoundResponse)
	}
	body, _ := io.ReadAll(req.Body)
	_, err = f.Write(body)
	if err != nil {
		return []byte(notFoundResponse)
	}
	return []byte("HTTP/1.1 201 Created\r\n\r\n")
}

func createResponseString(r ResponseItems) string {
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length:%d%s \r\n\r\n%s", r.contentType, r.contentLength,
		r.contentEncoding, r.content)
}
