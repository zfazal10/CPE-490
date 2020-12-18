package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room //map of the rooms
	commands chan command     //channel for sending cpommands from client to server
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() { //blocking function to receive and process messages, called from main func in main.go
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_USERNAME:
			s.username(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("%s has joined", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		username: "anonymous", //default name if user does not specify
		commands: s.commands,
	}

	c.readInput()
}

//when a new client is connected, its initialized and starts listening for any incoming messages

func (s *server) username(c *client, username string) { //set client username and send confirmation message
	c.username = username
	c.msg(fmt.Sprintf("Username: %s", username))
}

func (s *server) join(c *client, roomName string) { //joins room, creates it if non existent
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c) //remove user from old rooms
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined", c.username))

	c.msg(fmt.Sprintf("You have joined %s", roomName))
}

func (s *server) listRooms(c *client) { //displays current rooms
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Open Rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) { //sends message to all other clients in current room
	msg := strings.Join(args[1:len(args)], " ")
	c.room.broadcast(c, c.username+": "+msg)
}

func (s *server) quit(c *client) { //closes connection, but not without saying goodbye first!
	log.Printf("%s has left the chat", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("Goodbye")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.username))
	}
}
