package albedo

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Albedo, NewChar)
}

type char struct {
	*tmpl.Character
	lastConstruct   int
	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot
	bloomSnapshot   combat.Snapshot
	icdSkill        int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum

	c.icdSkill = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillHook()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons == 6 {
		c.c6()
	}
	return nil
}
