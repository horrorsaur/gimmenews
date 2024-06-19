package nntp_test

import (
	"bufio"
	"log"
	"net"
	"net/textproto"
	"os"
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

var testServer TestServer
var testClient nntp.Client

type TestServer struct {
	net.Conn
}

func setup() {
	c, s := net.Pipe()
	log.Printf("creating fake TCP testing pipe")
	testClient = nntp.Client{Conn: textproto.NewConn(s)}

	testServer = TestServer{Conn: c}
	testServer.SetReadDeadline(time.Now().Add(5 * time.Second))

	// server writes the initial response
	n, err := testClient.Conn.W.Write([]byte(responses[0]))
	if err != nil {
		panic(err)
	}
	log.Printf("wrote %d bytes", n)

	// // a blocking channel, ensuring the server reads until timeout
	done := make(chan error)

	var response string
	scanner := bufio.NewScanner(testServer.Conn)
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

		testServer.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		done <- err
	}

	<-done

	log.Printf("closing connections...")

	testServer.Close()
	testClient.Close()
}

// func TestInitialSetup(t *testing.T) {
// 	setup()
// 	// go func() {
// 	// 	setup()
// 	// }()
//
// 	// time.Sleep(5 * time.Second)
//
// 	resp, err := testClient.ReadLine()
// 	if err != nil {
// 		t.Fatalf(err.Error(), t)
// 	}
//
// 	log.Printf(resp)
// }

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

func TestPipeStuff(t *testing.T) {
	c, s := net.Pipe()

	client := nntp.Client{Conn: s}
	tlogger := log.New(os.Stdout, "[TEST] ", 1)
	client := nntp.NewClient("", true, tlogger)
	log.Printf("client: %v", client)
	// testClient = nntp.Client{Conn: textproto.NewConn(s)}
}
