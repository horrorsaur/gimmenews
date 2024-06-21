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
var ready = make(chan bool)

func setup(t *testing.T, ready *bool) {
	c, s := net.Pipe()
	nntpClient = &nntp.Client{Conn: textproto.NewConn(s)}

	// // a blocking channel, ensuring the server reads until timeout
	done := make(chan error)

	// add a 1 second write deadline
	s.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// server writes the initial response
	n, err := c.Write([]byte(responses[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes", n)

	fmt.Print("ready ", *ready)
	// server is ready
	*ready = true

	fmt.Print("ready ", *ready)

	var response string
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanWords)

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

	t.Cleanup(func() {
		log.Printf("closing connections... got err: %s", <-done)
		// clean up the underlying in-memory pipes
		if err := c.Close(); err != nil {
			log.Print(err)
		}

		if err := s.Close(); err != nil {
			log.Print(err)
		}
	})
}

func TestEnsureInitialResponse(t *testing.T) {
	var (
		r     bool
		retry int = 0
	)

	go func() {
		setup(t, &r)
	}()

	fmt.Printf("retry counter %d ready? %v", retry, r)
	//
	// for {
	// 	select {
	// 	case <-ready:
	// 		dat, err := io.ReadAll(nntpClient.R)
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 		require.Equal(t, string(dat), responses[0])
	// 	default:
	// 		if retry != 3 {
	// 			retry++
	// 			continue
	// 		}
	//
	// 		fmt.Printf("reached max retries!")
	// 		t.Fail()
	// 	}
	// }
}
