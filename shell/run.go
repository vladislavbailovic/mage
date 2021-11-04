package shell

import (
	"fmt"
	"os/exec"
)

type Command struct {
	raw      []string
	executed bool
	out      []byte
}

func NewCommand(cmd []string) *Command {
	return &Command{cmd, false, []byte{}}
}

func (c *Command) GetStdout() (string, error) {
	if c.executed {
		return string(c.out[:]), nil
	}
	err := c.Run()
	if err != nil {
		return "", err
	}
	return c.GetStdout()
}

func (c *Command) Run() error {
	c.executed = true
	if len(c.raw) < 1 {
		return fmt.Errorf("not enough arguments for command")
	}
	args := []string{"-c"}
	out, err := exec.Command(GetShellBinary(), append(args, c.raw...)...).Output()
	if err != nil {
		return err
	}
	c.out = out
	return nil
}

func GetShellBinary() string {
	return "/bin/sh"
}
