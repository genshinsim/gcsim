package jahoda

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Jahoda, NewChar)
}

type char struct {
	*tmpl.Character
	flaskAbsorbCheckLocation info.AttackPattern
	flaskAbsorb              attributes.Element
	flaskGauge               int
	flaskGaugeMax            int
	pursuitDuration          int
	skillSrc                 int
	skillTravel              int
	burstAbsorbCheckLocation info.AttackPattern
	burstSrc                 int
	a1HighestEle             attributes.Element
	robotAi                  info.AttackInfo
	robotHi                  info.HealInfo
	robotCount               int
	robotHitmarkInterval     float64
	c2NextHighestEle         attributes.Element
	a4Buff                   []float64
	c6Buff                   []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.Moonsign = 1
	c.HasArkhe = false

	c.flaskAbsorb = attributes.NoElement
	c.flaskGauge = 0
	c.flaskGaugeMax = 100

	c.a1HighestEle = attributes.NoElement

	c.c2NextHighestEle = attributes.NoElement

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1Init()

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
