package mq

import (
	"fmt"
	"net"
)

type conn struct {
	id int
	net.Conn
}

func newconn(id int, nc net.Conn) *conn {
	return &conn{id: id, Conn: nc}
}

func (c *conn) String() string {
	return fmt.Sprintf("%s:%d", c.RemoteAddr().String(), c.id)
}
