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

		connInfo ServerConnectionInfo

		log *log.Logger

		// indicates whether we are using TLS on DEFAULT_TLS_PORT (563)
		secure bool
		// TODO: Concurrent connections
		maxConnections uint8
		// Describes the capabilities of the server state. This can change based on server state (e.g., preauth -> auth)
		capabilities []string
	}

	ServerConnectionInfo struct {
		ServerName string
		TLSVersion uint16
	}

	TLS struct {
		*tls.Conn
	}
)

func NewClient(host, port string, insecure *bool, clogger *log.Logger) *Client {
	addr := net.JoinHostPort(host, strconv.Itoa(DEFAULT_TLS_PORT))

	if *insecure {
		clogger.Printf(`
			================ WARNING ================
			Client has been flagged for an unencrypted connection!
			=========================================
		`)
		clogger.Printf("connecting to '%s'...", addr)

		conn, err := textproto.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			clogger.Fatal(err)
		}

		return &Client{
			Conn:   conn,
			secure: false,
			log:    clogger,
		}
	}

	clogger.Printf("connecting to '%s'... (TLS)", addr)
	tlsConf := &tls.Config{InsecureSkipVerify: false}
	conn, err := tls.Dial("tcp", addr, tlsConf)
	if err != nil {
		clogger.Fatal(err)
	}

	connState := conn.ConnectionState()
	clogger.Printf("Server Name: %s, TLS Version: %d", connState.ServerName, connState.Version)
	connInfo := ServerConnectionInfo{ServerName: connState.ServerName, TLSVersion: connState.Version}

	// textproto gives us some nice, higher level abstractions from the protocol
	txConn := textproto.NewConn(conn)

	code, msg, err := txConn.ReadCodeLine(STATUS_AVAILABLE_POSTING_ALLOWED)
	if err != nil {
		clogger.Fatal(err)
	}

	if code != STATUS_AVAILABLE_POSTING_ALLOWED {
		clogger.Fatal(err)
	}

	clogger.Printf("GOT %d %s", code, msg)

	return &Client{
		Conn:     txConn,
		connInfo: connInfo,
		secure:   true,
		log:      clogger,
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
func (c *Client) List(keyword string) {
	resp, err := c.sendCommand(215, "LIST "+keyword)
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
