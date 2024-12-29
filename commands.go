package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, fn func(*state, command) error) {
	c.cmds[name] = fn
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	return f(s, cmd)
}
