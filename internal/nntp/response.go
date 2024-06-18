package nntp

import (
	"bufio"
)

type (
	NNTPResponse struct {
		// Responses MUST begin with a 3-digit status indicator
		//
		// 1xx - Informative message
		//
		// 2xx - Command completed OK
		//
		// 3xx - Command OK so far; send the rest of it
		//
		// 4xx - Command was syntactically correct but failed for some reason
		//
		// 5xx - Command unknown, unsupported, unavailable, or syntax error
		Status  int
		Message ResponseMessage
	}

	Msg      string
	MultiMsg []string
)

type ResponseMessage interface {
	Message()
}

func (m Msg) Message()      {}
func (m MultiMsg) Message() {}

var (
	// An array to check for response codes that are noted as multiline
	// remove this in favor of better response parsing
	multilineResponseCodes = [10]int{100, 101, 231}

	testResponse = []string{
		"101 Capability list:",
		"VERSION 2",
		"",
	}
)

const (
	STATUS_CAPABILITIES int = 101
	STATUS_SERVER_DATE  int = 111

	STATUS_AVAILABLE_POSTING_ALLOWED     int = 200 // Service is OK and indicates we are allowed to POST articles
	STATUS_AVAILABLE_POSTING_NOT_ALLOWED int = 201 // Service is OK but indicates that we are not allowed to POST articles
	STATUS_QUIT                          int = 205 // Indicates the server closed the connection immediately
	STATUS_AUTHENTICATION_REQUIRED       int = 480 // Indicates the server is using an authentication extension

	// Authentication via NNTP-AUTH
	AUTHENTICATION_ACCEPTED                 int = 281 // Successfully authenticated
	STATUS_PASSWORD_REQUIRED                int = 381 // Username accepted, pending password
	AUTHENTICATION_FAILED                   int = 481 // No go, you failed
	AUTHENTICATION_COMMANDS_OUT_OF_SEQUENCE int = 482

	UNAVAILABLE int = 502
)

// If the response is multi-line, ReadCodeLine returns an error.
//
// An expectCode <= 0 disables the check of the status code.
//
// response = simple-response / multi-line-response
// simple-response = initial-response-line
// multi-line-response = initial-response-line multi-line-data-block
func (c *Client) parseResponse(expectedCode int) (NNTPResponse, error) {
	code, msg, err := c.ReadCodeLine(expectedCode)
	if err != nil {
		c.log.Printf("received error on initial response: %s", err)
	}

	c.log.Printf("GOT INITIAL '%d %s'", code, msg)

	var multi bool
	// known multiline codes
	for _, c := range multilineResponseCodes {
		if code == c {
			multi = true
		}
	}

	// switch code {
	// case STATUS_PASSWORD_REQUIRED:
	// 	return NNTPResponse{Status: code, Message: msg}, nil
	// default:
	// 	multi = true
	// }

	var scanResult []string
	if multi {

		s := bufio.NewScanner(c.R)
		for s.Scan() {
			if s.Bytes()[0] == byte('.') {
				c.log.Printf("reached multi-line response EOF")
				break
			}

			scanResult = append(scanResult, s.Text())
		}

		if err := s.Err(); err != nil {
			c.log.Printf("got scan error %s", err)
		}

	}

	c.log.Printf("scan result len: %d", len(scanResult))

	if len(scanResult) == 0 {
		return NNTPResponse{Status: code, Message: Msg(msg)}, nil
	}

	// result := strings.Join(scanResult, " ")
	return NNTPResponse{
		Status:  code,
		Message: MultiMsg(scanResult),
	}, nil
}
