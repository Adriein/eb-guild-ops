package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DiscordRepository struct {
	token   string
	baseUrl string
	version string
}

type DiscordGuildChannelsResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

const DISCORD_BOT_TOKEN = "DISCORD_BOT_TOKEN"
const DISCORD_BASE_URL = "DISCORD_BOT_TOKEN"
const DISCORD_API_VERSION = "DISCORD_BOT_TOKEN"

func New() (*DiscordRepository, error) {
	token, hasEnvVar := os.LookupEnv(DISCORD_BOT_TOKEN)

	if !hasEnvVar {
		return nil, errors.New("DISCORD_BOT_TOKEN is not set")
	}

	authorizationToken := fmt.Sprintf("Bot %s", token)

	return &DiscordRepository{token: authorizationToken}, nil
}

func (discord *DiscordRepository) sendMessage(channelID string) error {
	return nil
}

func (discord *DiscordRepository) fetchChannel(guildID string, name string) (DiscordGuildChannelsResponse, error) {
	request, requestCreationError := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s/%s",
			TIBIA_DATA_API_BASE_URL,
			TIBIA_DATA_API_VERSION,
			TIBIA_DATA_API_CHARACTER_URL,
			name,
		),
		nil,
	)

	if requestCreationError != nil {
		return TibiaDataAPICharacter{}, requestCreationError
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return TibiaDataAPICharacter{}, requestError
	}

	defer response.Body.Close()

	var apiResponse = TibiaDataAPICharacterJSON{}

	byteRawResponse, readBodyError := io.ReadAll(response.Body)

	if readBodyError != nil {
		return TibiaDataAPICharacter{}, readBodyError
	}

	jsonUnmarshalError := json.Unmarshal(byteRawResponse, &apiResponse)

	if jsonUnmarshalError != nil {
		return TibiaDataAPICharacter{}, jsonUnmarshalError
	}
}
