package repository

import (
	"fmt"
	"net/http"
	"time"
)

type Date = time.Time

type TibiaDataAPICharacter struct {
	Name      string
	LastLogin Date
}

type TibiaDataApiResponse[T any] struct {
	Data    T
	Success bool
}

type TibiaDataAPI interface {
	Character() TibiaDataApiResponse[TibiaDataAPICharacter]
}

const TIBIA_DATA_API_BASE_URL string = "https://api.tibiadata.com"
const TIBIA_DATA_API_VERSION string = "v4"
const TIBIA_DATA_API_CHARACTER_URL string = "character"
const HTTP_GET_METHOD string = "GET"

func Character(name string) TibiaDataApiResponse[TibiaDataAPICharacter] {
	response, err := http.NewRequest(
		HTTP_GET_METHOD,
		fmt.Sprintf(
			"%s/%s/%s/%s",
			TIBIA_DATA_API_BASE_URL,
			TIBIA_DATA_API_VERSION,
			TIBIA_DATA_API_CHARACTER_URL,
			name,
		),
		nil,
	)

	if err != nil {
		return TibiaDataApiResponse[TibiaDataAPICharacter]{nil, false}
	}

	return TibiaDataApiResponse[TibiaDataAPICharacter]{}
}
