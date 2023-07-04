package character

import (
	"os"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (g *Generator) DumpJSON(path string) error {
	//delete existing
	err := g.writeCharDataJSON(path + "/char_data.generated.json")
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) writeCharDataJSON(path string) error {
	data := make(map[string]*model.AvatarData)
	for _, v := range g.data {
		//hide promodata from ui json; not needed
		x := proto.Clone(v).(*model.AvatarData)
		x.Stats = nil
		data[v.Key] = x
	}
	m := &model.AvatarDataMap{
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
