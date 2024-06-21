package nntp

import (
	"bufio"
	"strconv"
	"strings"
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
	initial, err := c.ReadLine()
	if err != nil {
		return NNTPResponse{Status: -1, Message: Msg(initial)}, err
	}

	r := strings.NewReader(initial)
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	var (
		code     int
		messages []string
		multi    bool
	)

	for sc.Scan() {

		token := sc.Text()
		// check multi-line dot EOF
		if token[0] == byte('.') {
			clogger.Debug("reached multi-line EOF", "code", code, "msg len", len(messages), "messages", messages)

			multi = true
			break
		}

		if len(token) == 3 {
			parsedCode, _ := strconv.Atoi(token)
			if parsedCode != expectedCode {
				clogger.Error("parsed status code mismatch", "parsed", parsedCode, "expected", expectedCode)
			}
			code = parsedCode
			continue // dont append the status code to the messages slice
		}

		messages = append(messages, strings.TrimSpace(token))

		if err := sc.Err(); err != nil {
			return NNTPResponse{Status: -1, Message: nil}, err
		}
	}

	joined := strings.Join(messages, " ")
	clogger.Info("GOT RESPONSE", "code", code, "joined", joined, "multi", multi)

	if multi {
		return NNTPResponse{Status: code, Message: MultiMsg{joined}}, nil
	}

	return NNTPResponse{Status: code, Message: Msg(joined)}, nil
}
