package dahlia

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Dahlia, NewChar)
}

type char struct {
	*tmpl.Character
	maxBenisonStacks int
	burstAttackSpeed float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.maxBenisonStacks = 4

	c.burstAttackSpeed = c.MaxHP() * 0.5 // unlocked by A4 only so maybe not here
	if c.Base.Cons >= 6 {
		c.burstAttackSpeed += 10
	}

	c.setupBurst()
	c.a1()
	c.a4()
	c.c6()

	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 10
	case info.AnimationYelanN0StartDelay:
		return 0
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
