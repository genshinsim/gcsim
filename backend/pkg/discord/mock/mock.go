package mock

import (
	"github.com/genshinsim/gcsim/backend/pkg/discord"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockBackend struct {
}

func NewMock() discord.Backend {
	return &MockBackend{}
}

func (*MockBackend) Submit(link, desc, sender string) (string, error) {
	return "", status.Error(codes.InvalidArgument, "bad link")
}
