package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/genshinsim/gcsim/backend/pkg/api"
	"go.uber.org/zap"
)

type Config struct {
	DBPath      string
	ResultStore api.ShareStore
}

type Store struct {
	cfg Config
	db  *badger.DB
	Log *zap.SugaredLogger
}

type User struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Role       int            `json:"role"`
	Permalinks []string       `json:"permalinks,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

const (
	RoleDefaultUser = 0
	RoleDBAdmin     = 50
	RoleSuperAdmin  = 99
)

func New(cfg Config, cust ...func(*Store) error) (*Store, error) {
	s := &Store{
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

	db, err := badger.Open(badger.DefaultOptions(cfg.DBPath))
	if err != nil {
		return nil, err
	}
	s.db = db

	return s, nil
}

func (s *Store) Create(ctx context.Context, id, name string) error {
	key := []byte(id)

	return s.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			s.Log.Infow("creating user - already exists", "id", id)
			return api.ErrUserAlreadyExists
		}

		u := User{
			ID:   id,
			Name: name,
		}
		b, err := json.Marshal(&u)
		if err != nil {
			s.Log.Warnw("error marshalling user", "id", id, "err", err)
		}
		return txn.Set(key, b)
	})
}

func (s *Store) Has(ctx context.Context, id string) (bool, error) {
	has := false
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(id))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		has = true
		return nil
	})
	return has, err
}

func (s *Store) Read(ctx context.Context, id string) ([]byte, error) {
	var data []byte
	err := s.db.View(func(txn *badger.Txn) error {
		x, err := s.getRequester(ctx, txn)
		if err != nil {
			return err
		}
		if x.ID != id && x.Role < RoleSuperAdmin {
			return api.ErrAccessDenied
		}
		key := []byte(id)
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				//TODO: responding with bad request here but not sure if that's ideal..
				s.Log.Infow("bad request; user does not exist", "id", id)
				return api.ErrInvalidRequest
			}
			s.Log.Errorw("unexpected error retrieving requested user", "id", id)
			return err
		}
		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)
			return nil
		})
	})

	return data, err
}

func (s *Store) UpdateData(ctx context.Context, data []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		// as long as we have permission this operation is ok; we don't need to check
		// what's in the data here
		u, err := s.getRequester(ctx, txn)
		if err != nil {
			return err
		}

		var next map[string]any
		err = json.Unmarshal(data, &next)
		if err != nil {
			s.Log.Infow("bad request; update did not supply valid data", "id", u.ID, "err", err)
			return api.ErrInvalidRequest
		}

		u.Data = next

		d, err := json.Marshal(u)
		if err != nil {
			s.Log.Errorw("unexpected error marshalling user back into json", "id", u.ID, "err", err)
			return api.ErrServerError
		}

		err = txn.Set([]byte(u.ID), d)
		if err != nil {
			s.Log.Errorw("unexpected error updating user data", "id", u.ID, "user", u, "err", err)
			return api.ErrServerError
		}
		return nil
	})
}

func (s *Store) getRequester(ctx context.Context, txn *badger.Txn) (User, error) {
	id, err := extractUser(ctx)
	if err != nil {
		s.Log.Infow("bad request; no valid user set in context", "ctx", ctx)
		return User{}, api.ErrInvalidRequest
	}
	return s.getUser(id, txn)
}

func (s *Store) getUser(id string, txn *badger.Txn) (User, error) {
	var u User
	item, err := txn.Get([]byte(id))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.Infow("bad request; requester user does not exist", "user", id)
			return u, api.ErrInvalidRequest
		}
		s.Log.Errorw("unexpected error retrieving requester info", "user", id)
		return u, err
	}
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &u)
	})
	if err != nil {
		s.Log.Errorw("could not decode user json", "id", id, "err", err)
		return u, err
	}

	return u, nil
}

func extractUser(ctx context.Context) (string, error) {
	x := ctx.Value("user")
	if x == nil {
		return "", fmt.Errorf("no user supplied; access denied")
	}
	val, ok := x.(string)
	if !ok {
		return "", fmt.Errorf("invalid user id supplied; access denied")
	}
	return val, nil
}
