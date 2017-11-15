package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/melonmanchan/go-http-server/common"
	"github.com/melonmanchan/go-http-server/mimetypes"
	"github.com/melonmanchan/go-http-server/statuscodes"
)

var port = flag.Int("port", 8081, "port to run under")
var basePath = flag.String("path", ".", "base path")
var CR = byte('\r')

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening at port :%d", *port)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	headers := common.ReadAllHeaders(*reader)

	path, possibleError := common.GetPathFromHeader(headers[0])

	if possibleError != nil {
		fmt.Fprint(conn, possibleError.ToHeader())
	}

	sanitizedPath := common.SafePath(path, *basePath)

	dat, err := ioutil.ReadFile(sanitizedPath)

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
