package main

import (
	"encoding/json"
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/handler"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
	"os"
)

func main() {
	params := handler.CreateReportCommandParameters{GuildName: "Elite Brotherhood"}
	command := handler.CreateReportCommand{Repository: repository.TibiaDataAPI{}, Params: params}

	report, err := handler.Execute(command)

	if err != nil {
		fmt.Printf("Received unexpected error:\n%+v\n", err)

		os.Exit(1)
	}

	reportJSON, err := json.MarshalIndent(report, "", "  ")

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Report %s\n", string(reportJSON))
}
