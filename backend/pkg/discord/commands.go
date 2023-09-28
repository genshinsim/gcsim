package discord

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

var commands = []api.CreateCommandData{
	{
		Name:        "echo",
		Description: "echo back the argument",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "argument",
				Description: "what's echoed back",
				Required:    true,
			},
		},
	},
}

func (b *Bot) routes() error {
	b.Router = cmdroute.NewRouter()
	// Automatically defer handles if they're slow.
	b.Use(cmdroute.Deferrable(b.s, cmdroute.DeferOpts{}))
	b.AddFunc("echo", b.cmdEcho)
	b.AddFunc("submit", b.cmdSubmit)
	b.AddFunc("list", b.cmdList)
	b.AddFunc("approve", b.cmdApprove)
	b.AddFunc("reject", b.cmdReject)
	b.AddFunc("rejectall", b.cmdRejectAll)
	b.AddFunc("randsim", b.cmdRandom)
	b.AddFunc("mine", b.cmdListUserSubs)
	b.AddFunc("delete", b.cmdUserDelete)
	b.AddFunc("replace", b.cmdReplaceConfig)
	b.AddFunc("reword", b.cmdReplaceDesc)
	b.AddFunc("dbstatus", b.cmdDBStatus)
	b.AddFunc("status", b.cmdEntryStatus)

	return nil
}

func (b *Bot) cmdEcho(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var options struct {
		Arg string `discord:"argument"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content:         option.NewNullableString(options.Arg),
		AllowedMentions: &api.AllowedMentions{}, // don't mention anyone
	}
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}
