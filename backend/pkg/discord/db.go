package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
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
	var opts struct {
		Page float64 `discord:"page"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}
	b.Log.Infow("list request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID, "page", opts.Page)

	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	b.Log.Infow("list request for tag", "tag", tag)

	entries, err := b.Backend.GetPending(tag, int(opts.Page))
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Oops we encountered an error: %v", err)),
		}
	}

	b.Log.Infow("entries received", "len", len(entries))

	if len(entries) == 0 {
		if opts.Page <= 1 {
			return &api.InteractionResponseData{
				Content: option.NewNullableString("No more pending entries!"),
			}
		}
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("No entries found for page %v!", int(opts.Page))),
		}
	}

	embeds := listEmbed(entries, int(opts.Page))

	return &api.InteractionResponseData{
		AllowedMentions: &api.AllowedMentions{
			Users: []discord.UserID{
				data.Event.SenderID(),
			},
		},
		Embeds: &embeds,
	}
}

func listEmbed(entries []*db.Entry, page int) []discord.Embed {
	var result []discord.Embed
	row := discord.NewEmbed()
	for _, v := range entries {
		name := fmt.Sprintf("%v: %v", v.Id, v.Description)
		if len(name) > 254 {
			name = name[:254]
		}
		desc := fmt.Sprintf("<@%v>: https://gcsim.app/sh/%v", v.Submitter, v.ShareKey)
		row.Fields = append(row.Fields, discord.EmbedField{
			Name:  name,
			Value: desc,
		})
	}

	row.Title = fmt.Sprintf("Pending submissions (page %v)", page)
	result = append(result, *row)

	return result
}

func (b *Bot) cmdApprove(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("approve request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	var opts struct {
		Id string `discord:"id"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}

	b.Log.Infow("approve request for tag", "tag", tag, "id", opts.Id)

	err := b.Backend.Approve(opts.Id, tag)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Approve failed due to error: %v", err)),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("%v approved!", opts.Id)),
	}
}

func (b *Bot) cmdReject(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("reject request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)
	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	var opts struct {
		Id string `discord:"id"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}

	b.Log.Infow("reject request for tag", "tag", tag, "id", opts.Id)

	err := b.Backend.Reject(opts.Id, tag)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Reject failed due to error: %v", err)),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("%v rejected!", opts.Id)),
	}
}

func (b *Bot) cmdRejectAll(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("reject all request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)
	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	b.Log.Infow("reject all request for tag", "tag", tag)

	count, err := b.Backend.RejectAll(tag)
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString(fmt.Sprintf("Reject failed due to error: %v", err)),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("%v entries rejected!", count)),
	}
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
	b.Log.Infow("db status request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	s, err := b.Backend.GetDBStatus()
	if err != nil {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Sorry! I encountered an error"),
		}
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("There are a total of %v entries in the database, including unapproved entries. %v is pending simulation run.", s.DbTotalCount, s.ComputeCount)),
	}
}

func (b *Bot) cmdEntryStatus(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("entry status request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	var opts struct {
		Ids string `discord:"ids"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}

	b.Log.Infow("entry status request for id", "id", opts.Ids)

	var embeds []discord.Embed
	row := discord.NewEmbed()
	row.Title = "DB Entry Status"
	row.Description = fmt.Sprintf("<@%v>, here are the status of the entries as requested", data.Event.SenderID())

	ids := strings.Split(opts.Ids, ",")
	for _, id := range ids {
		v, err := b.Backend.GetDBEntry(id)
		var desc string
		if err != nil {
			desc = fmt.Sprintf("Error retrieving entry: %v", err)
		} else {
			status := "UNKNOWN"
			link := "Unknown submission; No link"
			switch {
			case v.Summary == nil:
				status = "Pending Compute"
				link = "No link available yet"
			case !v.IsDbValid:
				link = fmt.Sprintf("https://gcsim.app/sh/%v", v.ShareKey)
				status = "Pending Review"
			case v.IsDbValid:
				link = fmt.Sprintf("https://gcsim.app/db/%v", v.Id)
				status = "Added"
			}

			desc = fmt.Sprintf("%v: %v", status, link)
		}

		row.Fields = append(row.Fields, discord.EmbedField{
			Name:  id,
			Value: desc,
		})
	}

	embeds = append(embeds, *row)

	return &api.InteractionResponseData{
		AllowedMentions: &api.AllowedMentions{
			Users: []discord.UserID{
				data.Event.SenderID(),
			},
		},
		Embeds: &embeds,
	}
}

func (b *Bot) cmdReplaceConfig(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("replace config request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	var opts struct {
		Link string `discord:"link"`
		Id   string `discord:"id"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}
	b.Log.Infow("replace options", "opts", opts)

	err := b.Backend.ReplaceConfig(opts.Id, opts.Link, tag)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("DB entry with id %v has been updated", opts.Id)),
	}
}

func (b *Bot) cmdReplaceDesc(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	b.Log.Infow("replace desc request received", "from", data.Event.Sender().Username, "channel", data.Event.ChannelID)

	channelId := data.Event.ChannelID.String()
	if data.Event.Channel != nil {
		if data.Event.Channel.Type == discord.GuildPublicThread || data.Event.Channel.Type == discord.GuildPrivateThread {
			channelId = data.Event.Channel.ParentID.String()
		}
	}

	tag, ok := b.TagMapping[channelId]
	if !ok {
		return &api.InteractionResponseData{
			Content: option.NewNullableString("Oops you don't have permission to do this"),
		}
	}

	var opts struct {
		Desc string `discord:"desc"`
		Id   string `discord:"id"`
	}
	if err := data.Options.Unmarshal(&opts); err != nil {
		return errorResponse(err)
	}
	b.Log.Infow("replace options", "opts", opts)

	err := b.Backend.ReplaceDesc(opts.Id, opts.Desc, tag)
	if err != nil {
		return errorResponse(err)
	}

	return &api.InteractionResponseData{
		Content: option.NewNullableString(fmt.Sprintf("Description for DB entry with id %v has been reworded", opts.Id)),
	}
}
