package commands

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

type (
	Command struct {
		Name string
		Path string

		Output io.Writer
	}

	CommandOpts struct {
		// The name of the NNCP command
		Name string
		// The path to the command's binary
		Path string
		// Where the command will write to
		Output io.Writer
	}
)

var (
	newCfg Command = Command{
		Name: "nncp-cfgnew",
		Path: "/usr/bin/",
	}

	stat Command = Command{
		Name: "nncp-cfgstat",
		Path: "/usr/bin/",
	}

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

func NewCommand(opts CommandOpts) (Command, error) {
	var (
		cmd Command
		err error
	)

	if opts.Name == "" || opts.Path == "" {
		err = errors.New("Received blank name or path value")
		return cmd, err
	}

	cmd.Name = opts.Name
	cmd.Path = opts.Path
	cmd.Output = os.Stdout

	return cmd, err
}

func CfgNew() {}

func Stat() error {
	cmd := stat.load()
	if cmd.Err != nil {
		return cmd.Err
	}

	return nil
}

func (c *Command) execute() {
	// TODO: Add in arguments
}

// Load builds an 'os/exec' Command
func (c *Command) load() *exec.Cmd {
	return exec.Command(c.binaryPath())
}

func (c *Command) binaryPath() string {
	return c.Path + c.Name
}
