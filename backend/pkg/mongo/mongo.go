package mongo

import (
	"context"
	"os"
	"strconv"

	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var generateID func() string

func init() {
	var err error
	// dictionary from https://github.com/CyberAP/nanoid-dictionary#nolookalikessafe
	generateID, err = nanoid.CustomASCII("6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz", 12)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	URL      string
	Database string
	//collections
	Collection string
	ValidView  string
	SubView    string
	//auth
	Username string
	Password string
	//compute
	CurrentHash string
	BatchSize   int
	Iterations  int
}

type Server struct {
	cfg          Config
	client       *mongo.Client
	Log          *zap.SugaredLogger
	maxPageLimit int64
}

func NewServer(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		cfg:          cfg,
		maxPageLimit: 100,
	}

	limitStr := os.Getenv("MONGO_STORE_MAX_LIMIT")
	if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && limit > 0 {
		s.maxPageLimit = limit
	}

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	credential := options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URL).SetAuth(credential))
	if err != nil {
		return nil, err
	}

	//check connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		s.Log.Errorw("mongodb ping failed", "err", err)
		return nil, err
	}

	//default sanity check
	if s.cfg.BatchSize == 0 {
		s.cfg.BatchSize = 5
	}
	if s.cfg.Iterations == 0 {
		s.cfg.Iterations = 1000
	}

	s.Log.Info("mongodb connected sucessfully")

	s.client = client

	return s, nil
}
