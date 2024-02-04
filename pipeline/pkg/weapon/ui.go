package weapon

import (
	"bytes"
	"encoding/json"
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
	s, err := protojson.Marshal(m)
	if err != nil {
		return err
	}
	dst := &bytes.Buffer{}
	err = json.Indent(dst, s, "", "  ")
	if err != nil {
		return err
	}
	os.Remove(path)
	err = os.WriteFile(path, dst.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}
