package nntp

type (
	// The character set for all NNTP commands is UTF-8
	//
	// A CRLF pair MUST terminate all commands
	// Keywords & arguments must be separated by one or more space or TAB characters
	Command struct {
		Keyword   string
		Arguments string
	}
)
