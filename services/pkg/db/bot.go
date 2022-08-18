package db

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
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

type SimStoreWithDB interface {
	store.SimStore
	store.SimDBStore
}

type Bot struct {
	cfg   Config
	db    *badger.DB
	Log   *zap.SugaredLogger
	Store SimStoreWithDB
}

func Run(cfg Config, s SimStoreWithDB) error {
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
var reReplace = regexp.MustCompile(`(?m)\!replace.+([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}).+([0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})`)

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
	case strings.HasPrefix(m.Content, "!replace"):
		if b.adminChanCheck(m) {
			b.Replace(s, m)
		}
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
	discordID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	if err != nil {
		b.Log.Warnw("submit - err decoding user id to int64", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}
	sub := store.DBEntry{
		Key:          match[1],
		Description:  match[2],
		AuthorString: m.Author.Username + "#" + m.Author.Discriminator,
		Author:       discordID,
	}
	b.Log.Infow("submission received", "author", sub.Author, "key", sub.Key, "description", sub.Description)

	//make sure sim is exist. we want to save just the config for rerunning later
	sim, err := b.Store.Fetch(sub.Key)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error retrieving submitted sim: %v", sub.Key))
		b.Log.Infow("error retrieiving submitted sim", "key", sub.Key, "err", err)
		return
	}

	res, err := sim.DecodeViewer()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Internal server error decoding submitted sim: %v", sub.Key))
		b.Log.Infow("error decoding submitted sim", "key", sub.Key, "err", err)
		return
	}

	sub.HashedConfig = res.Config

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
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission recorded! Thanks %v", sub.Author))
}

func (b *Bot) List(s *discordgo.Session, m *discordgo.MessageCreate) {
	var subs []store.DBEntry
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var x store.DBEntry
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
			sb.WriteString(fmt.Sprintf("%v - %v: [%v](https://gcsim.app/v3/viewer/share/%v)\n", v.Author, v.Description, v.Key, v.Key))
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

func (b *Bot) retrieveSubmission(key string) (store.DBEntry, error) {
	var data store.DBEntry

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &data)
		})
		return err
	})

	return data, err
}

func (b *Bot) Approve(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!ok <key>
	//does not have to be a valid link, just needs to be a valid uuid
	match := reApprove.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid !ok command")
		return
	}
	sub, err := b.retrieveSubmission(match[1])
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v not found", match[1]))
			return
		}
		b.Log.Warnw("approve - err retrieving key", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}

	id, err := b.Store.Add(sub)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v approval failed with error: %v", match[1], err))
		return
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(match[1]))
		return err
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleting submission %v after approval, not found", match[1]))
			return
		}
		b.Log.Warnw("approve - err deleting key", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v sucessfully added to db with id %v", match[1], id))
}

func (b *Bot) Replace(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!replace <new> <old>
	//replace existing submission with new
	match := reReplace.FindStringSubmatch(m.Content)
	if len(match) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Invalid !ok command")
		return
	}
	sub, err := b.retrieveSubmission(match[1])
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v not found", match[1]))
			return
		}
		b.Log.Warnw("approve - err retrieving key", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}

	id, err := b.Store.Replace(match[2], sub)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v replacement for %v failed with error: %v", match[1], match[2], err))
		return
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(match[1]))
		return err
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleting submission %v after approval, not found", match[1]))
			return
		}
		b.Log.Warnw("approve - err deleting key", "err", err)
		s.ChannelMessageSend(m.ChannelID, "Internal server error processing request")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Submission %v sucessfully replaced %v in db with id %v", match[1], match[2], id))
}

func (b *Bot) ShowConfig(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!showconfig <key>
	//replace existing submission with new
}

func (b *Bot) Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	//!dbhelp
}
