package main

import (
	"errors"
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app-eb-guild-ops/handler"
	"github.com/adriein/eb-guild-ops/internal/app-eb-guild-ops/markdown"
	"github.com/adriein/eb-guild-ops/internal/app-eb-guild-ops/repository"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	discord := getDiscordApi()

	guildID, existEnv := os.LookupEnv("ELITE_BROTHERHOOD_GUILD_ID")

	if !existEnv {
		noEnvVarError := errors.New("ELITE_BROTHERHOOD_GUILD_ID is not set")
		fmt.Printf("Received unexpected error:\n %+v\n", noEnvVarError)

		os.Exit(1)
	}

	channel := fetchGuildChannel(guildID, discord)

	report := generateReport()

	markdownReport := generateMarkdownReport(report)

	sendReport(channel, markdownReport, discord)
}

func generateReport() handler.EbGuildReport {
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

	return report
}

func fetchGuildChannel(guildID string, discord *repository.DiscordApi) repository.DiscordChannel {
	channel, fetchChannelError := discord.FetchChannel(guildID, "inactive-members-bot")

	if fetchChannelError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", fetchChannelError)

		os.Exit(1)
	}

	return channel
}

func getDiscordApi() *repository.DiscordApi {
	dotenvError := godotenv.Load()

	if dotenvError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", dotenvError)
		os.Exit(1)
	}

	discord, instantiateDiscordError := repository.NewDiscordRepository()

	if instantiateDiscordError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", instantiateDiscordError)

		os.Exit(1)
	}

	return discord
}

func generateMarkdownReport(report handler.EbGuildReport) string {
	converter, instantiateConverterError := markdown.NewMarkdownConverter()

	if instantiateConverterError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", instantiateConverterError)

		os.Exit(1)
	}

	markdownReport, converterError := converter.ConvertReport(report)

	if converterError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", converterError)

		os.Exit(1)
	}

	return markdownReport
}

func sendReport(channel repository.DiscordChannel, report string, discord *repository.DiscordApi) {
	if sendMessageError := discord.Message(channel.Id, report); sendMessageError != nil {
		fmt.Printf("Received unexpected error:\n %+v\n", sendMessageError)

		os.Exit(1)
	}
}
