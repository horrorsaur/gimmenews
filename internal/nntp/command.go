package nntp

type (
	// The character set for all NNTP commands is UTF-8
	//
	// Keywords & arguments must be separated by one or more space or TAB characters
	//
	// A CRLF pair MUST terminate all commands
	Command struct {
		// Commands may have variants (e.g., LIST or MODE)
		//
		// The second keyword is immediately after the first
		Keyword   string
		Arguments string

		// '*' and '?' are always wildcards.
	}
)
