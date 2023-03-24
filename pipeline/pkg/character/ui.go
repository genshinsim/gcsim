package character

import (
	"os"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
)

func (g *Generator) DumpUIJSON(path string) error {
	//delete existing
	err := g.writeCharDataJSON(path + "/char_data.generated.json")
	if err != nil {
		return err
	}
	err = g.writeCharKeyMapJSON(path + "/char_key_to_id_map.generated.json")
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) writeCharDataJSON(path string) error {
	m := &model.AvatarDataMap{
		Data: g.data,
	}
	d, err := protojson.Marshal(m)
	if err != nil {
		return err
	}
	os.Remove(path)
	err = os.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) writeCharKeyMapJSON(path string) error {
	data := make(map[string]int64)
	for k, v := range g.data {
		data[v.Key] = k
	}
	m := &model.AvatarKeyMap{
		Data: data,
	}

	d, err := protojson.Marshal(m)
	if err != nil {
		return err
	}
	os.Remove(path)
	err = os.WriteFile(path, d, 0644)
	if err != nil {
		return err
	}
	return nil
}
