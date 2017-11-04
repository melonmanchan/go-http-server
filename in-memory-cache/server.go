package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/melonmanchan/go-http-server/common"
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

func gzipBytes(inputBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)

	if _, err := gz.Write(inputBytes); err != nil {
		return nil, err
	}

	if err := gz.Flush(); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
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

	info, err := os.Stat(sanitizedPath)

	if err != nil {
		fmt.Fprint(conn, statuscodes.NotFound.ToHeader())
		return
	}

	fileKey := sanitizedPath + info.ModTime().String()

	dat := []byte{}

	dat, ok := fileMap.Load(fileKey)

	if !ok {
		log.Printf("Cache empty! key %s", fileKey)
		contents, _ := ioutil.ReadFile(sanitizedPath)
		dat, _ = gzipBytes(contents)
		fileMap.Store(fileKey, dat)
	}

	fileType := mimetypes.GetContentType(sanitizedPath)

	fmt.Fprint(conn, statuscodes.Ok.ToHeader())
	fmt.Fprintf(conn, "Content-Type: %s\r\n", fileType)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(dat))
	fmt.Fprint(conn, "Content-Encoding: gzip\r\n")
	fmt.Fprintf(conn, "Last-Modified: %s\r\n", info.ModTime().Format(http.TimeFormat))

	fmt.Fprint(conn, "\r\n")
	conn.Write(dat)
}
