package discord

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func init() {
	commands = append(commands,
		api.CreateCommandData{
			Name:        "mine",
			Description: "list your submissions",
			Options: []discord.CommandOption{
				&discord.NumberOption{
					OptionName:  "page",
					Description: "page number to list, min 1",
					Required:    true,
					Min:         option.NewFloat(1),
				},
			},
		},
		api.CreateCommandData{
			Name:        "delete",
			Description: "request delete of a pending submission",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "id",
					Description: "id of the submission",
					Required:    true,
				},
			},
		},
	)
}

func (b *Bot) cmdListUserSubs(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "mine")
}

func (b *Bot) cmdUserDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "delete")
}
