package main

import (
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
	"os"
)

func main() {
	guild, err := repository.Guild("Elite BrotherHood")

	if err != nil {
		fmt.Printf("Received unexpected error:\n%+v\n", err)

		os.Exit(1)
	}

	fmt.Printf("%#v\n", guild)
}
