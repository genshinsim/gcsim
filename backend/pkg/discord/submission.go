package discord

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
)

func init() {
	commands = append(commands, api.CreateCommandData{
		Name:        "submit",
		Description: "submit a sim to the db",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "link",
				Description: "sim viewer link",
				Required:    true,
			},
			&discord.StringOption{
				OptionName:  "desc",
				Description: "description of the sim",
				Required:    true,
			},
		},
	})
}

func (b *Bot) cmdSubmit(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("submission received while db archiving", "from", data.Event.Sender().Username)

	resp := discord.NewEmbed()
	resp.Title = "Submissions are closed"
	resp.Description = "The DB is in the process of being archived and as such we are not taking any new submissions. Stay tuned for more updates."

	e := []discord.Embed{*resp}

	return &api.InteractionResponseData{
		AllowedMentions: &api.AllowedMentions{
			Users: []discord.UserID{
				data.Event.SenderID(),
			},
		},
		Embeds: &e,
	}
}
