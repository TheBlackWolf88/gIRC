package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("An error has occured: %s\n", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("An error has occured: %s\n", err)
		}

		go handleConn(conn, "msg")
	}
}

func handleConn(conn net.Conn, prefix string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			return
		}

		line := fmt.Sprintf("%s %s", prefix, bytes)
		conn.Write([]byte(line))

	}
}
