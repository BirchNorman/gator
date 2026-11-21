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
		log.Fatal("a username is required")
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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		log.Fatal("Unsucessfully deleted users")
	}
	fmt.Print("Sucessfully deleted users")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.ptrCfg.CurrentUserName {
			fmt.Printf("\n* %v (current)", user)
		} else {
			fmt.Printf("\n* %v", user)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("handlerAgg fetchFeed error: %w", err)
	}
	fmt.Print(rssFeed)

	return nil
}
