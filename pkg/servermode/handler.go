package servermode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/go-chi/chi"
)

func (s *Server) ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		running := s.isRunning(id)
		if running {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) running() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		running := s.isRunning(id)
		if running {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Write([]byte(errorRecover(r).Error()))
				w.WriteHeader(http.StatusBadRequest)
			}
		}()
		id := chi.URLParam(r, "id")
		s.Log.Info("request to run sample", "id", id)
		var payload struct {
			Config string `json:"config"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			s.Log.Info("body did not decode to json", "id", id, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		cfg, _, err := simulator.Parse(payload.Config)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		marshalled, err := json.Marshal(cfg)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(marshalled)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) sample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				w.Write([]byte(errorRecover(r).Error()))
				w.WriteHeader(http.StatusBadRequest)
			}
		}()
		id := chi.URLParam(r, "id")
		s.Log.Info("request to run sample", "id", id)
		var payload struct {
			Config string `json:"config"`
			Seed   uint64 `json:"seed"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			s.Log.Info("body did not decode to json", "id", id, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		opts := simulator.Options{
			GZIPResult:       false,
			ResultSaveToPath: "",
			ConfigPath:       "",
		}
		data, err := simulator.GenerateSampleWithSeed(payload.Config, payload.Seed, opts)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		marshalled, err := data.MarshalJSON()
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(marshalled)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) run() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		s.Log.Info("request to run sim", "id", id)
		var payload struct {
			Config string `json:"config"`
		}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			s.Log.Info("body did not decode to json", "id", id, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		// don't run if already running
		if s.isRunning(id) {
			s.Log.Info("run request failed; already running", "id", id)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("already running!!"))
			return
		}
		s.Log.Info("config decoded ok; running", "id", id)

		// create worker
		wk := &worker{
			id:     id,
			cfg:    payload.Config,
			log:    s.Log,
			cancel: make(chan bool),
		}
		s.Lock()
		s.pool[id] = wk
		s.Unlock()

		// start worker
		go wk.run(s.WorkerCount, s.FlushInterval)

		// add a timeout
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
			defer cancel()
			for {
				select {
				case <-wk.cancel:
					// someone cancelled so we're done
					return
				case <-ctx.Done():
					// context must have timed out
					close(wk.cancel)
					wk.done = true
					if wk.err == nil {
						wk.err = fmt.Errorf("execution timed out after %s", s.Timeout)
					}
					return
				}
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) latest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		s.Log.Info("request for latest", "id", id)

		wk, ok := s.pool[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var res struct {
			Result string `json:"result"`
			Hash   string `json:"hash"`
			Done   bool   `json:"done"`
			Error  string `json:"error"`
		}
		res.Done = wk.done

		// regardless of what results looks like, we should delete worker if done
		defer func() {
			if res.Done {
				s.Lock()
				delete(s.pool, id)
				s.Unlock()
			}
		}()

		s.Log.Info("found data", "id", id, "res.Error", res.Error, "res.Done", res.Done)
		if wk.err == nil {
			// if no error then we expect there to be some kind of result
			if wk.result == nil {
				s.Log.Info("unexpected result is nil", "id", id)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("unexpected result is blank"))
				return
			}
			b, err := wk.result.MarshalJSON()
			if err != nil {
				s.Log.Info("error marshalling result to json", "id", id, "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			res.Result = string(b)
			res.Hash, err = wk.result.Sign(s.ShareKey)
			if err != nil {
				s.Log.Info("error signing result", "id", id, "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			res.Error = wk.err.Error()
		}

		msg, err := json.Marshal(res)
		if err != nil {
			s.Log.Info("error marshalling final response to json", "id", id, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
	}
}

func (s *Server) cancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		wk, ok := s.pool[id]
		if !ok {
			s.Log.Info("cancel request received; worker does not exist", "id", id)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if wk.done {
			s.Log.Info("cancel request received; already done", "id", id)
			w.WriteHeader(http.StatusOK)
			return
		}
		s.Log.Info("cancelling run", "id", id)
		close(wk.cancel)
		wk.done = true
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type Progress struct {
	Curr int
	Max  int
	Err  error
}

func (s *Server) Progress() map[string]Progress {
	progress := make(map[string]Progress)
	s.Lock()
	for id, wk := range s.pool {
		prog := Progress{}
		if wk.err != nil {
			prog.Err = wk.err
			continue
		}
		if wk.result == nil {
			s.Log.Info("unexpected result is nil", "id", id)
		}

		prog.Max = int(wk.result.GetSimulatorSettings().Iterations)
		if wk.result.Statistics == nil {
			prog.Curr = 0
		} else {
			prog.Curr = int(wk.result.Statistics.Iterations)
		}

		progress[id] = prog
	}
	s.Unlock()
	return progress
}
