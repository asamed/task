package main

import (
	"fmt"
	"log"
	"net"
)

type Controller interface {
	start() (net.Conn, error)
	stop(con net.Conn) error
	status() (string, error)
}

type Connection struct{}

func (c Connection) start() (net.Conn, error) {
	return net.Dial("tcp", "localhost:8081")
}

func (c Connection) stop(con net.Conn) error {
	err := con.Close()
	return err
}

func main() {
	var c Connection
	conn, err := c.start()
	if err != nil {
		log.Fatal(err)
	}
	var msg, ch string
	fmt.Println("A message for the server?")
	for {
		fmt.Scanln(&msg)
		_, err = conn.Write([]byte(msg))
		fmt.Println("Anything else?")
		fmt.Scanln(&ch)
		if ch == "No" || ch == "no" || ch == "NO" {
			break
		}
	}
	fmt.Println("Close the connection?")
	for {
		var ch string
		fmt.Scanln(&ch)
		if ch == "Yes" || ch == "yes" || ch == "y" {
			err = c.stop(conn)
			break
		}
	}
}
