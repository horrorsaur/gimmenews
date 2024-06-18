package nntp

// For servers using the NNTP-AUTH extension (https://datatracker.ietf.org/doc/html/rfc4643)
//
// Sends the authentication commands to the remote server.
//
// If the

type (
	AuthResponse struct {
	}

	AuthenticationPhase struct {
		Code int
	}
)

// If the server requires only a username, it MUST NOT give a 381 response to AUTHINFO USER
// and MUST give a 482 response to AUTHINFO PASS.
// if authUserResp, err := c.sendCommand(42, "AUTHINFO USER %s", user); err != nil {}
// Sends the authentication commands, AUTHINFO USER and AUTHINFO PASS
//
//	Parameters
//
// username
// password
func (c *Client) SendAuthInfo(user, pass string) {
	c.sendCommand(381, "AUTHINFO USER %s", user)
	c.sendCommand(281, "AUTHINFO PASS %s", pass)

	// refresh capabilities
	// c.SendCapabilities("")
}

func (c *Client) SupportsAuth() bool {
	return c.capabilities["AUTHINFO USER"]
}
