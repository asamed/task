package main

import (
	"fmt"
	"net"
	"os"
)

type Controler interface {
	start() error
	stop() error
	status() (string, error)
}

func processClient(conn net.Conn) {
	buf := make([]byte, 1024)
	msLn, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err)
	}
	fmt.Println("Server received: ", string(buf[:msLn]))
	_, err = conn.Write([]byte("Got the message: " + string(buf[:msLn])))
	conn.Close()
}

func serverStart() {
	fmt.Println("Starting server...")
	server, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Error while listening: ", err)
	}
	defer server.Close()
	fmt.Println("Listening on localhost:8081")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected.")
		go processClient(conn)
	}
}

func main() {
	// go serverStart()
	// conn, err := net.Dial("tcp", "localhost:8081")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err = conn.Write([]byte("Hello server."))
	// buf := make([]byte, 1024)
	// mLen, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println("Error reading: ", err)
	// }
	// fmt.Println("Client received: ", string(buf[:mLen]))
	// defer conn.Close()
}
