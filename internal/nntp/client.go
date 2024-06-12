package nntp

import (
	"net/textproto"
	"net"
		ClientOption func(*Client)
	)

type (
	Client struct {
		*textproto.Conn

		// Describes the max amount of concurrent connections to a given server
		maxConnections uint8
	}
)

func NewClient(options ...ClientOption) *Client {
	c := &Client{}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithHostAddr(addr string) ClientOption {
	return func(c *Client) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err) // Handle error appropriately in real code
		}
		c.Conn = textproto.NewConn(conn)
	}
}
