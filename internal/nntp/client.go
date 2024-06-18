package nntp

import (
	"bufio"
	"crypto/tls"
	"fmt"
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
		capabilities map[string]bool
		// Keeps track of capabilities ignored by the client
		ignoredCapabilities map[string]bool
	}

	ServerConnectionInfo struct {
		ServerName string
		TLSVersion uint16
	}

	TLS struct {
		*tls.Conn
	}
)

var (
	supportedCapabilities map[string]bool = map[string]bool{
		"AUTHINFO USER": true,
		"LIST":          true,
		"READER":        true,
	}
)

func NewClient(host string, insecure bool, clogger *log.Logger) *Client {
	var (
		addr     string
		c        *textproto.Conn
		err      error
		connInfo *ServerConnectionInfo
	)

	addr = net.JoinHostPort(host, strconv.Itoa(DEFAULT_TLS_PORT))
	clogger.Printf("connecting to '%s'...", addr)

	if insecure {
		addr = net.JoinHostPort(host, strconv.Itoa(DEFAULT_TCP_PORT))
		clogger.Printf(`
================ WARNING ================
Client has been flagged for an insecure connection!
Now connecting to '%s'
=========================================`, addr)
		c, err = textproto.Dial("tcp", addr)
		if err != nil {
			clogger.Fatal(err)
		}
		connInfo = &ServerConnectionInfo{ServerName: "", TLSVersion: 0}

	} else {

		tlsConf := &tls.Config{InsecureSkipVerify: false}
		tlsConn, err := tls.Dial("tcp", addr, tlsConf)
		if err != nil {
			clogger.Fatal(err)
		}

		connState := tlsConn.ConnectionState()
		clogger.Printf("Server Name: %s, TLS Version: %d", connState.ServerName, connState.Version)
		connInfo = &ServerConnectionInfo{ServerName: connState.ServerName, TLSVersion: connState.Version}

		// this lets us use the underlying TLS connection
		// while taking advantage of textproto abstractions
		c = textproto.NewConn(tlsConn)

	}

	code, msg, err := c.ReadCodeLine(STATUS_AVAILABLE_POSTING_ALLOWED)
	if err != nil {
		clogger.Fatal(err)
	}

	if code != STATUS_AVAILABLE_POSTING_ALLOWED {
		clogger.Fatal(err)
	}

	clogger.Printf("GOT %d %s", code, msg)

	return &Client{
		Conn:     c,
		connInfo: *connInfo,
		secure:   !insecure,
		log:      clogger,
	}
}

// Sends a CAPABILITIES command to the server
//
// Caches a list of capabilities both supported and not supported by the client
func (c *Client) SendCapabilities(keyword string) {
	// keyword is reserved for extension modules
	resp := c.sendCommand(STATUS_CAPABILITIES, fmt.Sprint("CAPABILITIES ", keyword))

	switch msg := resp.Message.(type) {
	case Msg:
		c.log.Printf("Status: %d, Message: %s", resp.Status, resp.Message)
	case MultiMsg:
		supportedCapas := make(map[string]bool)
		ignoredCapas := make(map[string]bool)

		for _, capa := range msg {
			if supportedCapabilities[capa] {
				supportedCapas[capa] = true
			} else {
				ignoredCapas[capa] = true
			}
		}

		c.log.Printf("Supported Capabilities: %v", supportedCapas)
		c.log.Printf("Ignored Capabilities: %v", ignoredCapas)

		c.capabilities = supportedCapas
		c.ignoredCapabilities = ignoredCapas
	}
}

// Sends a LIST command to the server
//
// Syntax
//
//	LIST [keyword [wildmat|argument]]
func (c *Client) List(keyword string) {
	resp := c.sendCommand(215, "LIST "+keyword)
	c.log.Printf("Status: %d, Message: %s", resp.Status, resp.Message)
}

func (c *Client) Quit() {
	c.log.Print(c.sendCommand(205, "QUIT"))
}

// Send sends command raw, that is, unchecked and returns the full response.
//
// This will exit if the default bufio Scanner encounters any non-EOF error.
func (c *Client) Send(command string, expectedStatusCode int) string {
	if err := c.PrintfLine(command); err != nil {
		c.log.Fatal("error sending cmd: ", err)
	}

	var response string
	s := bufio.NewScanner(c.R)
	for s.Scan() {
		response += s.Text()
	}

	if err := s.Err(); err != nil {
		c.log.Fatal(err)
	}

	return response
}

func (c *Client) sendCommand(expectedCode int, format string, args ...any) NNTPResponse {
	c.log.Printf("sending command '%s' with args '%s'", format, args)

	if err := c.PrintfLine(format, args...); err != nil {
		c.log.Printf("error sending cmd: %s", err)
		return NNTPResponse{Status: -1, Message: Msg(err.Error())}
	}

	resp, err := c.parseResponse(expectedCode)
	if err != nil {
		return NNTPResponse{Status: -1, Message: Msg(err.Error())}
	}

	return resp
}
