package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	L, err := net.Listen("tcp4", ":6379")

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println("Listening on port :6379")

	for {
		conn, err := L.Accept()

		if err != nil {
			fmt.Printf("%v", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		logged := io.TeeReader(conn, os.Stdout)
		
		resp := NewResp(logged)
		value, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}

			fmt.Println("Error:", err)
			break
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		
		writer := NewWriter(conn)

		command := strings.ToUpper(value.array[0].bulk)
		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		args := value.array[1:]

		result := handler(args)
		writer.Write(result)
	}
}