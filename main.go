package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/tepidmilk/gator/internal/config"
	"github.com/tepidmilk/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(db)
	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	c := commands{
		cmd: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAgg)
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
