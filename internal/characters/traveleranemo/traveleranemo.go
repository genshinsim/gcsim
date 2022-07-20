package traveleranemo

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
	core.RegisterCharFunc(keys.TravelerAnemo, NewChar)
}

type char struct {
	*tmpl.Character
	qInfuse             attributes.Element
	qICDTag             combat.ICDTag
	eInfuse             attributes.Element
	eICDTag             combat.ICDTag
	infuseCheckLocation combat.AttackPattern
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum
	c.infuseCheckLocation = combat.NewDefCircHit(0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}
