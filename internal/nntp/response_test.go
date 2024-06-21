package nntp_test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"testing"
	"time"

	"github.com/horrorsaur/gimmenews/internal/nntp"
	"github.com/stretchr/testify/require"
)

var responses = []string{
	`
200 news.eternal-september.org InterNetNews NNRP server INN 2.8.0 (20240323 snapshot) ready (posting ok)
	`,
	`
101 Capabilities List:
VERSION 2
COMPRESS DEFLATE
IMPLEMENTATION INN 2.8.0 (20240323 snapshot)
	`,
	`
205    Connection closing
	`,
}

// https://stackoverflow.com/questions/30688685/how-does-one-test-net-conn-in-unit-tests-in-golang

var nntpClient *nntp.Client

func setup() {
	c, s := net.Pipe()
	nntpClient = &nntp.Client{Conn: textproto.NewConn(s)}

	// // a blocking channel, ensuring the server reads until timeout
	done := make(chan error)

	var response string
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanWords)

	// add a 1 second write deadline
	s.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// server writes the initial response
	n, _ := c.Write([]byte(responses[0]))
	fmt.Printf("wrote %d bytes", n)
	s.Close()

	go func() {
		for scanner.Scan() {
			switch scanner.Text() {
			case "CAPABILITIES":
				response = responses[1] // hard-coded response
			case "QUIT":
				response = responses[2]
			default:
				response = ""
				continue
			}

			if response != "" {
				c.Write([]byte(response))
			}
		}

		if err := scanner.Err(); err != nil {
			done <- err
		}
	}()

	scanErr := <-done
	log.Printf("closing connections... got err: %s", scanErr)
	// clean up the underlying in-memory pipes
	if err := c.Close(); err != nil {
		log.Print(err)
	}

	if err := s.Close(); err != nil {
		log.Print(err)
	}
}

func TestEnsureInitialResponse(t *testing.T) {
	setup()
	dat, err := io.ReadAll(nntpClient.R)
	if err != nil {
		log.Print(err)
	}
	require.Equal(t, string(dat), responses[0])
}

// func TestInitialResponse(t *testing.T) {
// 	setup()
//
// 	r := bufio.NewReader(strings.NewReader(initialResponse))
//
// 	c.SendCapabilities("")
//
// 	initialResponseParse := NNTPResponse{Status: 200, Message: nntp.Msg("news.eternal-september.org InterNetNews NNRP server INN 2.8.0 (20240323 snapshot) ready (posting ok)")}
//
// 	require.Equal(t, NNTPResponse{}, initialResponseParse)
// }

// func TestPipeStuff(t *testing.T) {
// 	c, s := net.Pipe()
//
// 	client := nntp.Client{Conn: textproto.NewConn(s)}
//
// 	go func() {
// 		s.SetWriteDeadline(time.Now().Add(5 * time.Second))
//
// 		n, err := c.Write([]byte(responses[0]))
// 		if err != nil {
// 			log.Print(err)
// 		}
// 		log.Printf("wrote %d", n)
//
// 		s.Close()
// 	}()
//
// 	time.Sleep(5 * time.Second)
//
// 	dat, err := io.ReadAll(client.R)
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	log.Printf(string(dat))
// }
