package commands

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// path, err := exec.LookPath("nncp-cfgnew")

type (
	Command struct {
		*exec.Cmd

		// The actual name of the command invoked (e.g., nncp-toss)
		Name string

		// Embedded through exec.Cmd. This field must be set
		//
		// The path to the command's binary
		Path string
	}

	CommandOpts struct {
		// The actual name of the command invoked (e.g., nncp-toss)
		Name string
		// The path to the command's binary
		Path string
	}
)

var (
	logger log.Logger = *log.New(os.Stdout, "[COMMANDS] ", 1)

	daemon Command = Command{
		Name: "nncp-daemon",
		Path: "/usr/bin/",
	}

	call Command = Command{
		Name: "nncp-call",
		Path: "/usr/bin/",
	}

	caller Command = Command{
		Name: "nncp-caller",
		Path: "/usr/bin/",
	}

	toss Command = Command{
		Name: "nncp-toss",
		Path: "/usr/bin/",
	}
)

// New builds a Command wrapper, taking care of calling Command.load()
func NewCommand(opts CommandOpts) (Command, error) {
	var (
		c   Command
		err error
	)

	if opts.Name == "" || opts.Path == "" {
		err = errors.New("Received blank name or path value")
		return c, err
	}

	c.Name = opts.Name
	c.Path = opts.Path

	c.load()

	return c, err
}

// Load passes callers args to the Command struct
func (c *Command) load(args ...string) {
	c.Cmd = exec.Command(c.fullPath(), args...)
}

func (c *Command) fullPath() string {
	return c.Path + c.Name
}
