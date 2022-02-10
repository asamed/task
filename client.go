package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write([]byte("Greetings server."))
	buf := make([]byte, 1024)
	mLn, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err)
	}
	fmt.Println("Client received: ", string(buf[:mLn]))
	defer conn.Close()
}
