package character

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (g *Generator) DumpJSON(path string) error {
	// delete existing
	err := g.writeCharDataJSON(path + "/char_data.generated.json")
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) writeCharDataJSON(path string) error {
	data := make(map[string]*model.AvatarData)
	for _, v := range g.data {
		// hide promodata from ui json; not needed
		x := proto.Clone(v).(*model.AvatarData)
		x.Stats = nil
		// trim stat scaling too
		x.SkillDetails.SkillScaling = nil
		x.SkillDetails.AttackScaling = nil
		x.SkillDetails.BurstScaling = nil
		x.SkillDetails.A1 = 0
		x.SkillDetails.A4 = 0
		x.SkillDetails.A1Scaling = nil
		x.SkillDetails.A4Scaling = nil
		data[v.Key] = x
	}
	m := &model.AvatarDataMap{
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
	dst.WriteString("\n")
	os.Remove(path)
	err = os.WriteFile(path, dst.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}
