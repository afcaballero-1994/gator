package main

import (
	"context"
	"fmt"
	"github.com/afcaballero-1994/gator/internal/config"
	"github.com/afcaballero-1994/gator/internal/database"
	"github.com/google/uuid"
	"html"
	"os"
	"time"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	ars  []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.ars) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.name)
	}
	name := cmd.ars[0]

	ctx := context.Background()

	r, err := s.db.GetUser(ctx, name)
	if err != nil {
		return err
	}

	s.cfg.SetUser(r.Name)

	fmt.Printf("Change: %s has been set as user\n", r.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	ctx := context.Background()
	uname := cmd.ars[0]

	ruser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      uname,
	})

	if err != nil {
		return err
	}

	s.cfg.SetUser(uname)
	fmt.Println("User with name:", ruser.Name, "was created")
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.ResetTable(ctx)
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	rusers, err := s.db.GetUsers(ctx)

	if err != nil {
		return err
	}
	for _, u := range rusers {
		pu := "* " + u.Name
		if u.Name == s.cfg.Current_username {
			pu += " (current)"
		}
		fmt.Println(pu)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()

	feed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println("Channel Title:", html.UnescapeString(feed.Channel.Title))
	fmt.Println("Channel Description:", html.UnescapeString(feed.Channel.Description))

	for _, post := range feed.Channel.Item {
		fmt.Println("*", html.UnescapeString(post.Title))
		fmt.Println("-", html.UnescapeString(post.Description))
	}

	return nil
}

func handlerAddFeed(user database.User, s *state, cmd command) error {
	if len(cmd.ars) != 2 {
		os.Exit(1)
	}
	ctx := context.Background()
	rf, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.ars[0],
		Url:       cmd.ars[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	ff, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    rf.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println(ff)

	fmt.Println("Feed created successfully:")
	printFeed(rf)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	ctx := context.Background()

	fds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		fmt.Println("Name:", fd.Name)
		fmt.Println("Url:", fd.Url)
		fmt.Println("Username:", fd.Username)
		fmt.Println("------//---------------")
	}

	return nil
}

func handlerFollow(user database.User, s *state, cmd command) error {
	if len(cmd.ars) != 1 {
		fmt.Println("Need one argument")
		os.Exit(1)
	}

	ctx := context.Background()

	f, err := s.db.GetFeed(ctx, cmd.ars[0])
	if err != nil {
		return fmt.Errorf("Error getting feed: %s, err: %w", cmd.ars[0], err)
	}

	ff, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    f,
	})
	if err != nil {
		return fmt.Errorf("Error adding follow: %w", err)
	}

	fmt.Println("User name:", ff.UserName)
	fmt.Println("Feed name:", ff.FeedName)

	return nil
}
func handlerFollowing(user database.User, s *state, cmd command) error {
	ctx := context.Background()

	ffs, err := s.db.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return err
	}
	if len(ffs) > 0 {
		fmt.Println("User:", ffs[0].UserName)
		for _, feed := range ffs {
			fmt.Println("Feed name:", feed.FeedName)
		}
	}

	return nil
}

func handlerUnfollow(user database.User, s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.DeleteFollowForUser(ctx, database.DeleteFollowForUserParams{
		Name: user.Name,
		Url:  cmd.ars[0],
	})
	if err != nil {
		return fmt.Errorf("Error unfollowing. Err: %w", err)
	}

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handle, exist := c.cmds[cmd.name]
	if !exist {
		return fmt.Errorf("%s command does not exist", cmd.name)
	}
	return handle(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
