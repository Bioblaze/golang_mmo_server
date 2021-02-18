package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	connHost = "localhost"
	connPort = "6969"
	connType = "tcp"
)

func main() {
	go server()
}

func server() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	/*
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c // This will block until you manually exists with CRl-C
	*/

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			continue
		}

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		clientRequest := scanner.Text()

		switch scanner.Err() {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			reqParsed := strings.Fields(clientRequest)
			switch reqParsed[0] {
			case "test":
				log.Println("testing")
				return
			case "ping":
				_, err := conn.Write([]byte("pong\n"))
				if err != nil {
					log.Printf("failed to respond to client: %v\n", err)
				}
				return
			default:
				log.Println(clientRequest)
				return
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", scanner.Err())
			return
		}
	}
}
