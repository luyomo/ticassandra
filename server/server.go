package server

import (
	"fmt"
//	"io"
	"net"
)

// TCPServer struct
type TCPServer struct {
	Bind string
	Port int
}

// Start TCPServer
func (s *TCPServer) Start() {
	fmt.Printf("started tcp echo server... ... \n")
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Bind, s.Port))
	defer ln.Close()
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
        fmt.Println("------------------")
		go func(conn net.Conn ) {
			defer conn.Close()
            for {
                msg := make([]byte, 1024)
                fmt.Println("Starting to collect data")
                _, err = conn.Read(msg)
                fmt.Println("Fetched all the data ")
                if err != nil {
                    panic(err)
                }
                fmt.Println(string(msg))
            }
		}(conn)
        fmt.Println(">>>>>>>>>>>")
	}
}
