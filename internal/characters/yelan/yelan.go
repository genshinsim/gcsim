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
	a4buff  []float64
	c2icd   int
	c6count int
	c4count int //keep track of number of enemies tagged
	naHook  *common.NAHook
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

	barb, ok := p.Params["barb"]
	if ok && barb > 0 {
		c.SetTag(breakthroughStatus, 1)
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
	c.naHook = common.NewNAHook(c.CharWrapper, c.Core, "yelan burst", burstKey, 60, burstICDKey, common.Get0PercentN0Delay, c.summonExquisiteThrow)
	c.naHook.NAStateHook()
	return nil
}
