package discord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func init() {
	commands = append(commands,
		api.CreateCommandData{
			Name:        "dbstatus",
			Description: "return current db status",
		},
		api.CreateCommandData{
			Name:        "list",
			Description: "list pending sims",
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
			Name:        "randsim",
			Description: "give me a random sim link!",
		},
		api.CreateCommandData{
			Name:        "approve",
			Description: "approve sim",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "id",
					Description: "id of the entry",
					Required:    true,
				},
			},
		},
		api.CreateCommandData{
			Name:        "reject",
			Description: "reject sim",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "id",
					Description: "id of the entry",
					Required:    true,
				},
				&discord.StringOption{
					OptionName:  "reason",
					Description: "reason for rejection (will ping submitter if not blank)",
					Required:    false,
				},
			},
		},
		api.CreateCommandData{
			Name:        "rejectall",
			Description: "reject all unapproved sim",
		},
		api.CreateCommandData{
			Name:        "replace",
			Description: "replace sim config (admin only)",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "id",
					Description: "id of the entry",
					Required:    true,
				},
				&discord.StringOption{
					OptionName:  "link",
					Description: "viewer link of new config",
					Required:    true,
				},
			},
		},
		api.CreateCommandData{
			Name:        "reword",
			Description: "reword entry desc (admin only)",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "id",
					Description: "id of the entry",
					Required:    true,
				},
				&discord.StringOption{
					OptionName:  "desc",
					Description: "new description",
					Required:    true,
				},
			},
		},
		api.CreateCommandData{
			Name:        "status",
			Description: "status db entry",
			Options: []discord.CommandOption{
				&discord.StringOption{
					OptionName:  "ids",
					Description: "comma separated list of ids of entry to look up",
					Required:    true,
				},
			},
		},
	)
}

func (b *Bot) cmdList(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "list")
}

func (b *Bot) cmdApprove(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "approve")
}

func (b *Bot) cmdReject(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "reject")
}

func (b *Bot) cmdRejectAll(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "rejectall")
}

func (b *Bot) cmdRandom(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("random sim request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	id := b.Backend.GetRandomSim()

	if id == "" {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Sorry! I couldn't find anything :("),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Here you go: https://gcsim.app/sh/%v", id)),
	}
}

func (b *Bot) cmdDBStatus(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "dbstatus")
}

func (b *Bot) cmdEntryStatus(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "status")
}

func (b *Bot) cmdReplaceConfig(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "replace")
}

func (b *Bot) cmdReplaceDesc(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	return b.cmdDisabled(data, "reword")
}
