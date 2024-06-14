package nntp

// https://datatracker.ietf.org/doc/html/rfc3977

// The official TCP port for the NNTP service is 119.  However, if a
// host wishes to offer separate servers for transit and reading
// clients, port 433 SHOULD be used for the transit server and 119 for
// the reading server.

const (
	DEFAULT_TCP_PORT     int = 119
	DEFAULT_TLS_PORT     int = 563
	DEFAULT_TRANSIT_PORT int = 433
)
