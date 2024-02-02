package collei

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

var dendroEvents = []event.Event{event.OnBurning, event.OnQuicken, event.OnAggravate, event.OnSpread, event.OnBloom, event.OnHyperbloom, event.OnBurgeon}

func init() {
	core.RegisterCharFunc(keys.Collei, NewChar)
}

type char struct {
	*tmpl.Character
	burstPos           geometry.Point
	burstExtendCount   int
	sproutShouldExtend bool
	sproutShouldProc   bool
	sproutSrc          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.sproutShouldProc = false
	c.sproutShouldExtend = false

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 2 && c.Base.Ascension >= 1 {
		c.c2()
	}
	return nil
}
