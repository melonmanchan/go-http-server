package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path"
	"strings"

	"github.com/melonmanchan/go-http-server/mimetypes"
	"github.com/melonmanchan/go-http-server/statuscodes"
)

var port = flag.Int("port", 8081, "port to run under")
var basePath = flag.String("path", ".", "base path")

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening at port :%d", *port)

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func getPathFromBytes(bytes []byte) (string, *statuscodes.HTTPStatus) {
	s := string(bytes[:])
	splat := strings.Split(s, "\r\n")
	paths := strings.Split(splat[0], " ")

	if paths[0] != "GET" {
		return "", &statuscodes.MethodNotAllowed
	}

	return paths[1], nil
}

func safePath(reqPath string) string {
	return "./" + path.Join(*basePath, strings.Replace(reqPath, "..", "", -1))
}

func readFile(reqPath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(reqPath)
	return dat, err
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)

	if err != nil {
		fmt.Fprint(conn, statuscodes.ServerError.ToHeader())
		return
	}

	path, possibleError := getPathFromBytes(buf)

	if possibleError != nil {
		fmt.Fprint(conn, possibleError.ToHeader())
	}

	sanitizedPath := safePath(path)

	dat, err := readFile(sanitizedPath)

	if err != nil {
		fmt.Fprint(conn, statuscodes.NotFound.ToHeader())
		return
	}

	fileType := mimetypes.GetContentType(sanitizedPath)

	fmt.Fprint(conn, statuscodes.Ok.ToHeader())
	fmt.Fprintf(conn, "Content-Type: %s\r\n", fileType)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(dat))
	fmt.Fprint(conn, "\r\n")
	conn.Write(dat)
}
