package jahoda

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

var elePriority = []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo}

type char struct {
	*tmpl.Character

	flaskAbsorbCheckLocation info.AttackPattern
	flaskAbsorb              attributes.Element
	flaskGauge               int
	pursuitDuration          int
	skillSrc                 int
	meowballSrc              int
	meowballTravel           int

	burstAbsorbCheckLocation info.AttackPattern
	burstSrc                 int
	a1HighestEle             attributes.Element
	a1DmgMult                float64
	a1HealMult               float64
	robotCount               int
	robotHitmarkInterval     float64

	c2NextHighestEle attributes.Element
	a4Buff           []float64
	c6Buff           []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.Moonsign = 1

	c.flaskAbsorb = attributes.NoElement
	c.a1HighestEle = attributes.NoElement
	c.c2NextHighestEle = attributes.NoElement

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1HealMult = 1.0
	c.a1Init()
	c.a4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 13
	case info.AnimationYelanN0StartDelay:
		return 11
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) enemyAuraInArea(area info.AttackPattern, priority []attributes.Element) attributes.Element {
	for _, ele := range priority {
		for _, enemy := range c.Core.Combat.Enemies() {
			if enemy == nil || !enemy.IsAlive() {
				continue
			}

			target, ok := enemy.(info.TargetWithAura)
			if !ok {
				continue
			}

			collision, _ := target.AttackWillLand(area)
			if !collision {
				continue
			}

			if target.AuraContains(ele) {
				return ele
			}
		}
	}

	return attributes.NoElement
}
