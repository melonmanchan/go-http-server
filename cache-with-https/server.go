package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
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

	cer, err := tls.LoadX509KeyPair("./cache-with-https/server.crt", "./cache-with-https/server.key")

	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{
		Certificates:             []tls.Certificate{cer},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", *port), cfg)

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

	fmt.Fprint(conn, "\r\n")
	conn.Write(dat)
}
