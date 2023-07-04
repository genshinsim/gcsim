package mocks

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/jaevor/go-nanoid"
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

type ShareStore struct {
	data map[string]*model.SimulationResult
}

func NewMockShareStore() api.ShareStore {
	return &ShareStore{
		data: make(map[string]*model.SimulationResult),
	}
}

func (s *ShareStore) SetTTL(ctx context.Context, id string) error {
	return nil
}

func (s *ShareStore) Delete(ctx context.Context, id string) error {
	delete(s.data, id)
	return nil
}

func (s *ShareStore) Random(ctx context.Context) (string, error) {
	for k := range s.data {
		return k, nil
	}
	return "", status.Error(codes.NotFound, "no records found")
}

func (s *ShareStore) Read(ctx context.Context, id string) (*model.SimulationResult, uint64, error) {
	v, ok := s.data[id]
	if !ok {
		return nil, 0, status.Error(codes.NotFound, "record not foudn")
	}
	return v, 0, nil // never expires
}

func (s *ShareStore) Create(ctx context.Context, data *model.SimulationResult, ttl uint64, user string) (string, error) {
	id := generateID()
	if _, ok := s.data[id]; ok {
		return "", status.Error(codes.Internal, "error generating id - got collision")
	}
	s.data[id] = data
	return id, nil
}
