package commands

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
)

type (
	Command struct {
		Name string
		Path string

		out io.Writer
		err io.Writer
	}

	CommandOpts struct {
		// The name of the NNCP command
		Name string
		// The path to the command's binary
		Path string
	}
)

var (
	logger log.Logger = *log.New(os.Stdout, "[COMMANDS ]", 1)

	newCfg Command = Command{
		Name: "nncp-cfgnew",
		Path: "/usr/bin/",
	}

	stat Command = Command{
		Name: "nncp-stat",
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

	// commands map[string]Command = map[string]Command{
	// 	"stat": stat,
	// }
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
	cmd.out = os.Stdout

	return cmd, err
}

func (c *Command) SetOutput(w io.Writer) {
	c.out = w
}

func Stat() ([]byte, error) {
	cmd := stat.load()

	dat, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return dat, nil
}

// Load builds an exec.Cmd passing 'args' for additional nncp options
func (c *Command) load(args ...string) *exec.Cmd {
	cmd := exec.Command(c.fullPath(), args...)

	if c.out != nil {
		logger.Print("setting stdout")
		cmd.Stdout = c.out
	}

	if c.err != nil {
		logger.Print("setting stderr")
		cmd.Stderr = c.err
	}

	logger.Print("using default in-memory reader")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (c *Command) fullPath() string {
	return c.Path + c.Name
}
