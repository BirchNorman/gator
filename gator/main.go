package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/BirchNorman/gator/internal/config"
	"github.com/BirchNorman/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	ptrCfg *config.Config
	db     *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)

	s := &state{
		ptrCfg: &cfg,
		db:     database.New(db),
	}

	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

	userInput := os.Args
	if len(userInput) < 2 {
		log.Fatalf("not enough arguments were provided")
	}

	cmd := command{
		name: userInput[1],
		args: userInput[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
