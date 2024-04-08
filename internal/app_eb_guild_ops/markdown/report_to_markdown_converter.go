package markdown

import (
	"fmt"
	"github.com/adriein/eb-guild-ops/internal/app_eb_guild_ops/handler"
)

type MarkdownConverter struct{}

func NewMarkdownConverter() (*MarkdownConverter, error) {
	return &MarkdownConverter{}, nil
}

func (converter *MarkdownConverter) convertReport(report handler.EbGuildReport) (string, error) {
	markdown := fmt.Sprintf("# Marlock Police State\n### Inactive Members:\n* Elite Fordrin, 315 days (29-05-2023)\n* Ame Damnee, 69 days (30-01-2024)")
}
