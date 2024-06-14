package nntp

type (
	Response struct {
		// A response MUST begin with a 3-digit status indicator
		Status  StatusCode
		Message string
		Err     error
	}

	MultilineResponse struct {
		Status  StatusCode
		Message []string
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
	// An array to check for response codes that are noted as multiline
	multilineResponseCodes = [10]int{100, 101, 231}
)

const (
	STATUS_CAPABILITIES int = 101
	STATUS_SERVER_DATE  int = 111

	STATUS_AVAILABLE_POSTING_ALLOWED     int = 200 // Service is OK and indicates we are allowed to POST articles
	STATUS_AVAILABLE_POSTING_NOT_ALLOWED int = 201 // Service is OK but indicates that we are not allowed to POST articles
	STATUS_QUIT                          int = 205 // Indicates the server closed the connection immediately
	STATUS_AUTHENTICATION_REQUIRED       int = 480 // Indicates the server is using an authentication extension
)
