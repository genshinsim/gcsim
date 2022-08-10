package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

//1006934337612173413

type Config struct {
	Token          string
	AdminChannelID string //id for the admin channel
}

type Bot struct {
	cfg Config
}

func (b *Bot) Run(cfg Config) error {

	b.cfg = cfg

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

func (b *Bot) msgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case strings.HasPrefix(m.Content, "!list"):
		if b.adminChanCheck(m) {
			b.List(s)
		}
	case strings.HasPrefix(m.Content, "!approve"):
	case strings.HasPrefix(m.Content, "!reject"):
	case strings.HasPrefix(m.Content, "!submit"):
	}

}

func (b *Bot) adminChanCheck(m *discordgo.MessageCreate) bool {
	return b.cfg.AdminChannelID == m.ChannelID
}
