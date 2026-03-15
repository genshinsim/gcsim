package nefer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Nefer, NewChar)
}

const (
	shadowDanceKey   = "nefer-shadow-dance"
	seedWindowKey    = "nefer-seed-window"
	veilEMBuffKey    = "nefer-veil-em-buff"
	seedAbsorbRadius = 6
)

type char struct {
	*tmpl.Character
	ascendantGleam bool
	veilstacks     int
	maxVeilStacks  int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.Moonsign = 1
	c.SetNumCharges(action.ActionSkill, 2)
	c.maxVeilStacks = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.ascendantGleam = c.Core.Player.GetMoonsignLevel() >= 2
	if c.Base.Cons >= 2 {
		c.maxVeilStacks = 5
	}
	c.lunarbloomInit()
	c.p1Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationYelanN0StartDelay:
		return 10
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) addVeilStacks(count int) {
	if count <= 0 {
		return
	}
	prev := c.veilstacks
	c.veilstacks = min(c.veilstacks+count, c.maxVeilStacks)
	c.applyVeilThresholdBuff(prev, c.veilstacks)
}

func (c *char) consumeVeilStacks() int {
	stacks := c.veilstacks
	c.veilstacks = 0
	return stacks
}

func (c *char) applyVeilThresholdBuff(prev, next int) {
	if next <= prev {
		return
	}

	amount := 0.0
	if c.Base.Cons >= 2 && next >= 5 && prev < 5 {
		amount = 200
	} else if next >= 3 && prev < 3 {
		amount = 100
	}
	if amount <= 0 {
		return
	}

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.EM] = amount
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(veilEMBuffKey, 8*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return buff
		},
	})
}
