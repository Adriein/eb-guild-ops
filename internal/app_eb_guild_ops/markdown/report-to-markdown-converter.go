package markdown

import (
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/handler"
	"strings"
)

type Converter struct{}

func NewMarkdownConverter() (*Converter, error) {
	return &Converter{}, nil
}

func (converter *Converter) ConvertReport(report handler.EbGuildReport) (string, error) {
	var reportElements []string

	reportElements = append(reportElements, "# Marlock Police State", "### Inactive Members:")

	for _, inactiveMember := range report.InactiveMembers {
		reportLine := fmt.Sprintf(
			"* %s, %g days (%s)",
			inactiveMember.Name,
			inactiveMember.DaysElapsed,
			inactiveMember.LastLoginDate,
		)

		reportElements = append(reportElements, reportLine)
	}

	markdown := fmt.Sprintf(strings.Join(reportElements, "\n"))

	return markdown, nil
}
