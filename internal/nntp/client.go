package nntp

import "net/textproto"

type (
	Client struct {
		*textproto.Conn

		// Describes the max amount of concurrent connections to a given server
		maxConnections uint8
	}
)

func NewClient() *Client {
	var c *Client

	return c
}
