package collei

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Collei, NewChar)
}

type char struct {
	*tmpl.Character
	burstExtendCount int
	a1StartFrame     int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.burstExtendCount = 0
	c.a1StartFrame = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	return nil
}
