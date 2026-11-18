package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BirchNorman/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Fatal("a username is required")
	}

	dbUser, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.ptrCfg.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Fatalf("a username is required")
	}
	db := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	dbUser, err := s.db.CreateUser(context.Background(), db)
	if err != nil {
		return err
	}
	err = s.ptrCfg.SetUser(dbUser.Name)
	fmt.Println(s.ptrCfg.CurrentUserName)

	return nil
}
