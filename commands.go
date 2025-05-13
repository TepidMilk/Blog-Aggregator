package main

import (
	"errors"
	"fmt"

	"github.com/tepidmilk/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmd[cmd.name]
	if !ok {
		return errors.New("invalid command name")
	}
	err := val(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) <= 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	err := s.cfg.SetUser(cmd.args[1])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %W", err)
	}

	fmt.Printf("User has been set to %s\n", cmd.args[1])
	return nil
}
