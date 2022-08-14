package embed

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/genshinsim/gcsim/services/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

//Server contains an authentication server
type Server struct {
	Router    *chi.Mux
	db        *badger.DB
	Log       *zap.SugaredLogger
	Store     store.SimStore
	Generator ImageGenerator
	cfg       Config
	work      chan string //key to generate embed for
}

type Config struct {
	DataFolder string
}

type ImageGenerator interface {
	Generate(sim store.Simulation, outpath string) error
}

func New(cfg Config, cust ...func(*Server) error) (*Server, error) {

	s := &Server{
		cfg: cfg,
	}

	//read config
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	db, err := badger.Open(badger.DefaultOptions(cfg.DataFolder))
	if err != nil {
		return nil, err
	}
	s.db = db

	s.Router = chi.NewRouter()
	s.routes()

	s.work = make(chan string)
	go s.pool()

	return s, nil
}

func (s *Server) pool() {
	var queue []string
	done := make(chan bool)

	count := 0

	for {
		select {
		case w := <-s.work:
			queue = append(queue, w)
		case <-done:
			count--
		}

		s.Log.Infow("worker status", "depth", len(queue), "count", count, "queue", queue)
		if count == 0 && len(queue) > 0 {
			//send out worker
			go s.worker(queue[0], done)
			count++
			queue = queue[1:]
		}
	}
}

func (s *Server) worker(key string, done chan bool) {
	res, err := s.Store.Fetch(key)
	if err != nil {
		s.Log.Warnw("error getting simulation", "err", err)
		done <- true
		return
	}

	filepath := fmt.Sprintf("./%v.png", key)
	os.Remove(filepath)

	err = s.Generator.Generate(res, filepath)
	if err != nil {
		s.Log.Warnw("error generating image", "key", key, "path", filepath, "err", err)
		done <- true
		return
	}

	//try reading file and storing in db
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		s.Log.Warnw("error loading image", "key", key, "path", filepath, "err", err)
		done <- true
		return
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), b).WithTTL(time.Hour * 24 * 60)
		err := txn.SetEntry(e)
		return err
	})

	if err != nil {
		s.Log.Warnw("error saving image to db", "key", key, "path", filepath, "err", err)
	}

	s.Log.Infow("generating embed ok", "key", key)

	//clean up after ourselves
	os.Remove(filepath)

	done <- true
}

func (s *Server) routes() {
	r := s.Router
	r.Use(middleware.Logger)

	r.Post("/embed/{key}", s.handleGenerateEmbed())
	r.Get("/embed/{key}", s.handleRetrieveEmbed())
}

func (s *Server) handleGenerateEmbed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		s.Log.Infof("we got work: %v", key)
		s.work <- key
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleRetrieveEmbed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		var data []byte

		err := s.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return err
			}
			err = item.Value(func(val []byte) error {
				data = append([]byte{}, val...)
				return nil
			})
			return err
		})

		if err != nil {
			if err == badger.ErrKeyNotFound {
				s.Log.Warnw("error getting image - not found", "key", key)
				//we should try to regenerate it
				s.work <- key
				//in which case we return 404 anyways
				w.WriteHeader(http.StatusNotFound)
				return
			}
			s.Log.Warnw("unexpected error reading image", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
