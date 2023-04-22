package discord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
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
	var opts struct {
		Page float64 `discord:"page"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}
	b.Log.Infow("list user sims request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID, "page", opts.Page)

	entries, err := b.Backend.GetBySubmitter(data.Event.SenderID().String(), int(opts.Page))
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Oops we encountered an error: %v", err)),
		}
	}

	b.Log.Infow("user entries received", "len", len(entries))

	if len(entries) == 0 {
		if opts.Page <= 1 {
			return &api.InteractionResponseData{
				Content: option.NewNullableString("No submissions found :( Klee is sad"),
			}
		}
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("No entries found for page %v!", int(opts.Page))),
		}
	}

	embeds := userSubEmbeds(entries, int(opts.Page), data.Event.SenderID())

	return &api.InteractionResponseData{
		AllowedMentions: &api.AllowedMentions{
			Users: []discord.UserID{
				data.Event.SenderID(),
			},
		},
		Embeds: &embeds,
	}
}

func userSubEmbeds(entries []*db.Entry, page int, sender discord.UserID) []discord.Embed {
	var result []discord.Embed
	row := discord.NewEmbed()
	row.Title = fmt.Sprintf("DB Submissions (Page %v)", page)
	row.Description = fmt.Sprintf("<@%v>, here's all of your db submissions as requested", sender)
	for _, v := range entries {
		//TODO: insert create date
		title := fmt.Sprintf("%v", v.Description)
		if len(title) > 254 {
			title = title[:254]
		}
		status := "UNKNOWN"
		switch {
		case v.Summary == nil:
			status = "Pending Compute"
		case !v.IsDbValid:
			status = "Pending Review"
		case v.IsDbValid:
			status = "Added"
		}

		desc := fmt.Sprintf("id: %v (%v):\nhttps://simimpact.app/viewer/share/%v", v.Id, status, v.ShareKey)
		row.Fields = append(row.Fields, discord.EmbedField{
			Name:  title,
			Value: desc,
		})
	}

	result = append(result, *row)

	return result
}

func (b *Bot) cmdUserDelete(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("user delete request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	var opts struct {
		Id string `discord:"id"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}

	b.Log.Infow("delete request", "user", data.Event.SenderID(), "id", opts.Id)

	err := b.Backend.DeletePending(opts.Id, data.Event.SenderID().String())
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Delete failed due to error: %v", err)),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Submission id %v has been deleted!", opts.Id)),
	}
}
