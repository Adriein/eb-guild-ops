package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DiscordApi struct {
	token   string
	baseUrl string
	version string
}

type DiscordGuildChannelsAPIResponse []struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FetchDiscordChannelResponse struct {
	Id   string
	Name string
}

type SendMessageDiscordAPIBody struct {
	Content string `json:"content"`
}

const (
	DiscordBotToken             string = "DISCORD_BOT_TOKEN"
	DiscordBaseUrl              string = "https://discord.com/api"
	DiscordApiVersion           string = "v10"
	DiscordApiGuildsResources   string = "guilds"
	DiscordApiChannelsResources string = "channels"
	DiscordApiMessagesResources string = "messages"
)

func NewDiscordRepository() (*DiscordApi, error) {
	token, isSet := os.LookupEnv(DiscordBotToken)

	if !isSet {
		return nil, fmt.Errorf("> Function: NewDiscordRepository\n > Error: DISCORD_BOT_TOKEN is not set\n")
	}

	authorizationToken := fmt.Sprintf("Bot %s", token)

	return &DiscordApi{token: authorizationToken, baseUrl: DiscordBaseUrl, version: DiscordApiVersion}, nil
}

func (discord *DiscordApi) FetchChannel(guildID string, name string) (FetchDiscordChannelResponse, error) {
	request, requestCreationError := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s/%s/%s",
			discord.baseUrl,
			discord.version,
			DiscordApiGuildsResources,
			guildID,
			DiscordApiChannelsResources,
		),
		nil,
	)

	request.Header.Set("Authorization", discord.token)

	if requestCreationError != nil {
		err := fmt.Errorf(
			"> Function: FetchChannel\n > Error: RequestCreationError -> %s\n",
			requestCreationError.Error(),
		)

		return FetchDiscordChannelResponse{}, err
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	defer response.Body.Close()

	if requestError != nil {
		err := fmt.Errorf("> Function: FetchChannel\n > Error: RequestError -> %s\n", requestCreationError.Error())

		return FetchDiscordChannelResponse{}, err
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf(
			"> Function: FetchChannel\n > Error: RequestError -> Discord API responded with http code: %d\n",
			response.StatusCode,
		)

		return FetchDiscordChannelResponse{}, err
	}

	var apiResponse = DiscordGuildChannelsAPIResponse{}

	if decodeError := json.NewDecoder(response.Body).Decode(&apiResponse); decodeError != nil {
		err := fmt.Errorf("> Function: FetchChannel\n > Error: DecodeBodyError -> %s\n", decodeError.Error())

		return FetchDiscordChannelResponse{}, err
	}

	for _, channel := range apiResponse {
		if channel.Name == name {
			return FetchDiscordChannelResponse{channel.Id, channel.Name}, nil
		}
	}

	err := fmt.Errorf("> Function: FetchChannel\n > Error: Channel with name %s not found\n", name)

	return FetchDiscordChannelResponse{}, err
}

func (discord *DiscordApi) Message(channelId string, message string) error {
	request, requestCreationError := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/%s/%s/%s/%s",
			discord.baseUrl,
			discord.version,
			DiscordApiChannelsResources,
			channelId,
			DiscordApiMessagesResources,
		),
		nil,
	)

	request.Header.Set("Authorization", discord.token)

	if requestCreationError != nil {
		err := fmt.Errorf(
			"> Function: FetchChannel\n > Error: RequestCreationError -> %s\n",
			requestCreationError.Error(),
		)

		return err
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	defer response.Body.Close()

	if requestError != nil {
		err := fmt.Errorf("> Function: FetchChannel\n > Error: RequestError -> %s\n", requestCreationError.Error())

		return err
	}

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf(
			"> Function: FetchChannel\n > Error: RequestError -> Discord API responded with http code: %d\n",
			response.StatusCode,
		)

		return err
	}

	return nil
}
