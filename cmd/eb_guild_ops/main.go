package main

import (
	"encoding/json"
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
)

func main() {
	char, err := repository.Character("Marlock")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", char)
	a, err := json.Marshal(char)
	fmt.Println(string(a))
}
