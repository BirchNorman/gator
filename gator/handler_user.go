package main

import (
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatalf("a username is required")
	}
	err := s.ptrConfig.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	return nil
}
