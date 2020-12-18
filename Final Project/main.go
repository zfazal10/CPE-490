package main

import (
	"log"
	"net"
)

//initialize a TCP listener to look out for any new messages
func main() {
	s := newServer()
	go s.run() //server proccesses the messages and commands in a centralized manner so that messages arrive in order

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection failed: %s", err.Error())
			continue
		}

		go s.newClient(conn) //seperate go routine for each client
	}
}
