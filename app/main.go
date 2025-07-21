package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type entry struct {
	val string
}

var database map[string]entry = make(map[string]entry)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		if text == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if text == "ECHO" {
			scanner.Scan()
			encode_length := scanner.Text()
			scanner.Scan()
			echo := scanner.Text()
			conn.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", encode_length, echo)))
		} else if text == "GET" {
			scanner.Scan()
			_ = scanner.Text()
			scanner.Scan()
			key := scanner.Text()
			if entry, ok := database[key]; ok {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(entry.val), entry.val)))
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
			}
		} else if text == "SET" {
			scanner.Scan()
			_ = scanner.Text()
			scanner.Scan()
			key := scanner.Text()
			scanner.Scan()
			_ = scanner.Text()
			scanner.Scan()
			val := scanner.Text()

			database[key] = entry{ val }
			conn.Write([]byte(fmt.Sprintf("+OK\r\n")))

		}
	}
}
