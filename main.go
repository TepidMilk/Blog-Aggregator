package main

import (
	"fmt"

	"github.com/tepidmilk/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error while reading config file: ", err)
	}
	cfg.SetUser("Noah")
	fmt.Println(config.Read())
}
