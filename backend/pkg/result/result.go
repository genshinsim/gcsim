package result

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/pb"
	"github.com/dgraph-io/ristretto/z"
	"github.com/genshinsim/gcsim/backend/pkg/api"
	"github.com/jaevor/go-nanoid"
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
	DBPath string
}

type Store struct {
	cfg Config
	db  *badger.DB
	Log *zap.SugaredLogger
}

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

func (s *Store) Create(data []byte, ctx context.Context) (string, error) {
	id := generateID()
	ttl := extractTTL(ctx)
	s.Log.Infow("received create request", "id", id, "ttl", ttl)

	err := s.db.Update(func(txn *badger.Txn) error {
		key := []byte(id)

		//sanity check that uuid is not already used...
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			s.Log.Warnw("unexpected id collision", "id", id)
			return fmt.Errorf("unexpected key already exists: %v", id)
		}

		if ttl > 0 {
			e := badger.NewEntry(key, data).WithTTL(time.Hour * time.Duration(ttl))
			return txn.SetEntry(e)
		}
		return txn.Set(key, data)
	})

	return id, err
}

func (s *Store) Random(ctx context.Context) ([]byte, error) {
	var keys [][]byte
	count := 0
	stream := s.db.NewStream()
	stream.NumGo = 16

	// overide stream.KeyToList as we only want keys. Also
	// we can take only first version for the key.
	stream.KeyToList = func(key []byte, itr *badger.Iterator) (*pb.KVList, error) {
		l := &pb.KVList{}
		// Since stream framework copies the item's key while calling
		// KeyToList, we can directly append key to list.
		l.Kv = append(l.Kv, &pb.KV{Key: key})
		return l, nil
	}

	// The bigger the sample size, the more randomness in the outcome.
	sampleSize := 1000
	c, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream.Send = func(buf *z.Buffer) error {
		l, err := badger.BufferToKVList(buf)
		if err != nil {
			s.Log.Infow("error converting buffer to list", "err", err)
			return nil
		}
		if count >= sampleSize {
			return nil
		}
		// Collect "keys" equal to sample size
		for _, kv := range l.Kv {
			keys = append(keys, kv.Key)
			count++
			if count >= sampleSize {
				cancel()
				return nil
			}
		}
		return nil
	}

	if err := stream.Orchestrate(c); err != nil && err != context.Canceled {
		panic(err)
	}
	s.Log.Infow("random key selection done", "len", len(keys))
	if len(keys) == 0 {
		return nil, api.ErrKeyNotFound
	}
	// Pick a random key from the list of keys
	n := keys[rand.Intn(len(keys))]

	return s.Read(string(n), ctx)
}

func (s *Store) Read(key string, ctx context.Context) ([]byte, error) {
	var res []byte
	s.Log.Infow("received request to view result", "key", key)
	err := s.db.Update(func(txn *badger.Txn) error {
		k := []byte(key)
		//sanity check that uuid is not already used...
		item, err := txn.Get(k)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				s.Log.Infow("requested key does not exist", "key", key)
				return api.ErrKeyNotFound
			}
			return err
		}
		if item.IsDeletedOrExpired() {
			s.Log.Infow("requested key has expired", "key", key)
			return api.ErrKeyNotFound
		}
		item.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})
		diff := item.ExpiresAt() - uint64(time.Now().Unix())
		if diff < uint64(60*60*24) {
			//if expiring in less than 24 hours, reset ttl for another 14 days
			s.Log.Infow("requested key will expire in less than 24 hours; resetting TTL", "key", key, "expiry", item.ExpiresAt(), "expires_in_s", diff)
			e := badger.NewEntry(k, res).WithTTL(time.Hour * time.Duration(api.DefaultTLL))
			err := txn.SetEntry(e)
			//Read shouldn't error here
			if err != nil {
				s.Log.Warnw("error updating TTL for requested temp key", "key", key, "err", err)
			}
		}

		return nil
	})

	return res, err
}

func (s *Store) Update(key string, data []byte, ctx context.Context) error {
	ttl := extractTTL(ctx)
	s.Log.Infow("received update request", "key", key, "ttl", ttl)
	return s.db.Update(func(txn *badger.Txn) error {
		k := []byte(key)
		//sanity check that uuid is not already used...
		_, err := txn.Get(k)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return api.ErrKeyNotFound
			}
			return err
		}
		if ttl > 0 {
			e := badger.NewEntry(k, data).WithTTL(time.Hour * time.Duration(ttl))
			return txn.SetEntry(e)
		}
		return txn.Set(k, data)
	})
}

func (s *Store) Delete(key string, ctx context.Context) error {
	s.Log.Infow("received delete request", "key", key)
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				s.Log.Debugw("deleting a key that does not exist", "key", key)
				return nil
			}
			return err
		}
		return nil
	})
}

func extractTTL(ctx context.Context) int {
	x := ctx.Value("ttl")
	//expecting ttl to be an integer value >= 0; if not int then default to
	//14 days; if ttl = 0 then assume to be permanent
	ttl := api.DefaultTLL
	if val, ok := x.(int); ok {
		if val >= 0 {
			ttl = val
		}
	}
	return ttl
}
