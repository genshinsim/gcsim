package chevreuse

import (
	_ "embed"
	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/proto"
)

//go:embed data_gen.pb
var pbData []byte
var base *model.AvatarData

func init() {
	base = &model.AvatarData{}
	err := proto.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
}

func (x *char) Data() *model.AvatarData {
	return base
}
