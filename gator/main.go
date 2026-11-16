package main

import (
	"log"
	"os"

	"github.com/BirchNorman/gator/internal/config"
)

type state struct {
	ptrCfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := &state{
		ptrCfg: &cfg,
	}

	cmds := &commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

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
