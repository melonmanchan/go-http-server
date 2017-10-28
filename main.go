package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path"
	"strings"
)

type HTTPStatus struct {
	status int
	msg    string
}

func (h HTTPStatus) ToHeader() string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", h.status, h.msg)
}

var notFound = HTTPStatus{404, "not found"}
var ok = HTTPStatus{200, "OK"}
var methodNotAllowed = HTTPStatus{405, "method not allowed"}
var serverError = HTTPStatus{500, "internal server error"}

func main() {
	ln, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

var filetypes = map[string]string{
	".html": "text/html",
	".htm":  "text/html",
	".js":   "application/javascript",
	".css":  "text/css",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	"":      "text/plain",
}

func getPathFromBytes(bytes []byte) (string, *HTTPStatus) {
	s := string(bytes[:])
	splat := strings.Split(s, "\r\n")
	paths := strings.Split(splat[0], " ")

	if paths[0] != "GET" {
		return "", &methodNotAllowed
	}

	return paths[1], nil
}

func safePath(reqPath string) string {
	return path.Clean(strings.Replace(reqPath, "..", "", -1))
}

func readFile(reqPath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(reqPath)
	return dat, err
}

func getContentType(fileName string) string {
	name := path.Ext(fileName)
	return filetypes[name]
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)

	if err != nil {
		fmt.Fprint(conn, serverError.ToHeader())
		return
	}

	path, possibleError := getPathFromBytes(buf)

	if possibleError != nil {
		fmt.Fprint(conn, possibleError.ToHeader())
	}

	sanitizedPath := safePath(path)

	dat, err := readFile("." + sanitizedPath)

	if err != nil {
		fmt.Fprint(conn, notFound.ToHeader())
		return
	}

	fileType := getContentType(sanitizedPath)

	fmt.Fprint(conn, ok.ToHeader())
	fmt.Fprintf(conn, "Content-Type: %s\r\n", fileType)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(dat))
	fmt.Fprint(conn, "\r\n")
	conn.Write(dat)
}
