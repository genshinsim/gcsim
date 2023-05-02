package baizhu

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Baizhu, NewChar)
}

type char struct {
	*tmpl.Character
	skillAtk *combat.AttackEvent
	c6done   bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.BurstCon = 3
	c.SkillCon = 5
	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}
