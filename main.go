package main

import (
	"fmt"
	"github.com/afcaballero-1994/gator/internal/config"
)

func main() {
	con, err := config.Read()
	if err != nil{
		fmt.Println("Error:", err)
	}
	con.SetUser("andres")

	con, err = config.Read()
	fmt.Println(con)
}
