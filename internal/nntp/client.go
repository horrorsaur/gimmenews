package nntp

import (
	"log"
	"net"
	"net/textproto"
)

type (
	// A generic nntp client that relies on a textproto.Conn for handling
	Client struct {
		*textproto.Conn

		log *log.Logger

		// Describes the max amount of concurrent connections to a given server
		maxConnections uint8
	}

	usenetOpts struct {
		capabilities []string
	}
)

func NewClient(host, port string, clogger *log.Logger) *Client {
	addr := net.JoinHostPort(host, port)
	clogger.Printf("attempting connection to '%s'...", addr)
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		clogger.Fatal(err)
	}

	code, msg, err := conn.ReadCodeLine(200)
	if err != nil {
		clogger.Fatal(err)
	}

	if code != STATUS_INITIAL_OK {
		clogger.Fatal(err)
	}

	clogger.Printf("%d %s", code, msg)
	return &Client{
		Conn: conn,
		log:  clogger,
	}
}

// Prompts the client to send a 'CAPABILITIES' command.
//
// Returns a list of capabilities that the server supports
func (c *Client) GetCapabilities() []string {
	resp, err := c.Command(STATUS_CAPABILITIES, "CAPABILITIES")
	if err != nil {
		c.log.Fatalf("error receiving capabilities response: %s", err)
	}

	c.log.Printf("Response: %s", resp)
	return nil
}

func (c *Client) Command(expectedCode int, format string, args ...any) ([]string, error) {
	c.log.Printf("sending command '%s' with args '%s'", format, args)
	if err := c.PrintfLine("CAPABILITIES"); err != nil {
		c.log.Fatal("error sending cmd: ", err)
		return nil, err
	}
	return c.parseResponse(expectedCode)
}

// If the response is multi-line, ReadCodeLine returns an error.
//
// An expectCode <= 0 disables the check of the status code.
func (c *Client) parseResponse(expectedCode int) ([]string, error) {
	_, _, err := c.ReadCodeLine(expectedCode)
	if err != nil {
		c.log.Print("error parsing response: ", err)
		return nil, err
	}
	return c.ReadDotLines()
}

func (c *Client) parseMulti() ([]string, error) {
	return nil, nil
}
