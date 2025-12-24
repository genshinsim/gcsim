package lauma

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Lauma, NewChar)
}

type char struct {
	*tmpl.Character
	ascendantGleam       bool
	deerStateReady       bool
	c6SkillPaleHymnCount int
	paleHymnStacks       []paleHymnStack
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.deerStateReady = true
	c.c6SkillPaleHymnCount = 8

	w.Character = &c
	w.Moonsign = 1

	return nil
}

func (c *char) Init() error {
	c.lunarbloomInit()

	chars := c.Core.Player.Chars()
	count := 0
	for _, ch := range chars {
		count += ch.Moonsign
	}
	if count >= 2 {
		c.ascendantGleam = true
	} else {
		c.ascendantGleam = false
	}

	c.a1()
	c.a4()

	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 6 {
		c.c6Elevation()
	}

	c.setupPaleHymnBuff()

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 10
	case info.AnimationYelanN0StartDelay:
		return 10
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
