package noelle

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Noelle, NewChar)
}

type char struct {
	*tmpl.Character
	shieldTimer int
	a4Counter   int
	burstBuff   []float64
	caState     ChargeState
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstBuff = make([]float64, attributes.EndStatType)
	c.a1()
	c.c2()
	return nil
}

// Noelle Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *info.AttackInfo) info.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.StatModIsActive(burstBuffKey) {
		// infusion to attacks only
		switch ai.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagPlunge:
		case attacks.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Geo
	}

	return ds
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 24
	case info.AnimationYelanN0StartDelay:
		return 18
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
