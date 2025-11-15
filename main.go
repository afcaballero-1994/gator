package main

import _ "github.com/lib/pq"

import (
	"os"
	"fmt"
	"database/sql"
	"github.com/afcaballero-1994/gator/internal/config"
	"github.com/afcaballero-1994/gator/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Need at least one argument")
		os.Exit(1)
	}
	var s state
	cfg, err := config.Read()
	if err != nil{
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
		cmds : make(map[string]func(*state, command) error),
	}
	cms.register("login", handlerLogin)
	cms.register("register", handlerRegister)
	cms.register("reset", handlerReset)
	cms.register("users", handlerUsers)

	c := command{
		name: os.Args[1],
		ars: os.Args[2:],
	}
	err = cms.run(&s, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}
