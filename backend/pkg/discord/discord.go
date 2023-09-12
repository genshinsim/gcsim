package discord

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
)

type Backend interface {
	Submit(link, desc, sender string) (string, error)
	GetPending(model.DBTag, int) ([]*db.Entry, error)
	Approve(id string, tag model.DBTag) error
	Reject(id string, tag model.DBTag) error
	RejectAll(model.DBTag) (int64, error)
	GetBySubmitter(id string, page int) ([]*db.Entry, error)
	DeletePending(id, sender string) error
	GetRandomSim() string
	GetDBStatus() (*model.DBStatus, error)
	GetDBEntry(id string) (*db.Entry, error)
	ReplaceConfig(id, link string) error
}

type DBStatus struct {
	ComputeTodoCount int32
	DBTotalCount     int32
}

type Config struct {
	Token      string
	Backend    Backend
	TagMapping map[string]model.DBTag
}

type Bot struct {
	Config
	Log *zap.SugaredLogger
	// discord stuff
	*cmdroute.Router
	s *state.State
}

func New(cfg Config, cust ...func(*Bot) error) (*Bot, error) {
	b := &Bot{
		Config: cfg,
	}

	for _, f := range cust {
		err := f(b)
		if err != nil {
			return nil, err
		}
	}

	if b.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		b.Log = sugar
	}

	if b.Token == "" {
		return nil, errors.New("no token provided")
	}

	if b.Backend == nil {
		return nil, errors.New("no backend provided")
	}

	return b, nil
}

func (b *Bot) Run() error {
	b.s = state.New("Bot " + b.Token)
	err := b.routes()
	if err != nil {
		return nil
	}
	b.s.AddInteractionHandler(b)
	b.s.AddIntents(gateway.IntentGuilds)
	b.s.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := b.s.Me()
		b.Log.Infow("connected", "username", me.Username)
	})

	if err := overwriteCommands(b.s); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// this is blocking
	if err := b.s.Connect(ctx); err != nil {
		return err
	}

	return nil
}

func overwriteCommands(s *state.State) error {
	return cmdroute.OverwriteCommands(s, commands)
}
