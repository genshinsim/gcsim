package embed

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/genshinsim/gcsim/services/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

//Server contains an authentication server
type Server struct {
	Router *chi.Mux
	db     *badger.DB
	Log    *zap.SugaredLogger
	Store  store.SimStore
	cfg    Config
	work   chan string //key to generate embed for
}

type Config struct {
	AssetFolder string
	DataFolder  string
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

type simulation struct {
	Metadata string `json:"metadata"`
}

func (s *Server) worker(key string, done chan bool) {
	//grab simulation meta data from postgrest server
	//http://localhost:3000/simulations\?simulation_key\=eq.57baa25d-27e6-4ffd-a152-71dd86272acf
	res, err := s.Store.Fetch(key)
	if err != nil {
		s.Log.Warnw("error getting simulation", "err", err)
		done <- true
		return
	}
	meta := res.Metadata
	meta = strings.TrimPrefix(meta, `"`)
	meta = strings.TrimSuffix(meta, `"`)

	//pass metadata off to python script
	filepath := fmt.Sprintf("./%v.png", key)
	os.Remove(filepath)
	ioutil.WriteFile("test.json", []byte(meta), fs.ModePerm)
	cmd := exec.Command("./embed.py", filepath)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// cmd.Stdin = strings.NewReader(result[0].Metadata)
	// out, err := cmd.Output()
	err = cmd.Run()
	if err != nil {
		// s.Log.Warnw("error generating image", "err", err, "out", string(out))
		s.Log.Warnw("error generating image", "err", err)
	}
	log.Println(stdBuffer.String())

	//try reading file and storing in db
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		s.Log.Warn("error loading image", "key", key, "path", filepath, "err", err)
		done <- true
		return
	}
	err = s.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), b).WithTTL(time.Hour * 24 * 60)
		err := txn.SetEntry(e)
		return err
	})
	if err != nil {
		s.Log.Warn("error saving image to db", "key", key, "path", filepath, "err", err)
	}

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
			item, err := txn.Get([]byte("answer"))
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
