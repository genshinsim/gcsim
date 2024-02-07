package furina

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Furina, NewChar)
}

type Arkhe int

const (
	ousia Arkhe = iota
	pneuma
)

func (a Arkhe) String() string {
	switch a {
	case ousia:
		return "Ousia"
	case pneuma:
		return "Pneuma"
	}
	return "unknown"
}

type char struct {
	*tmpl.Character
	curFanfare          float64
	maxQFanfare         float64
	maxC2Fanfare        float64
	burstBuff           []float64
	a4Buff              []float64
	a4IntervalReduction float64
	lastSummonSrc       int
	arkhe               Arkhe
	c6Count             int
	c6HealSrc           int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3
	c.HasArkhe = true

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.arkhe = ousia

	c.a1()

	c.a4()
	c.a4Tick()

	c.burstInit()

	if c.Base.Cons >= 2 {
		c.c2()
	}

	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "ousia":
		return c.arkhe == ousia, nil
	case "fanfare":
		if c.StatusIsActive(burstKey) {
			return c.curFanfare, nil
		}
		return 0, nil
	case "c6-count":
		return c.c6Count, nil
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) NextQueueItemIsValid(next action.Eval) error {
	// can use charge without attack beforehand unlike most of the other sword users
	switch next.Action {
	case action.ActionCharge:
		return nil
	}
	return c.NextQueueItemIsValid(next)
}
