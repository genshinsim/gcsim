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
	DBPath string
}

type Store struct {
	cfg Config
	db  *badger.DB
	Log *zap.SugaredLogger
	UnimplementedResultStoreServer
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

func (s *Store) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	id := generateID()
	s.Log.Infow("received create request", "id", id, "ttl", req.Ttl)

	err := s.db.Update(func(txn *badger.Txn) error {
		key := []byte(id)

		//sanity check that uuid is not already used...
		_, err := txn.Get(key)
		switch err {
		case badger.ErrKeyNotFound:
		case nil:
			s.Log.Warnw("unexpected id collision", "id", id)
			return fmt.Errorf("unexpected key already exists: %v", id)
		default:
			return err
		}

		if req.Ttl > 0 {
			e := badger.NewEntry(key, req.Result).WithTTL(time.Hour * time.Duration(req.Ttl))
			return txn.SetEntry(e)
		}
		return txn.Set(key, req.Result)
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &CreateResponse{
		Key: id,
	}, nil
}

func (s *Store) Random(ctx context.Context, req *RandomRequest) (*RandomResponse, error) {
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
		//unexpected error trying to grab key
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.Log.Infow("random key selection done", "len", len(keys))
	if len(keys) == 0 {
		return nil, status.Error(codes.NotFound, "no key found - collection size is 0")
	}
	// Pick a random key from the list of keys
	n := keys[rand.Intn(len(keys))]

	return &RandomResponse{
		Key: string(n),
	}, nil
}

func (s *Store) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	var res []byte
	var expiry uint64
	s.Log.Infow("received request to view result", "key", req.Key)
	err := s.db.Update(func(txn *badger.Txn) error {
		k := []byte(req.Key)
		//sanity check that uuid is not already used...
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		if item.IsDeletedOrExpired() {
			return badger.ErrKeyNotFound
		}
		item.Value(func(val []byte) error {
			res = append([]byte{}, val...)
			return nil
		})
		expiry = item.ExpiresAt()
		diff := expiry - uint64(time.Now().Unix())
		s.Log.Infow("result retrieved ok", "key", req.Key, "expiry", expiry, "left", diff)
		if diff < uint64(60*60*24) {
			//if expiring in less than 24 hours, reset ttl for another 14 days
			s.Log.Infow("requested key will expire in less than 24 hours; resetting TTL", "key", req.Key, "expiry", item.ExpiresAt(), "expires_in_s", diff)
			e := badger.NewEntry(k, res).WithTTL(time.Hour * time.Duration(api.DefaultTLL))
			err := txn.SetEntry(e)
			//Read shouldn't error here
			if err != nil {
				s.Log.Errorw("error updating TTL for requested temp key", "key", req.Key, "err", err)
			}
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.Infow("requested key does not exist", "key", req.Key)
			return nil, status.Error(codes.NotFound, "key not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &ReadResponse{
		Key:    req.Key,
		Result: res,
		Ttl:    expiry,
	}, nil
}

func (s *Store) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	s.Log.Infow("received update request", "key", req.Key, "ttl", req.Ttl)
	err := s.db.Update(func(txn *badger.Txn) error {
		k := []byte(req.Key)
		//sanity check that uuid is not already used...
		_, err := txn.Get(k)
		if err != nil {
			return err
		}
		if req.Ttl > 0 {
			e := badger.NewEntry(k, req.Result).WithTTL(time.Hour * time.Duration(req.Ttl))
			return txn.SetEntry(e)
		}
		return txn.Set(k, req.Result)
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.Infow("requested key does not exist", "key", req.Key)
			return nil, status.Error(codes.NotFound, "key not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &UpdateResponse{Key: req.Key}, nil
}

func (s *Store) SetTTL(ctx context.Context, req *SetTTLRequest) (*SetTTLResponse, error) {
	s.Log.Infow("received request to set ttl", "key", req.Key, "ttl", req.Ttl)
	err := s.db.Update(func(txn *badger.Txn) error {
		k := []byte(req.Key)
		var data []byte
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			data = append([]byte{}, val...)
			return nil
		})
		if req.Ttl > 0 {
			e := badger.NewEntry(k, data).WithTTL(time.Hour * time.Duration(req.Ttl))
			return txn.SetEntry(e)
		}
		return txn.Set(k, data)
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.Infow("requested key does not exist", "key", req.Key)
			return nil, status.Error(codes.NotFound, "key not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &SetTTLResponse{Key: req.Key}, nil
}

func (s *Store) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	s.Log.Infow("received delete request", "key", req.Key)
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(req.Key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				s.Log.Debugw("deleting a key that does not exist", "key", req.Key)
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			s.Log.Infow("requested key does not exist", "key", req.Key)
			return nil, status.Error(codes.NotFound, "key not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &DeleteResponse{Key: req.Key}, nil
}
