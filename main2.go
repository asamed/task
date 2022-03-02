package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const socket = "sock"

type Component interface {
	Start()
	Stop()
	Status() string
}

type Server struct {
	l   net.Listener
	run bool
}

type Client struct {
	con net.Conn
	run bool
}

type Status struct{}

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

func (s *Server) Start() {
	if s.l != nil {
		fmt.Println("Server already running.")
		return
	}
	var err error
	s.l, err = net.Listen("unix", socket)
	if err != nil {
		fmt.Println("Error while listening: ", err)
		return
	}
	s.run = true
	defer s.l.Close()
	fmt.Println("Listening on sock...")
	for {
		conn, err := s.l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err)
			return
		}
		fmt.Println("Connected.")
		go func(net.Conn) {
			buf := make([]byte, 1024)
			msL, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading: ", err)
			}
			fmt.Println("Server received: ", string(buf[:msL]))
			_, err = conn.Write([]byte("Server got: " + string(buf[:msL])))
			if err != nil {
				fmt.Println("Error writing: ", err)
			}
			conn.Close()
		}(conn)
	}
}

func (s *Server) Stop() {
	s.run = false
	s.l.Close()
	os.Remove(socket)
}

func (s Server) Status() string {
	_, err := net.Dial("unix", socket)
	if err != nil {
		return "NOT RUNNING"
	}
	return "RUNNING"
}

func (c *Client) Start() {
	var err error
	for {
		c.con, err = net.Dial("unix", socket)
		if err != nil {
			fmt.Println("Error connecting: ", err)
			fmt.Println("Trying again...")
			time.Sleep(time.Second)
		}
		if err == nil {
			break
		}
	}
	c.run = true
	_, err = c.con.Write([]byte("Client connected."))
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

func (c *Client) Stop() {
	c.run = false
	c.con.Close()
}

func (st Status) Start() {
	fmt.Println("Server is ", s.Status())
	fmt.Println("Client is ", c.Status())
	fmt.Println("Status is ", st.Status())
}

func (st Status) Stop() {

}

func (st Status) Status() string {
	return "RUNNING"
}

func (c *Client) Status() string {
	if c.con != nil {
		return "RUNNING"
	}
	return "NOT RUNNING"
}

func cleanup() {
	s.l.Close()
	c.con.Close()
}

var s, c, st = NewServer(), NewClient(), NewStatusComp()

func main() {
	os.Remove(socket)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Printf("\nCleaning...")
		cleanup()
		os.Exit(1)
	}()
	var chc int
	var wg sync.WaitGroup
	for {
		time.Sleep(time.Second)
		fmt.Println("Choose action: ")
		fmt.Println("1. Start server")
		fmt.Println("2. Start client")
		fmt.Println("3. Stop server")
		fmt.Println("4. Stop client")
		fmt.Println("5. Server status")
		fmt.Println("6. Client status")
		fmt.Println("7. All status check")
		fmt.Println("8. Exit")
		fmt.Printf("Enter your choice: ")
		fmt.Scanln(&chc)
		switch chc {
		case 1:
			wg.Add(1)
			go s.Start()
		case 2:
			wg.Add(1)
			go c.Start()
		case 3:
			if !s.run {
				fmt.Println("Server is not running, nothing to stop.")
				continue
			}
			wg.Done()
			s.Stop()
		case 4:
			if !c.run {
				fmt.Println("Client is not running, nothing to stop.")
				continue
			}
			wg.Done()
			c.Stop()
		case 5:
			fmt.Println(s.Status())
		case 6:
			fmt.Println(c.Status())
		case 7:
			st.Start()
		case 8:
			cleanup()
			os.Exit(1)
		}
	}
	wg.Wait()
}
