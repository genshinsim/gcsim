package neuvillette

import (
	"math"

	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Neuvillette, NewChar)
}

type char struct {
	*tmpl.Character
	lastThorn              int
	lastSkillParticle      int
	lastc4                 int
	chargeJudgeStartF      int
	chargeJudgeDur         int
	tickAnimLength         int
	tickAnimLengthC6Extend int
	chargeEarlyCancelled   bool
	a1Statuses             []NeuvA1Keys
	a4Buff                 []float64
	chargeAi               combat.AttackInfo
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	c.lastThorn = math.MinInt / 2
	c.lastc4 = math.MinInt / 2
	c.lastSkillParticle = math.MinInt / 2

	c.chargeEarlyCancelled = false
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Ascension >= 1 {
		c.a1()
	}

	if c.Base.Ascension >= 4 {
		c.a4Buff = make([]float64, attributes.EndStatType)
		c.a4()
	}

	if c.Base.Cons >= 1 {
		c.c1()
	}

	if c.Base.Cons >= 2 {
		c.c2()
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) getSourcewaterDroplets() []*common.SourcewaterDroplet {
	playerPos := c.Core.Combat.Player().Pos()
	droplets := make([]*common.SourcewaterDroplet, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*common.SourcewaterDroplet)
		if !ok {
			continue
		}
		if droplet.Pos().Distance(playerPos) <= 15 {
			droplets = append(droplets, droplet)
		}
	}
	return droplets
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "droplets":
		return len(c.getSourcewaterDroplets()), nil
	default:
		return c.Character.Condition(fields)
	}
}
