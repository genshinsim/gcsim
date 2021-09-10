package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/genshinsim/gsim"
	"github.com/genshinsim/gsim/pkg/core"
)

type runConfig struct {
	Options core.RunOpt `json:"options"`
	Config  string      `json:"config"`
}

func (s *Server) handleRun(ctx context.Context, r wsRequest) {

	s.Log.Debugw("handleRun: request to run received")

	var cfg runConfig
	err := json.Unmarshal([]byte(r.Payload), &cfg)

	if err != nil {
		s.Log.Debugw("handleRun: invalid request payload", "payload", r.Payload)
		handleErr(r, http.StatusBadRequest, "bad request payload")
		return
	}

	if cfg.Options.Debug {
		s.Log.Debugw("handleRun: running with debug")
		s.runDebug(cfg, r)
	} else {
		s.Log.Debugw("handleRun: running without debug")
		s.run(cfg, r)
	}

}

func (s *Server) runDebug(cfg runConfig, r wsRequest) {

	now := time.Now()
	logfile := fmt.Sprintf("./%v.txt", now.Format("2006-01-02-15-04-05"))

	cfg.Options.DebugPaths = []string{logfile}

	result, err := gsim.Run(cfg.Config, cfg.Options)
	if err != nil {
		handleErr(r, http.StatusBadRequest, err.Error())
		return
	}

	file, err := os.Open(logfile)
	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()
	var log strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log.WriteString(scanner.Text())
		log.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}

	err = file.Close()

	if err != nil {
		handleErr(r, http.StatusInternalServerError, err.Error())
		return
	}

	s.Log.Debugw("run complete", "result", result)

	result.Text = result.PrettyPrint()
	result.Debug = log.String()

	data, _ := json.Marshal(result)
	e := wsResponse{
		ID:      r.ID,
		Status:  http.StatusOK,
		Payload: string(data),
	}
	msg, _ := json.Marshal(e)
	r.client.send <- msg

	os.Remove(logfile)
}

func (s *Server) run(cfg runConfig, r wsRequest) {

	result, err := gsim.Run(cfg.Config, cfg.Options)
	if err != nil {
		handleErr(r, http.StatusBadRequest, err.Error())
		return
	}

	result.Text = result.PrettyPrint()

	s.Log.Debugw("run complete", "result", result)

	data, _ := json.Marshal(result)
	e := wsResponse{
		ID:      r.ID,
		Status:  http.StatusOK,
		Payload: string(data),
	}
	msg, _ := json.Marshal(e)
	r.client.send <- msg
}
