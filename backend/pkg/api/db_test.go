package api

import (
	"context"

	"github.com/genshinsim/gcsim/pkg/model"
)

type mockDBStore struct {
}

func (m *mockDBStore) Get(context.Context, *model.DBQueryOpt) (*model.DBEntries, error) {
	return nil, nil
}
func (m *mockDBStore) GetOne(ctx context.Context, id string) (*model.DBEntry, error) {
	return nil, nil
}
func (m *mockDBStore) Update(ctx context.Context, id string, result *model.SimulationResult) error {
	return nil
}
func (m *mockDBStore) ApproveTag(context.Context, string, model.DBTag) error {
	return nil
}
func (m *mockDBStore) RejectTag(context.Context, string, model.DBTag) error {
	return nil
}
