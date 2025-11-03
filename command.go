package main

import(
	"fmt"
	"github.com/afcaballero-1994/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	ars []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.ars) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}
	s.cfg.SetUser(cmd.ars[0])

	fmt.Printf("Change: %s has been set as user\n", cmd.ars[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handle, exist := c.cmds[cmd.name]
	if !exist {
		return fmt.Errorf("%s command does not exist", cmd.name)
	}
	return handle(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
