package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/genshinsim/gcsim/pkg/model"
)

// SubmissionStore is used to store both pre-compute submissions as well as pre-approval submissions.
// It's up to the store implementation to keep track of which tags has already approved/rejected a submission
// and purge any submission that has been fully approved/rejected for all tags
type SubmissionStore interface {
	Complete(context.Context, string, *model.SimulationResult) error
	Submit(context.Context, *model.Submission) (string, error)
}

type submitEntry struct {
	Config      string `json:"config"`
	Description string `json:"description"`
}

func (s *Server) submitEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserContextKey)
		s.Log.Infow("submission request received", "user", user)
		if user == nil {
			s.Log.Infow("submission request failed - user not logged in")
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}
		var o submitEntry
		err := json.NewDecoder(r.Body).Decode(&o)
		if err != nil {
			s.Log.Infow("submission request - bad body", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userStr, ok := user.(string)
		if !ok {
			//this should never happen
			s.Log.Infow("submission request - user context is not a string?", "user", user)
			http.Error(w, "unexpected server error", http.StatusInternalServerError)
			return
		}

		//TODO: check for invalid entries like blank config or blank description?

		id, err := s.cfg.SubmissionStore.Submit(r.Context(), &model.Submission{
			Config:      o.Config,
			Description: o.Description,
			Submitter:   userStr,
		})

		if err != nil {
			//shouldn't happen
			s.Log.Infow("submission request - add to queue failed", "err", err)
			http.Error(w, "unexpected serever error", http.StatusInternalServerError)
			return
		}

		s.Log.Infow("submission request - received succesfully", "id", id)

		//TODO: notify

		w.Write([]byte("ok"))
		w.WriteHeader(r.Response.StatusCode)
	}
}
