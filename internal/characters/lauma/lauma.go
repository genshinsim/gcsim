package lauma

import (
	"github.com/gammazero/deque"

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
	skillSrc             int
	moonSong             int
	c6SkillPaleHymnCount int
	paleHymnStacks       deque.Deque[int]
	c6PaleHymnStacks     deque.Deque[int]
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.deerStateReady = true

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

	c.a1Init()
	c.a4Init()

	c.c1()
	c.c2()
	c.c6Elevation()

	c.initBurst()

	c.chargeInit()

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

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.deerStateReady {
		return 0
	}
	return c.Character.ActionStam(a, p)
}
