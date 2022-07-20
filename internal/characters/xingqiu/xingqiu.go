package xingqiu

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Xingqiu, NewChar)
}

type char struct {
	*tmpl.Character
	numSwords     int
	nextRegen     bool
	burstCounter  int
	burstTickSrc  int
	orbitalActive bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.burstStateHook()
	return nil
}
