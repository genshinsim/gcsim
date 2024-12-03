package chasca

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Faruzan, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState *nightsoul.State
	nightsoulSrc   int
	partyTypes     []attributes.Element
	phecCount      int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	w.Character = &c

	c.partyTypes = nil

	return nil
}

func (c *char) Init() error {
	types := map[attributes.Element]bool{}
	for _, other := range c.Core.Player.Chars() {
		switch ele := other.Base.Element; ele {
		case attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro:
			c.phecCount++
			types[ele] = true
		}
	}
	for ele := range types {
		c.partyTypes = append(c.partyTypes, ele)
	}

	c.a4()
	if c.Base.Cons >= 6 {
		c.c6Collapse()
	}

	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}
