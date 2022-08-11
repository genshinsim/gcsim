package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/genshinsim/gcsim/pkg/result"
)

type Simulation struct {
	Metadata   string `json:"metadata"`
	ViewerFile string `json:"viewer_file"`
}

func (s *Simulation) DecodeViewer() (*result.Summary, error) {
	//base64 zlib encoded string

	return nil, nil
}

type SimStore interface {
	Fetch(url string) (Simulation, error)
}

type PostgRESTStore struct {
	URL string
}

func (b *PostgRESTStore) Fetch(key string) (Simulation, error) {
	url := fmt.Sprintf(`%v/simulations?simulation_key=eq.%v`, b.URL, key)
	r, err := http.Get(url)
	if err != nil {
		return Simulation{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return Simulation{}, err
		} else {
			return Simulation{}, fmt.Errorf("bad status code %v msg %v", r.StatusCode, string(msg))
		}
	}

	var result []Simulation

	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return Simulation{}, err
	}

	if len(result) == 0 {
		return Simulation{}, fmt.Errorf("unexpected result length is 0")
	}

	return result[0], nil
}
