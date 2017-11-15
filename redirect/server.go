package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	headers := common.ReadAllHeaders(*reader)

	host := common.FindValueFromHeaders(headers, "Host")

	fmt.Fprint(conn, statuscodes.Redirect.ToHeader())
	fmt.Fprintf(conn, "Location: https://%s\r\n", host)
}
