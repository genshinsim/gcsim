package mongo

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	URL        string
	Database   string
	Collection string
	Username   string
	Password   string
}

type Server struct {
	cfg    Config
	client *mongo.Client
	Log    *zap.SugaredLogger
}

func NewServer(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		cfg: cfg,
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
		s.Log.Errorw("result - mongodb ping failed", "err", err)
		return nil, err
	}

	s.Log.Info("result - mongodb connected sucessfully")

	s.client = client

	return s, nil
}

func (s *Server) Create(ctx context.Context, entry *share.ShareEntry) (string, error) {
	s.Log.Infow("create share request", "entry", entry.String())

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	entry.Id = generateID()
	s.Log.Infow("id generated for req ", "id", entry.GetId())

	res, err := col.InsertOne(ctx, entry)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.Log.Infow("create failed - duplicated id (unexpected)", "key", entry.GetId(), "err", err)
			return "", status.Error(codes.Internal, "duplicated id")
		}
		s.Log.Infow("create failed - unexpected error", "key", entry.GetId(), "err", err)
		return "", status.Error(codes.Internal, "internal server error")
	}

	s.Log.Infow("create successful", "id", res.InsertedID)

	return entry.GetId(), nil
}

func (s *Server) Read(ctx context.Context, key string) (*share.ShareEntry, error) {
	s.Log.Infow("get share request", "key", key)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	res := &share.ShareEntry{}
	err := col.FindOne(ctx, bson.D{{Key: "_id", Value: key}}).Decode(res)

	if err != nil {
		s.Log.Infow("error getting share", "err", err)
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "no records found")
		}
		return nil, status.Error(codes.Internal, "unexpected server error")
	}

	s.Log.Infow("get share request done", "key", key)

	return res, nil
}

func (s *Server) Update(ctx context.Context, entry *share.ShareEntry) (string, error) {
	key := entry.GetId()
	s.Log.Infow("update share request", "key", key)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	res, err := col.ReplaceOne(ctx, bson.D{{Key: "_id", Value: key}}, entry)

	if err != nil {
		return "", status.Error(codes.Internal, "unexpected server error")
	}

	if res.MatchedCount == 0 {
		s.Log.Infow("update request failed - no document found", "key", key)
		return "", status.Error(codes.NotFound, "document not found")
	}

	s.Log.Infow("update share request done", "key", key)

	return key, nil
}

func (s *Server) SetTTL(ctx context.Context, key string, until uint64) (string, error) {
	return "not implemented", status.Error(codes.Unimplemented, "set ttl not implemented")
}

func (s *Server) Delete(ctx context.Context, key string) error {
	s.Log.Infow("delete share request", "key", key)

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	res, err := col.DeleteOne(ctx, bson.D{{Key: "_id", Value: key}})

	if err != nil {
		return status.Error(codes.Internal, "unexpected server error")
	}

	if res.DeletedCount == 0 {
		s.Log.Infow("delete request failed - no document found", "key", key)
		return status.Error(codes.NotFound, "document not found")
	}

	s.Log.Infow("delete share request done", "key", key)

	return nil
}

func (s *Server) Random(ctx context.Context) (string, error) {
	s.Log.Infow("get random share request")

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	c, err := col.Aggregate(ctx, bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}})

	if err != nil {
		return "", status.Error(codes.Internal, "unexpected server error")
	}

	res := &share.ShareEntry{}

	err = c.Decode(res)

	if err != nil {
		return "", status.Error(codes.Internal, "unexpected server error")
	}

	s.Log.Infow("get random share request done", "key", res.GetId())

	return res.GetId(), nil
}
