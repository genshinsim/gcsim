package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/backend/pkg/services/preview"
	"github.com/genshinsim/gcsim/backend/pkg/services/queue"
	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
	"github.com/genshinsim/gcsim/backend/pkg/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sha1ver string
)

func main() {
	setHash()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	keys := getKeys()

	s, err := api.New(api.Config{
		ShareStore:      makeShareStore(),
		UserStore:       makeUserStore(sugar),
		DBStore:         makeDBStore(),
		SubmissionStore: makeSubStore(),
		PreviewStore:    makePreviewStore(),
		Discord: api.DiscordConfig{
			RedirectURL:  os.Getenv("REDIRECT_URL"),
			ClientID:     os.Getenv("DISCORD_ID"),
			ClientSecret: os.Getenv("DISCORD_SECRET"),
			JWTKey:       os.Getenv("JWT_KEY"),
		},
		AESDecryptionKeys: keys,
		MQTTConfig: api.MQTTConfig{
			MQTTUser: os.Getenv("MQTT_USERNAME"),
			MQTTPass: os.Getenv("MQTT_PASSWORD"),
			MQTTHost: os.Getenv("MQTT_URL"),
		},
		QueueService:  makeQueueService(),
		CurrentHash:   sha1ver,
		ComputeAPIKey: os.Getenv("COMPUTE_API_KEY"),
	}, func(s *api.Server) error {
		s.Log = sugar
		return nil
	})

	if err != nil {
		panic(err)
	}

	log.Println("API gateway starting to listen at port 3000")
	log.Fatal(http.ListenAndServe(":3000", s.Router))

}

func setHash() {
	info, _ := debug.ReadBuildInfo()
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			sha1ver = bs.Value
		}
	}
}

func makeShareStore() api.ShareStore {
	shareStore, err := share.NewClient(share.ClientCfg{
		Addr: os.Getenv("SHARE_STORE_URL"),
	})

	if err != nil {
		panic(err)
	}
	return shareStore
}

func makeQueueService() api.QueueService {
	store, err := queue.NewClient(os.Getenv("QUEUE_SERVICE_URL"))
	if err != nil {
		panic(err)
	}
	return store
}

func makeUserStore(sugar *zap.SugaredLogger) api.UserStore {
	store, err := user.New(user.Config{
		DBPath: os.Getenv("USER_DATA_PATH"),
	}, func(s *user.Store) error {
		s.Log = sugar
		return nil
	})

	if err != nil {
		panic(err)
	}

	return store
}

func makeDBStore() api.DBStore {
	store, err := db.NewClient(db.ClientCfg{
		Addr: os.Getenv("DB_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}
	return store
}

func makeSubStore() api.SubmissionStore {
	store, err := submission.NewClient(os.Getenv("SUBMISSION_STORE_URL"))
	if err != nil {
		panic(err)
	}
	return store
}

func makePreviewStore() api.PreviewStore {
	store, err := preview.NewClient(preview.ClientCfg{
		Addr: os.Getenv("PREVIEW_STORE_URL"),
	})
	if err != nil {
		panic(err)
	}
	return store
}

func getKeys() map[string][]byte {
	//read from key file
	var hexKeys map[string]string
	f, err := os.Open(os.Getenv("SHARE_KEY_FILE"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fv, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fv, &hexKeys)
	if err != nil {
		panic(err)
	}

	keys := make(map[string][]byte)
	//convert key from hex string into []byte
	for k, v := range hexKeys {
		key, err := hex.DecodeString(v)
		if err != nil {
			panic("invalid key provided - cannot decode hex to string")
		}
		keys[k] = key
	}

	log.Println("keys read sucessfully: ", hexKeys)
	return keys
}
