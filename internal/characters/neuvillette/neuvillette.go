package neuvillette

import (
	"math"

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
	lastThorn         int
	lastSkillParticle int
	lastc6            int

	a1Statuses []NeuvA1Keys
	a4Buff     []float64
	chargeAi   combat.AttackInfo
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	c.lastThorn = math.MinInt / 2
	c.lastc6 = math.MinInt / 2
	c.lastSkillParticle = math.MinInt / 2
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4()

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
