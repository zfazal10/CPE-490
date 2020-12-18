package main

import (
	"net"
)

type room struct {
	name    string               //name of room
	members map[net.Addr]*client //client remove address is used as their key
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}
