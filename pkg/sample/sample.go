package sample

import (
	"compress/zlib"
	"encoding/json"
	"os"

	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

type Sample struct {
	Config           string                       `json:"config"`
	InitialCharacter string                       `json:"initial_character"`
	CharacterDetails []simulation.CharacterDetail `json:"character_details"`
	TargetDetails    []enemy.EnemyProfile         `json:"target_details"`
	Seed             string                       `json:"seed"`
	Logs             []map[string]interface{}     `json:"logs"`
}

// TODO: this should probably just be a utility function for any json serializable struct
func (s *Sample) Save(fpath string, gz bool) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if gz {
		f, err := os.OpenFile(fpath+".gz", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		defer f.Close()
		zw := zlib.NewWriter(f)
		zw.Write(data)
		return zw.Close()
	}

	return os.WriteFile(fpath, data, 0644)
}
