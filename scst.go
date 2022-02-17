package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type Component interface {
	Start()
	Stop()
	Status() string
}

func NewServer() Server {
	s := new(Server)
	return *s
}

func NewClient() Client {
	c := new(Client)
	return *c
}

func NewStatusComp() Status {
	st := new(Status)
	return *st
}

func (s Server) Start() {
	srv, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Error while listening: ", err)
	}
	defer srv.Close()
	fmt.Println("Listening on localhost:8081...")
	for {
		conn, err := srv.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected.")
		go func(net.Conn) {
			buf := make([]byte, 1024)
			msL, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading: ", err)
			}
			fmt.Println("Server received: ", string(buf[:msL]))
			_, err = conn.Write([]byte("Server got: " + string(buf[:msL])))
			conn.Close()
		}(conn)
	}
}

func (s Server) Stop() {

}

func (s Server) Status() string {
	_, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		return "NOT RUNNING"
	}
	return "RUNNING"
}

func (c Client) Start() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	_, err = conn.Write([]byte("Client connected."))
	if err != nil {
		log.Println("Error writing: ", err)
	}
	fmt.Println("Waiting...")
	fmt.Scanln()
}

func (c Client) Stop() {

}

func (st Status) Start() {
	isRun := false
	_ = isRun
	for _, c := range *components {
		if c.Status() == "RUNNING" {
			isRun = true
		}
		fmt.Println(isRun)
	}
}

func (st Status) Stop() {

}

func (st Status) Status() string {
	return "ALL OK"
}

func (c Client) Status() string {
	return "RUNNING"
}

type Server struct{}
type Client struct{}
type Status struct{}

var components *[]Component

func main() {
	s := NewServer()
	c := NewClient()
	st := NewStatusComp()
	coms := []Component{s, c, st}
	*components = coms
	var wg sync.WaitGroup
	for i, c := range *components {
		wg.Add(1)
		comp := c
		go func() {
			defer wg.Done()
			comp.Start()
		}()
		fmt.Println(i + 1)
	}
	wg.Wait()
}
