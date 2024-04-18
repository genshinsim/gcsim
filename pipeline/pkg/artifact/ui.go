package artifact

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
	err := g.writeCharDataJSON(path + "/artifact_data.generated.json")
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) writeCharDataJSON(path string) error {
	data := make(map[string]*model.ArtifactData)
	for k, v := range g.data {
		x := proto.Clone(v).(*model.ArtifactData)
		data[k] = x
		data[k].ImageNames = nil
	}
	m := &model.ArtifactDataMap{
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
