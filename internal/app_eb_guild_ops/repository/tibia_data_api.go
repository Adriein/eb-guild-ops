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

type TibiaDataAPI interface {
	Character() (TibiaDataAPICharacter, error)
}

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
		API struct {
			Commit  string `json:"commit"`
			Release string `json:"release"`
			Version int    `json:"version"`
		} `json:"api"`
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

// Character make a new request to TibiaAPI and returns a TibiaDataAPICharacter struct
func Character(name string) (TibiaDataAPICharacter, error) {
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
