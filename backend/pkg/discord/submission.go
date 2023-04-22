package discord

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	var opts struct {
		Link string `discord:"link"`
		Desc string `discord:"desc"`
	}
	b.Log.Infow("submission received", "from", data.Event.Sender().Username)
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}
	b.Log.Infow("submission options", "opts", opts)

	if data.Event.SenderID() == 0 {
		b.Log.Info("unexpected sender id is 0")
		return &api.InteractionResponseData{
			Content: &option.NullableStringData{
				Val: "Command failed; Could not find sender information",
			},
		}
	}

	id, err := b.Backend.Submit(opts.Link, opts.Desc, data.Event.SenderID().String())
	resp := discord.NewEmbed()
	if err != nil {
		resp.Title = "Submission failed"
		//catch all desc
		resp.Description = "An unexpected error occured. Please contact administrator"
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				resp.Description = "The viewer link you have provided is either invalid or an error occured reading it."
			case codes.NotFound:
				resp.Description = "Could not find the viewer share provided"
			}
		} else {
			b.Log.Warnw("unexpected error submitting; error not a status", "err", err)
		}

	} else {
		resp.Title = fmt.Sprintf("Submission recorded successfully (id: %v)", id)
		resp.URL = opts.Link
		resp.Description = opts.Desc
	}

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
