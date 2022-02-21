package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
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
	// syscall.Unlink("/tmp/go.sock")
	srv, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Error while listening: ", err)
		return
	}
	s.l = srv
	defer s.l.Close()
	fmt.Println("Listening on localhost:8081...")
	for {
		conn, err := s.l.Accept()
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
	s.l.Close()
}

func (s Server) Status() string {
	_, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		return "NOT RUNNING"
	}
	return "RUNNING"
}

func (c Client) Start() {
	var err error
	for i := 0; i < 3; i++ {
		c.con, err = net.Dial("tcp", "localhost:8081")
		if err != nil {
			fmt.Println("Error connecting: ", err)
			fmt.Println("Trying again...")
			time.Sleep(time.Millisecond * 1500)
		}
		if err == nil {
			break
		}
	}
	// conn, err := net.Dial("tcp", "localhost:8081")
	// if err != nil {
	// 	log.Fatal("Error connecting: ", err)
	// }
	_, err = c.con.Write([]byte("Client connected."))
	*cst = true
	if err != nil {
		log.Println("Error writing: ", err)
	}
	buf := make([]byte, 1024)
	msr, err := c.con.Read(buf)
	if err != nil {
		fmt.Println("Error reading from server: ", err)
	}
	fmt.Println("Client received: ", string(buf[:msr]))
}

func (c Client) Stop() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	c.con = conn
	c.con.Close()

}

func (st Status) Start() {
	isRun := false
	_ = isRun
	for i, c := range *components {
		if c.Status() == "RUNNING" {
			isRun = true
		}
		fmt.Println("Component", i+1, "running is", isRun)
	}
}

func (st Status) Stop() {

}

func (st Status) Status() string {
	return "RUNNING"
}

func (c Client) Status() string {
	if *cst == true {
		return "RUNNING"
	}
	return ""
}

type Server struct {
	l net.Listener
}
type Client struct {
	con net.Conn
}
type Status struct{}

var components *[]Component
var cst *bool

func main() {
	s := NewServer()
	c := NewClient()
	st := NewStatusComp()
	coms := []Component{s, c, st}
	components = &coms
	cf := false
	cst = &cf
	var wg sync.WaitGroup
	for _, c := range *components {
		wg.Add(1)
		comp := c
		go func() {
			defer wg.Done()
			comp.Start()
		}()
	}
	wg.Wait()
}
