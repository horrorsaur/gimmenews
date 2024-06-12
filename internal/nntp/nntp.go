package nntp

import (
	"log"
	"net"
	"net/textproto"
)

// https://datatracker.ietf.org/doc/html/rfc3977

func Connect() {
	addr := net.JoinHostPort("news.eternal-september.org", "119")
	log.Printf("attempting connect to '%s'", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	txtConn := textproto.NewConn(conn)
	defer txtConn.Close()
}
