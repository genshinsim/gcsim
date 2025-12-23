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
	flaskAbsorbDuration      int
	flaskGauge               int
	flaskGaugeMax            int
	skillSrc                 int
	skillTravel              int
	burstAbsorbCheckLocation info.AttackPattern
	burstSrc                 int
	a1HighestEle             attributes.Element
	robotAi                  info.AttackInfo
	robotHi                  info.HealInfo
	robotCount               int
	robotInterval            float64
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

	c.flaskAbsorb = attributes.NoElement
	c.flaskAbsorbDuration = -1
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

	c.c6()

	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 10 // Frames needed
	}
	if k == info.AnimationYelanN0StartDelay {
		return 10 // Frames needed
	}
	return c.Character.AnimationStartDelay(k)
}
