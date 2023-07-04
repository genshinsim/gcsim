package share

import (
	"context"
	"math/rand"
	"time"

	"github.com/jaevor/go-nanoid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
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

// mockServer is a mock server for purpose of testing share RPC end points
type mockServer struct {
	Log  *zap.SugaredLogger
	Rand *rand.Rand
	data map[string]*ShareEntry
}

func newMock(cust ...func(*mockServer) error) (*mockServer, error) {
	s := &mockServer{
		data: make(map[string]*ShareEntry),
	}
	s.Rand = rand.New(rand.NewSource(time.Now().Unix()))

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

	return s, nil
}

func (s *mockServer) Create(ctx context.Context, e *ShareEntry) (string, error) {
	key := generateID()
	if _, ok := s.data[key]; ok {
		return "", status.Error(codes.Internal, "error creating nanoid")
	}
	s.data[key] = e
	return key, nil
}

func (s *mockServer) Read(ctx context.Context, key string) (*ShareEntry, error) {
	val, ok := s.data[key]
	if !ok {
		return nil, status.Error(codes.NotFound, "key not found")
	}
	n := proto.Clone(val)
	return n.(*ShareEntry), nil
}

func (s *mockServer) Update(ctx context.Context, entry *ShareEntry) (string, error) {
	key := entry.GetId()
	_, ok := s.data[key]
	if !ok {
		return "", status.Error(codes.NotFound, "key not found")
	}
	s.data[key] = entry
	return key, nil
}

func (s *mockServer) SetTTL(ctx context.Context, key string, until uint64) (string, error) {
	_, ok := s.data[key]
	if !ok {
		return "", status.Error(codes.NotFound, "key not found")
	}
	return key, nil
}

func (s *mockServer) Delete(ctx context.Context, key string) error {
	_, ok := s.data[key]
	if !ok {
		return status.Error(codes.NotFound, "key not found")
	}
	delete(s.data, key)
	return nil
}

func (s *mockServer) Random(context.Context) (string, error) {
	max := len(s.data)
	if max == 0 {
		return "", status.Error(codes.NotFound, "not found")
	}
	n := s.Rand.Intn(max)
	for k := range s.data {
		if n == 0 {
			return k, nil
		}
		n--
	}

	return "", status.Error(codes.NotFound, "not found")
}
