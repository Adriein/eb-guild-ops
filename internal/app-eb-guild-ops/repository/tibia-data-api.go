package repository

import (
	"encoding/json"
	"fmt"
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

const (
	TibiaDataApiBaseUrl      string = "https://api.tibiadata.com"
	TibiaDataApiVersion      string = "v4"
	TibiaDataApiCharacterUrl string = "character"
	TibiaDataApiGuildUrl     string = "guild"
)

func NewTibiaDataAPI() (*TibiaDataAPI, error) {
	return &TibiaDataAPI{}, nil
}

// Character make a new HTTP request to TibiaAPI and returns a TibiaDataAPICharacter struct
func (api *TibiaDataAPI) Character(name string) (TibiaDataAPICharacter, error) {
	request, requestCreationError := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/%s/%s/%s",
			TibiaDataApiBaseUrl,
			TibiaDataApiVersion,
			TibiaDataApiCharacterUrl,
			name,
		),
		nil,
	)

	if requestCreationError != nil {
		err := fmt.Errorf(
			"> Function: Character\n > Error: RequestCreationError -> %s\n",
			requestCreationError.Error(),
		)

		return TibiaDataAPICharacter{}, err
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		err := fmt.Errorf(
			"> Function: Character\n > Error: RequestError -> %s\n",
			requestError.Error(),
		)

		return TibiaDataAPICharacter{}, err
	}

	defer response.Body.Close()

	var apiResponse = TibiaDataAPICharacterJSON{}

	if decodeError := json.NewDecoder(response.Body).Decode(&apiResponse); decodeError != nil {
		err := fmt.Errorf(
			"> Function: Character\n > Error: DecodeError -> %s\n",
			decodeError.Error(),
		)

		return TibiaDataAPICharacter{}, err
	}

	lastLogin, timeParseError := time.Parse(time.RFC3339, apiResponse.Character.Character.LastLogin)

	if timeParseError != nil {
		err := fmt.Errorf(
			"> Function: Character\n > Error: TimeParseError -> %s\n",
			timeParseError.Error(),
		)

		return TibiaDataAPICharacter{}, err
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
			TibiaDataApiBaseUrl,
			TibiaDataApiVersion,
			TibiaDataApiGuildUrl,
			name,
		),
		nil,
	)

	if requestCreationError != nil {
		err := fmt.Errorf(
			"> Function: Guild\n > Error: RequestCreationError -> %s\n",
			requestCreationError.Error(),
		)

		return TibiaDataAPIGuild{}, err
	}

	client := &http.Client{}
	response, requestError := client.Do(request)

	if requestError != nil {
		err := fmt.Errorf(
			"> Function: Guild\n > Error: RequestError -> %s\n",
			requestError.Error(),
		)

		return TibiaDataAPIGuild{}, err
	}

	defer response.Body.Close()

	var apiResponse = TibiaDataAPIGuildJSON{}

	if decodeError := json.NewDecoder(response.Body).Decode(&apiResponse); decodeError != nil {
		err := fmt.Errorf(
			"> Function: Guild\n > Error: DecodeError -> %s\n",
			decodeError.Error(),
		)

		return TibiaDataAPIGuild{}, err
	}

	guildHall := apiResponse.Guild.Guildhalls[0]

	paidUntil, paidUntilParseError := time.Parse(time.DateOnly, guildHall.PaidUntil)

	if paidUntilParseError != nil {
		err := fmt.Errorf(
			"> Function: Guild\n > Error: PaidUntilParseError -> %s\n",
			paidUntilParseError.Error(),
		)

		return TibiaDataAPIGuild{}, err
	}

	var members []Member

	for _, member := range apiResponse.Guild.Members {
		joined, joinedTimeParseError := time.Parse(time.DateOnly, member.Joined)

		if joinedTimeParseError != nil {
			err := fmt.Errorf(
				"> Function: Guild\n > Error: JoinedTimeParseError -> %s\n",
				joinedTimeParseError.Error(),
			)

			return TibiaDataAPIGuild{}, err
		}

		members = append(members, Member{joined, member.Name, member.Rank, member.Status})
	}

	return TibiaDataAPIGuild{
		GuildHall{guildHall.Name, paidUntil},
		members,
		apiResponse.Guild.MembersTotal,
	}, nil
}
