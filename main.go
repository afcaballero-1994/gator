package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/afcaballero-1994/gator/internal/config"
	"github.com/afcaballero-1994/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Need at least one argument")
		os.Exit(1)
	}
	var s state
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error:", err)
	}
	s.cfg = &cfg
	db, err := sql.Open("postgres", s.cfg.DB_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	cms := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	cms.register("login", handlerLogin)
	cms.register("register", handlerRegister)
	cms.register("reset", handlerReset)
	cms.register("users", handlerUsers)
	cms.register("agg", handlerAgg)
	cms.register("addfeed", userLogin(handlerAddFeed))
	cms.register("feeds", handlerGetFeeds)
	cms.register("follow", userLogin(handlerFollow))
	cms.register("following", userLogin(handlerFollowing))
	cms.register("unfollow", userLogin(handlerUnfollow))
	cms.register("browse", userLogin(handlerBrowse))

	c := command{
		name: os.Args[1],
		ars:  os.Args[2:],
	}
	err = cms.run(&s, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func userLogin(handler func(user database.User, s *state, cmd command) error) func(*state, command) error {
	ctx := context.Background()
	return func(s *state, cmd command) error {
		u, err := s.db.GetUser(ctx, s.cfg.Current_username)
		if err != nil {
			return fmt.Errorf("Error login to user: %w", err)
		}
		return handler(u, s, cmd)
	}

}
