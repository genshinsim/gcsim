package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/genshinsim/gcsim/services/pkg/store"
	"go.uber.org/zap"
)

type Config struct {
	Token          string
	AdminChannelID string //id for the admin channel
	DBPath         string
}

type Bot struct {
	cfg   Config
	db    *badger.DB
	Log   *zap.SugaredLogger
	Store store.SimStore
}

type Submission struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Config      string `json:"config"` //this is to be pulled out of the database
}

func Run(cfg Config, s store.SimStore) error {
	b := &Bot{}
	b.cfg = cfg
	b.Store = s

	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	b.Log = sugar

	db, err := badger.Open(badger.DefaultOptions(cfg.DBPath))
	if err != nil {
		return err
	}
	defer db.Close()
	b.db = db

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %v", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(b.msgHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	return nil
}

var reSubmit = regexp.MustCompile(`\!submit.+([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}) (.+)`)
var reApprove = regexp.MustCompile(`\!ok.+([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})`)
var reReject = regexp.MustCompile(`\!reject.+([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}) (.+)`)

func (b *Bot) msgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!list"):
		if b.adminChanCheck(m) {
			b.List(s, m)
		}
	case strings.HasPrefix(m.Content, "!ok"):
		if b.adminChanCheck(m) {
			b.Approve(s, m)
		}
	case strings.HasPrefix(m.Content, "!reject"):
		if b.adminChanCheck(m) {
			b.Reject(s, m)
		}
	case strings.HasPrefix(m.Content, "!submit"):
		b.Submit(s, m)
	}

}

func (b *Bot) adminChanCheck(m *discordgo.MessageCreate) bool {
	return b.cfg.AdminChannelID == m.ChannelID
}

func (b *Bot) Submit(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!submit <link> <description>
	//does not have to be a valid link, just needs to be a valid uuid
	match := reSubmit.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid !submit command")
		return
	}
	author := m.Author.Username + "#" + m.Author.Discriminator

	b.Log.Infow("submission received", "author", author, "key", match[1], "description", match[2])

	sub := Submission{
		Key:         match[1],
		Description: match[2],
		Author:      author,
	}
	data, err := json.Marshal(sub)
	if err != nil {
		b.Log.Warnw("submit - err marshalling json", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(match[1]), data)
		return err
	})

	if err != nil {
		b.Log.Warnw("submit - err updating store", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission recorded! Thanks %v", author))
}

func (b *Bot) List(s *discordgo.Session, m *discordgo.MessageCreate) {
	var subs []Submission
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var x Submission
				err := json.Unmarshal(v, &x)
				if err != nil {
					return err
				}
				subs = append(subs, x)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		b.Log.Warnw("list - err iterating keys", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}
	//make a nice embed?
	embed := &discordgo.MessageEmbed{
		Color: 0x2a8fce,
		Title: "Submitted sims waiting for approval",
	}

	if len(subs) > 0 {
		var sb strings.Builder
		for _, v := range subs {
			sb.WriteString(fmt.Sprintf("%v: [%v](https://next.gcsim.app/v3/viewer/share/%v)\n", v.Author, v.Key, v.Key))
		}
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "Links",
				Value:  sb.String(),
				Inline: true,
			},
		)
	} else {
		embed.Fields = append(embed.Fields,
			&discordgo.MessageEmbedField{
				Name:   "Links",
				Value:  "Nothing here. Yay!",
				Inline: true,
			},
		)
	}

	msg := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	res, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
	if err != nil {
		b.Log.Warnw("list - sending msg", "err", err, "res", res)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
	}
}

func (b *Bot) Reject(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!reject <key> <reason>
	//does not have to be a valid link, just needs to be a valid uuid
	match := reReject.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid !reject command")
		return
	}
	if match[2] == "" {
		s.ChannelMessageSend(m.ChannelID, "Reason required for reject")
		return
	}
	err := b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(match[1]))
		return err
	})

	switch err {
	case badger.ErrKeyNotFound:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v not found", match[1]))
	case nil:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v has been rejected", match[1]))
	default:
		b.Log.Warnw("reject - err deleting key", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
	}

}

func (b *Bot) Approve(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!ok <key>
	//does not have to be a valid link, just needs to be a valid uuid
	match := reApprove.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid !ok command")
		return
	}
	//try grabbing the link from postgrest
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v approved", match[1]))
}
