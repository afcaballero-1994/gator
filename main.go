package main

import (
	"os"
	"fmt"
	"github.com/afcaballero-1994/gator/internal/config"
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

	cms := commands{
		cmds : make(map[string]func(*state, command) error),
	}
	cms.register("login", handlerLogin)

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
