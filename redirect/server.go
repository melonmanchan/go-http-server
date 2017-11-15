package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/melonmanchan/go-http-server/common"
	"github.com/melonmanchan/go-http-server/statuscodes"
)

func main() {
	ln, err := net.Listen("tcp", ":80")

	defer ln.Close()

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func findHostFromHeaders(headers []string) string {
	for _, v := range headers {
		if strings.HasPrefix(v, "Host") {
			return strings.Split(v, " ")[1]
		}
	}
	return ""
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	headers := common.ReadAllHeaders(*reader)

	host := findHostFromHeaders(headers)

	fmt.Fprint(conn, statuscodes.Redirect.ToHeader())
	fmt.Fprintf(conn, "Location: https://%s\r\n", host)
}
