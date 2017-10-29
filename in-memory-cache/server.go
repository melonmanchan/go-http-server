package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/melonmanchan/go-http-server/mimetypes"
	"github.com/melonmanchan/go-http-server/statuscodes"
	"github.com/melonmanchan/go-http-server/syncmap"
)

var port = flag.Int("port", 8081, "port to run under")
var basePath = flag.String("path", ".", "base path")
var CR = byte('\r')

var fileMap = syncmap.NewSyncByteMap()

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

func getPathFromHeader(header string) (string, *statuscodes.HTTPStatus) {
	paths := strings.Split(header, " ")

	if paths[0] != "GET" {
		return "", &statuscodes.MethodNotAllowed
	}

	return paths[1], nil
}

func safePath(reqPath string) string {
	return "./" + path.Join(*basePath, strings.Replace(reqPath, "..", "", -1))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	firstHeader, err := reader.ReadString(CR)

	if err != nil {
		fmt.Fprint(conn, statuscodes.ServerError.ToHeader())
		return
	}

	path, possibleError := getPathFromHeader(firstHeader)

	if possibleError != nil {
		fmt.Fprint(conn, possibleError.ToHeader())
	}

	sanitizedPath := safePath(path)

	info, err := os.Stat(sanitizedPath)

	if err != nil {
		fmt.Fprint(conn, statuscodes.NotFound.ToHeader())
		return
	}

	fileKey := sanitizedPath + string(info.ModTime().Unix())

	dat := []byte{}

	dat, ok := fileMap.Load(fileKey)

	if !ok {
		log.Printf("Cache empty! key %s", fileKey)
		dat, err = ioutil.ReadFile(sanitizedPath)
		fileMap.Store(fileKey, dat)
	}

	fileType := mimetypes.GetContentType(sanitizedPath)

	fmt.Fprint(conn, statuscodes.Ok.ToHeader())
	fmt.Fprintf(conn, "Content-Type: %s\r\n", fileType)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(dat))
	fmt.Fprint(conn, "\r\n")
	conn.Write(dat)
}
