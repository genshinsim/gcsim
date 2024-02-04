package testhelper

import (
	_ "embed"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/prototext"
)

//go:embed test_weapon_data.pb
var pbDataWeap []byte
var bweap *model.WeaponData

func init() {
	bweap = &model.WeaponData{}
	err := prototext.Unmarshal(pbDataWeap, bweap)
	if err != nil {
		panic(err)
	}
}

type Weapon struct {
	Index int
}

func (b *Weapon) SetIndex(idx int)        { b.Index = idx }
func (b *Weapon) Init() error             { return nil }
func (b *Weapon) Data() *model.WeaponData { return bweap }

func NewFakeWeapon(_ *core.Core, _ *character.CharWrapper, _ info.WeaponProfile) (info.Weapon, error) {
	return &Weapon{}, nil
}
