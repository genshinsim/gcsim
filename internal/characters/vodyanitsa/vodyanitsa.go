package vodyanitsa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character

	skillSrc int

	nearbyWanderingVortexExp int
	leadVocal                int
	chorus                   int

	c1Buff   []float64
	c2Buff   []float64
	c4Buff   []float64
	c6Buff   []float64
	c4Stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Init()
	c.a4Init()

	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "lead-vocal":
		if !c.StatusIsActive(a4Key) {
			return 0, nil
		}
		return c.leadVocal, nil
	case "chorus":
		if !c.StatusIsActive(a4Key) {
			return 0, nil
		}
		return c.chorus, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 5
	case info.AnimationYelanN0StartDelay:
		return 5
	}
	return c.Character.AnimationStartDelay(k)
}
