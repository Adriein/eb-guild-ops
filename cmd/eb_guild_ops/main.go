package main

import (
	"encoding/json"
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/handler"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
	"os"
)

func main() {
	tibiaApiRepository, instantiateRepositoryError := repository.NewTibiaDataAPI()

	if instantiateRepositoryError != nil {
		fmt.Printf("Received unexpected error:\n%+v\n", instantiateRepositoryError)

		os.Exit(1)
	}

	createReportHandler, instantiateHandlerError := handler.NewCreateReportHandler(tibiaApiRepository)

	if instantiateHandlerError != nil {
		fmt.Printf("Received unexpected error:\n%+v\n", instantiateHandlerError)

		os.Exit(1)
	}
	command := handler.CreateReportCommand{GuildName: "Elite Brotherhood"}

	report, err := createReportHandler.Execute(command)

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
