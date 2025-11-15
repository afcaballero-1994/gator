package main

import(
	"fmt"
	"context"
	"github.com/google/uuid"
	"time"
	"github.com/afcaballero-1994/gator/internal/config"
	"github.com/afcaballero-1994/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	ars []string
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
		return err;
	}
	
	s.cfg.SetUser(r.Name)

	fmt.Printf("Change: %s has been set as user\n", r.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	ctx := context.Background()
	uname := cmd.ars[0]

	ruser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: uname,
	})

	if err != nil{
		return err
	}

	s.cfg.SetUser(uname)
	fmt.Println("User with name:", ruser.Name, "was created")
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	err := s.db.ResetTable(ctx)
	if err != nil{
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
	for _,u := range rusers {
		pu := "* " + u.Name
		if u.Name == s.cfg.Current_username {
			pu += " (current)"
		}
		fmt.Println(pu)
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
