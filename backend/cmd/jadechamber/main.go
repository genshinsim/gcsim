package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sha1ver string

func main() {
	setHash()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("jadechamber started", "sha1ver", sha1ver)

	db := connectMongo().Database(os.Getenv("MONGODB_DATABASE"))
	s, err := api.New(api.Config{
		ShareStore:  api.NewMongoShareStore(db.Collection(mustGetenv("MONGODB_SHARE_COLLECTION"))),
		DBShareKeys: api.NewMongoDBShareKeys(db.Collection(os.Getenv("MONGODB_COLLECTION"))),
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

func connectMongo() *mongo.Client {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().
			ApplyURI(os.Getenv("MONGODB_URL")).
			SetAuth(options.Credential{
				Username: os.Getenv("MONGODB_USERNAME"),
				Password: os.Getenv("MONOGDB_PASSWORD"),
			}),
	)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}
	return client
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " is not set")
	}
	return v
}
