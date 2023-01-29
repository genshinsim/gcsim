package yelan

import (
	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

const (
	breakthroughStatus = "yelan_breakthrough"
	c6Status           = "yelan_c6"
	burstKey           = "yelanburst"
	burstICDKey        = "yelanburstICD"
)

func init() {
	core.RegisterCharFunc(keys.Yelan, NewChar)
}

type char struct {
	*tmpl.Character
	a4buff       []float64
	breakthrough bool // tracks breakthrough state
	c2icd        int
	c6count      int
	c4count      int // keep track of number of enemies tagged
}

func NewChar(s *core.Core, w *character.CharWrapper, p profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.c2icd = 0
	c.c6count = 0

	breakthrough, ok := p.Params["breakthrough"]
	if ok && breakthrough > 0 {
		c.breakthrough = true
	}

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a1()
	(&common.NAHook{
		C:           c.CharWrapper,
		AbilName:    "yelan burst",
		Core:        c.Core,
		AbilKey:     burstKey,
		AbilProcICD: 60,
		AbilICDKey:  burstICDKey,
		DelayFunc:   common.Get0PercentN0Delay,
		SummonFunc:  c.summonExquisiteThrow,
	}).Enable()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "breakthrough":
		return c.breakthrough, nil
	default:
		return c.Character.Condition(fields)
	}
}
