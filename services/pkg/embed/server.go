package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

//Server contains an authentication server
type Server struct {
	Router *chi.Mux
	Log    *zap.SugaredLogger
	cfg    Config
	work   chan string //key to generate embed for
}

type Config struct {
	PostgRESTPort int
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
	url := fmt.Sprintf(`http://localhost:%v/simulations?simulation_key=eq.%v`, s.cfg.PostgRESTPort, key)
	r, err := http.Get(url)
	if err != nil {
		s.Log.Warnw("error getting simulation", "err", err)
		done <- true
		return
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			s.Log.Warnw("error getting simulation", "status", r.StatusCode, "err_reading_body", err, "url", url)
		} else {
			s.Log.Warnw("error getting simulation", "status", r.StatusCode, "msg", string(msg), "url", url)
		}
		done <- true
		return
	}

	var result []simulation

	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		s.Log.Warnw("error parsing simulation", "err", err)
		done <- true
		return
	}

	if len(result) == 0 {
		s.Log.Warn("unexpected result length is 0")
		done <- true
		return
	}

	//pass metadata off to python script
	data := result[0].Metadata
	data = strings.TrimPrefix(data, `"`)
	data = strings.TrimSuffix(data, `"`)
	ioutil.WriteFile("test.json", []byte(data), fs.ModePerm)
	cmd := exec.Command("./embed.py", fmt.Sprintf("./images/%v.png", key))
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
		//if key not found return 404
		filepath := fmt.Sprintf("./images/%v.png", key)
		if stat, err := os.Stat(filepath); err != nil || stat.IsDir() {
			//it's possible file exist but we can't read it
			//we should try to regenerate it
			s.work <- key
			//in which case we return 404 anyways
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := ioutil.ReadFile(filepath)
		if err != nil {
			s.Log.Warn("error loading image", "key", key, "path", filepath, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
