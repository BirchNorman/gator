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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		log.Fatal("a feed name and URL are required")
	}

	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	dbUser, err := s.db.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
	if err != nil {
		return err
	}

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    dbUser.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), arg)
	if err != nil {
		return err
	}

	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		log.Fatal("no arguments are required")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {

		creatorName, err := s.db.GetName(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(creatorName)

	}

	return nil
}
