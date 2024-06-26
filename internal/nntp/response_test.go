package nntp_test

import (
	"bufio"
	"fmt"
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

var nntpClient nntp.Client
var clientConfigured bool

func set(c nntp.Client) bool {
	nntpClient = c
	return true
}

func listen(t *testing.T, s net.Conn, done chan bool) error {
	go func() {
		defer close(done)

		for {
			select {
			case <-time.After(15 * time.Second):
				done <- true
			}
		}
	}()

	// add a 5 second write deadline
	s.SetWriteDeadline(time.Now().Add(2 * time.Second))

	// server writes the initial response
	n, err := s.Write([]byte(responses[0]))
	require.NoError(t, err)

	log.Print("wrote bytes: ", n)
	// t.Parallel()

	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		var response string
		switch scanner.Text() {
		case "CAPABILITIES":
			response = responses[1] // hard-coded response
		case "QUIT":
			response = responses[2]
		default:
			continue
		}

		if response != "" {
			_, err := s.Write([]byte(response))
			require.NoError(t, err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner err: %v", err)
		done <- true
	}

	<-done
	return nil
}

func setup(t *testing.T) error {
	c, s := net.Pipe()

	cl := nntp.Client{Conn: textproto.NewConn(s)}
	if ok := set(cl); ok {
		clientConfigured = true
	}

	s.Write([]byte(responses[0]))

	// test clean up hooks
	t.Cleanup(func() {
		fmt.Print("cleaning up tests...")

		// clean up the underlying in-memory pipes
		if err := c.Close(); err != nil {
			log.Print(err)
		}

		if err := s.Close(); err != nil {
			log.Print(err)
		}
	})

	ready := make(chan bool, 1)
	if err := listen(t, s, ready); err != nil {
		return err
	}
	return nil
}

func TestEnsureInitialResponse(t *testing.T) {
	setup(t)

	if clientConfigured {
		log.Print("client is ready")
	}

	// dat, err := io.ReadAll(nntpClient.R)
	// require.NoError(t, err)
	// require.Equal(t, responses[0], string(dat))

	// if run := t.Run("setup", setup); run {
	// 	<-ready
	//
	// 	dat, err := io.ReadAll(nntpClient.R)
	// 	require.NoError(t, err)
	// 	require.Equal(t, responses[0], string(dat))
	// }
}
