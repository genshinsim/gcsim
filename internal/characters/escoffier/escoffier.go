package escoffier

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Escoffier, NewChar)
}

type char struct {
	*tmpl.Character
	skillSrc         int
	skillTravel      int
	a1Src            int
	a4HydroCryoCount int
	c1Active         bool
	c1Buff           []float64
	c2Count          int
	c4Count          int
	c6Count          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.SkillCon = 3
	c.BurstCon = 5

	c.EnergyMax = burstEnergy[c.TalentLvlBurst()]
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4Init()
	c.c1Init()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 11
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "c2-count":
		return c.c2Count, nil
	case "c4-count":
		return c.c4Count, nil
	case "c6-count":
		return c.c6Count, nil
	default:
		return c.Character.Condition(fields)
	}
}
