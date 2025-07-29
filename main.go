package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	L, err := net.Listen("tcp4", ":6379")

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println("Listening on port :6379")

	conn, err := L.Accept()

	if err != nil {
		fmt.Printf("%v", err)
	}

	defer conn.Close()

	for {
		logged := io.TeeReader(conn, os.Stdout)
		
		resp := NewResp(logged)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}

}