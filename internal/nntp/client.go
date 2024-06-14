package nntp

import (
	"crypto/tls"
	"log"
	"net"
	"net/textproto"
	"strconv"
)

// We mainly care for
// AUTHINFO USER(?)
// HEADERS(?)
// NEWSGROUPS(?)
// SUBSCRIPTIONS(?)
// READER
// STARTTLS

// TODO:
// Add support for:
// POST
// Mode Switching https://datatracker.ietf.org/doc/html/rfc3977#section-3.4.2

type (
	// A generic nntp client that relies on a textproto.Conn for handling
	Client struct {
		*textproto.Conn

		log *log.Logger

		// TODO: Concurrent connections
		maxConnections uint8
		// Describes the capabilities of the server state. This can change based on server state (e.g., preauth -> auth)
		capabilities []string
	}
)

func NewClient(host, port string, insecure *bool, clogger *log.Logger) *Client {
	addr := net.JoinHostPort(host, strconv.Itoa(DEFAULT_TLS_PORT))

	if *insecure {
		clogger.Printf("WARNING. connecting on unencrypted port!")
		addr = net.JoinHostPort(host, port)
	}

	clogger.Printf("attempting connection to '%s'...", addr)
	conn, err := textproto.Dial("tcp", addr)
	if err != nil {
		clogger.Fatal(err)
	}

	// tls.Dial("tcp", addr, &tls.Config{})

	code, msg, err := conn.ReadCodeLine(STATUS_AVAILABLE_POSTING_ALLOWED)
	if err != nil {
		clogger.Fatal(err)
	}

	if code != STATUS_AVAILABLE_POSTING_ALLOWED {
		clogger.Fatal(err)
	}

	clogger.Printf("GOT %d %s", code, msg)
	return &Client{
		Conn: conn,
		log:  clogger,
	}
}

// Prompts the client to send a 'CAPABILITIES' command.
//
// Returns a list of capabilities that the server supports
func (c *Client) GetCapabilities() []string {
	resp, err := c.sendCommand(STATUS_CAPABILITIES, "CAPABILITIES")
	if err != nil {
		c.log.Fatalf("error receiving capabilities response: %s", err)
	}
	c.log.Printf("Response: %s", resp)
	return nil
}

// Indicating capability: LIST
// Syntax
//
//	LIST [keyword [wildmat|argument]]
func (c *Client) List() {
	resp, err := c.sendCommand(215, "LIST alt")
	if err != nil {
		c.log.Fatal(err)
	}
	c.log.Printf("Response: %s", resp)
}

// the IHAVE command is designed for transit, while the
//    commands indicated by the READER capability are designed for reading
//    clients.

// Generic send
func (c *Client) Send(command string, expectedStatusCode int) string {
	return ""
}

func (c *Client) sendCommand(expectedCode int, format string, args ...any) ([]string, error) {
	c.log.Printf("sending command '%s' with args '%s'", format, args)
	if err := c.PrintfLine(format, args...); err != nil {
		c.log.Fatal("error sending cmd: ", err)
		return nil, err
	}
	return c.parseResponse(expectedCode)
}

// If the response is multi-line, ReadCodeLine returns an error.
//
// An expectCode <= 0 disables the check of the status code.
func (c *Client) parseResponse(expectedCode int) ([]string, error) {
	code, msg, err := c.ReadCodeLine(expectedCode)
	if err != nil {
		c.log.Print("error parsing response: ", err)
		return nil, err
	}

	c.log.Printf("GOT %d %s", code, msg)
	return c.ReadDotLines()
}

func (c *Client) parseMulti() ([]string, error) {
	return nil, nil
}

// TODO:
// s := bufio.NewScanner(c.R)
// for s.Scan() {
// 	fmt.Printf(s.Text())
// }
