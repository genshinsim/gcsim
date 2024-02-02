package weapon

import (
	"os"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (g *Generator) DumpUIJSON(path string) error {
	// delete existing
	err := g.writeDataJSON(path + "/weapon_data.generated.json")
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) writeDataJSON(path string) error {
	data := make(map[string]*model.WeaponData)
	for _, v := range g.data {
		// hide promodata from ui json; not needed
		x := proto.Clone(v).(*model.WeaponData)
		x.BaseStats = nil
		data[v.Key] = x
	}
	m := &model.WeaponDataMap{
		Data: data,
	}
	s := protojson.Format(m)
	os.Remove(path)
	err := os.WriteFile(path, []byte(s), 0o644)
	if err != nil {
		return err
	}

	return nil
}
