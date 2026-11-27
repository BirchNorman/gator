package main

import (
	"context"
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
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.ptrCfg.CurrentUserName)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		return err
	}
}
