package main

import (
	"encoding/json"
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
	"log"
	"os"
)

func main() {
	guild, err := repository.Guild("Elite BrotherHood")

	if err != nil {
		fmt.Printf("Received unexpected error:\n%+v\n", err)

		os.Exit(1)
	}

	guildJSON, err := json.MarshalIndent(guild, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Guild %s\n", string(guildJSON))
}
