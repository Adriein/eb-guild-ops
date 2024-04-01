package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TibiaDataAPICharacter struct {
	Name          string
	Traded        bool
	DeletionDate  string
	LastLogin     time.Time
	AccountStatus string
}

type GuildHall struct {
	Name      string
	PaidUntil time.Time
}

type Member struct {
	Joined time.Time
	Name   string
	Rank   string
	Status string
}

type TibiaDataAPIGuild struct {
	GuildHall    GuildHall
	Members      []Member
	MembersTotal int
}

type ITibiaDataAPI interface {
	Character(name string) (TibiaDataAPICharacter, error)
	Guild(name string) (TibiaDataAPIGuild, error)
}

type TibiaDataAPI struct{}

type TibiaDataAPICharacterJSON struct {
	Character struct {
		Character struct {
			AccountStatus string `json:"account_status"`
			DeletionDate  string `json:"deletion_date"`
			LastLogin     string `json:"last_login"`
			Traded        bool   `json:"traded"`
		} `json:"character"`
	} `json:"character"`
	Information struct {
		Status struct {
			Error    int    `json:"error"`
			HTTPCode int    `json:"http_code"`
			Message  string `json:"message"`
		} `json:"status"`
		Timestamp string `json:"timestamp"`
	} `json:"information"`
}

type TibiaDataAPIGuildJSON struct {
	Guild struct {
		Guildhalls []struct {
			Name      string `json:"name"`
			PaidUntil string `json:"paid_until"`
		} `json:"guildhalls"`
		Members []struct {
			Joined string `json:"joined"`
			Name   string `json:"name"`
			Rank   string `json:"rank"`
			Status string `json:"status"`
		} `json:"members"`
		MembersTotal int `json:"members_total"`
	} `json:"guild"`
	Information struct {
		Status struct {
			Error    int    `json:"error"`
			HTTPCode int    `json:"http_code"`
			Message  string `json:"message"`
		} `json:"status"`
		Timestamp string `json:"timestamp"`
	} `json:"information"`
}

const TIBIA_DATA_API_BASE_URL string = "https://api.tibiadata.com"
const TIBIA_DATA_API_VERSION string = "v4"
const TIBIA_DATA_API_CHARACTER_URL string = "character"
const TIBIA_DATA_API_GUILD_URL string = "guild"

func NewTibiaDataAPI() (*TibiaDataAPI, error) {
	return &TibiaDataAPI{}, nil
}

// Character make a new HTTP request to TibiaAPI and returns a TibiaDataAPICharacter struct
func (api *TibiaDataAPI) Character(name string) (TibiaDataAPICharacter, error) {
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

	lastLogin, timeParseError := time.Parse(time.RFC3339, apiResponse.Character.Character.LastLogin)

	if timeParseError != nil {
		return TibiaDataAPICharacter{}, timeParseError
	}

	return TibiaDataAPICharacter{
		name,
		apiResponse.Character.Character.Traded,
		apiResponse.Character.Character.DeletionDate,
		lastLogin,
		apiResponse.Character.Character.AccountStatus,
	}, nil
}

// Guild make a new HTTP request to TibiaAPI and returns a TibiaDataAPIGuild struct
func (api *TibiaDataAPI) Guild(name string) (TibiaDataAPIGuild, error) {
	request, requestCreationError := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s/%s",
			TIBIA_DATA_API_BASE_URL,
			TIBIA_DATA_API_VERSION,
			TIBIA_DATA_API_GUILD_URL,
			name,
		),
		nil,
	)

	if requestCreationError != nil {
		return TibiaDataAPIGuild{}, requestCreationError
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		return TibiaDataAPIGuild{}, requestError
	}

	defer response.Body.Close()

	var apiResponse = TibiaDataAPIGuildJSON{}

	byteRawResponse, readBodyError := io.ReadAll(response.Body)

	if readBodyError != nil {
		return TibiaDataAPIGuild{}, readBodyError
	}

	jsonUnmarshalError := json.Unmarshal(byteRawResponse, &apiResponse)

	if jsonUnmarshalError != nil {
		return TibiaDataAPIGuild{}, jsonUnmarshalError
	}

	guildHall := apiResponse.Guild.Guildhalls[0]

	paidUntil, paidUntilParseError := time.Parse(time.DateOnly, guildHall.PaidUntil)

	if paidUntilParseError != nil {
		return TibiaDataAPIGuild{}, fmt.Errorf(
			"paid until time parsing error with message: %s",
			paidUntilParseError.Error(),
		)
	}

	var members []Member

	for _, member := range apiResponse.Guild.Members {
		joined, joinedTimeParseError := time.Parse(time.DateOnly, member.Joined)

		if joinedTimeParseError != nil {
			return TibiaDataAPIGuild{}, fmt.Errorf(
				"joined time parsing error with message: %s",
				joinedTimeParseError.Error(),
			)
		}

		members = append(members, Member{joined, member.Name, member.Rank, member.Status})
	}

	return TibiaDataAPIGuild{
		GuildHall{guildHall.Name, paidUntil},
		members,
		apiResponse.Guild.MembersTotal,
	}, nil
}
