package main

import (
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("a username is required")
	}
	err := s.ptrCfg.SetUser(cmd.args[0])
	if err != nil {
		log.Fatalf("a username was not set")
	}
	return nil
}
