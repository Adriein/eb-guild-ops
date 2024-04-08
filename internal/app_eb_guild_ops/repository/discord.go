package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DiscordApi struct {
	token   string
	baseUrl string
	version string
}

type DiscordGuildChannelsResponse []struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

const DISCORD_BOT_TOKEN string = "DISCORD_BOT_TOKEN"
const DISCORD_BASE_URL string = "https://discord.com/api"
const DISCORD_API_VERSION string = "v10"
const DISCORD_API_GUILDS_RESOURCES string = "guilds"
const DISCORD_API_CHANNELS_RESOURCES string = "channels"

func NewDiscordRepository() (*DiscordApi, error) {
	token, hasEnvVar := os.LookupEnv(DISCORD_BOT_TOKEN)

	if !hasEnvVar {
		return nil, errors.New("DISCORD_BOT_TOKEN is not set")
	}

	authorizationToken := fmt.Sprintf("Bot %s", token)

	return &DiscordApi{token: authorizationToken, baseUrl: DISCORD_BASE_URL, version: DISCORD_API_VERSION}, nil
}

func (discord *DiscordApi) FetchChannel(guildID string, name string) (DiscordGuildChannelsResponse, error) {
	request, requestCreationError := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s/%s/%s",
			discord.baseUrl,
			discord.version,
			DISCORD_API_GUILDS_RESOURCES,
			guildID,
			DISCORD_API_CHANNELS_RESOURCES,
		),
		nil,
	)

	request.Header.Set("Authorization", discord.token)

	if requestCreationError != nil {
		return DiscordGuildChannelsResponse{}, requestCreationError
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return DiscordGuildChannelsResponse{}, requestError
	}

	defer response.Body.Close()

	var apiResponse = DiscordGuildChannelsResponse{}

	byteRawResponse, readBodyError := io.ReadAll(response.Body)

	if readBodyError != nil {
		return DiscordGuildChannelsResponse{}, readBodyError
	}

	jsonUnmarshalError := json.Unmarshal(byteRawResponse, &apiResponse)

	if jsonUnmarshalError != nil {
		return DiscordGuildChannelsResponse{}, jsonUnmarshalError
	}

	return apiResponse, nil
}
