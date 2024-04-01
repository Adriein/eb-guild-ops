package handler

import (
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/repository"
	"math"
	"time"
)

type CreateReportCommandParameters struct {
	GuildName string
}

type CreateReportCommand struct {
	Repository repository.ITibiaDataAPI
	Params     CreateReportCommandParameters
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

func Execute(command CreateReportCommand) (EbGuildReport, error) {
	var tibiaDataAPI = command.Repository
	var guildName = command.Params.GuildName

	guild, guildError := tibiaDataAPI.Guild(guildName)

	if guildError != nil {
		return EbGuildReport{}, guildError
	}

	var inactiveMembers []InactiveMember

	for _, member := range guild.Members {
		memberDetail, characterError := tibiaDataAPI.Character(member.Name)

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
