package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/tepidmilk/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	s := state{
		cfg: &cfg,
	}
	c := commands{
		cmd: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = c.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
