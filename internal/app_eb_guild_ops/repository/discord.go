package repository

import (
	"encoding/json"
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

const DISCORD_BOT_TOKEN string = "DISCORD_BOT_TOKEN"
const DISCORD_BASE_URL string = "https://discord.com/api"
const DISCORD_API_VERSION string = "v10"
const DISCORD_API_GUILDS_RESOURCES string = "guilds"
const DISCORD_API_CHANNELS_RESOURCES string = "channels"
const DISCORD_API_MESSAGES_RESOURCES string = "messages"

func NewDiscordRepository() (*DiscordApi, error) {
	token, isSet := os.LookupEnv(DISCORD_BOT_TOKEN)

	if !isSet {
		return nil, fmt.Errorf("> Function: NewDiscordRepository\n > Error: DISCORD_BOT_TOKEN is not set\n")
	}

	authorizationToken := fmt.Sprintf("Bot %s", token)

	return &DiscordApi{token: authorizationToken, baseUrl: DISCORD_BASE_URL, version: DISCORD_API_VERSION}, nil
}

func (discord *DiscordApi) FetchChannel(guildID string, name string) (FetchDiscordChannelResponse, error) {
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

	byteRawResponse, readBodyError := io.ReadAll(response.Body)

	if readBodyError != nil {
		err := fmt.Errorf("> Function: FetchChannel\n > Error: ReadBodyError -> %s\n", readBodyError.Error())

		return FetchDiscordChannelResponse{}, err
	}

	jsonUnmarshalError := json.Unmarshal(byteRawResponse, &apiResponse)

	if jsonUnmarshalError != nil {
		err := fmt.Errorf("> Function: FetchChannel\n > Error: JsonUnmarshalError -> %s\n", jsonUnmarshalError.Error())

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
			DISCORD_API_CHANNELS_RESOURCES,
			channelId,
			DISCORD_API_MESSAGES_RESOURCES,
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
