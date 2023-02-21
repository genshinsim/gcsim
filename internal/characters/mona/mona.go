package mona

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

const (
	bubbleKey = "mona-bubble"
	omenKey   = "omen-debuff"
)

func init() {
	core.RegisterCharFunc(keys.Mona, NewChar)
}

type char struct {
	*tmpl.Character
	c2icd    int
	c6Src    int
	c6Stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstHook()
	c.burstDamageBonus()
	c.a4()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}
