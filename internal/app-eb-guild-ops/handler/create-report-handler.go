package handler

import (
	"github.com/adriein/eb-guild-ops/internal/app-eb-guild-ops/repository"
	"math"
	"time"
)

type CreateReportHandler struct {
	Api repository.ITibiaDataAPI
}

type CreateReportCommand struct {
	GuildName string
}

type InactiveMember struct {
	Name          string
	LastLoginDate string
	DaysElapsed   float64
}

type EbGuildReport struct {
	Version                   string
	InactiveMembers           []InactiveMember
	InactiveMembersNumber     int
	MembersNumber             int
	NextGuildHousePaymentDate string
}

func NewCreateReportHandler(api repository.ITibiaDataAPI) (*CreateReportHandler, error) {
	return &CreateReportHandler{api}, nil
}

func (handler *CreateReportHandler) Execute(command CreateReportCommand) (EbGuildReport, error) {
	var guildName = command.GuildName

	guild, guildError := handler.Api.Guild(guildName)

	if guildError != nil {
		return EbGuildReport{}, guildError
	}

	var inactiveMembers []InactiveMember

	for _, member := range guild.Members {
		memberDetail, characterError := handler.Api.Character(member.Name)

		if characterError != nil {
			return EbGuildReport{}, characterError
		}

		nano := time.Since(memberDetail.LastLogin)

		days := math.RoundToEven(nano.Hours() / 24)

		if days >= 30 {
			inactiveMembers = append(
				inactiveMembers,
				InactiveMember{
					memberDetail.Name,
					memberDetail.LastLogin.Format("02-01-2006"),
					days,
				},
			)
		}

	}

	return EbGuildReport{
		"1.0",
		inactiveMembers,
		len(inactiveMembers),
		len(guild.Members),
		guild.GuildHall.PaidUntil.Format("02-01-2006"),
	}, nil
}
