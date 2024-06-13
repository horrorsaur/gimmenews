package nntp

type (
	Response struct {
		// A response MUST begin with a 3-digit status indicator
		Status StatusCode
	}

	// 1xx - Informative message
	//
	// 2xx - Command completed OK
	//
	// 3xx - Command OK so far; send the rest of it
	//
	// 4xx - Command was syntactically correct but failed for some reason
	//
	// 5xx - Command unknown, unsupported, unavailable, or syntax error
	StatusCode int
)

var (
	multilineResponseCodes = [10]int{101}
)

const (
	STATUS_INITIAL_OK              int = 200
	STATUS_CAPABILITIES            int = 101
	STATUS_SERVER_DATE             int = 111
	STATUS_AUTHENTICATION_REQUIRED int = 480
)
