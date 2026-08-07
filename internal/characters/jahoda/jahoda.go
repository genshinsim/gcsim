package jahoda

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	absorbPriority []attributes.Element

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
	robotAi                  info.AttackInfo
	robotHealCoeff           float64
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

	//c.absorbPriority = make([]attributes.Element, 4)
	c.absorbPriority = append(c.absorbPriority, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)

	c.flaskAbsorb = attributes.NoElement
	c.a1HighestEle = attributes.NoElement
	c.c2NextHighestEle = attributes.NoElement

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Init()

	c.robotHealCoeff = 1.0
	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.EM] = 100

	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 13
	}
	if k == info.AnimationYelanN0StartDelay {
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
