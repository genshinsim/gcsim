package neuvillette

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/sourcewaterdroplet"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Neuvillette, NewChar)
}

type char struct {
	*tmpl.Character
	lastSwap               int
	chargeJudgeStartF      int
	chargeJudgeDur         int
	tickAnimLength         int
	tickAnimLengthC6Extend int
	chargeEarlyCancelled   bool
	a1BaseStackCount       int
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
	c.HasArkhe = true

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

	c.onSwap()

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) getSourcewaterDroplets() []*sourcewaterdroplet.SourcewaterDroplet {
	player := c.Core.Combat.Player()

	// TODO: this is an approximation based on an ongoing KQM ticket (faster-neuvi-balls)
	// The fan is bigger than the 60 degrees in the ticket to account for the 10 degree camera tilt.
	segment := combat.NewCircleHitOnTargetFanAngle(player, nil, 14, 80)
	rect := combat.NewBoxHitOnTarget(player, geometry.Point{Y: -7}, 8, 18)

	droplets := make([]*sourcewaterdroplet.SourcewaterDroplet, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*sourcewaterdroplet.SourcewaterDroplet)
		if !ok {
			continue
		}
		if !droplet.IsWithinArea(rect) && !droplet.IsWithinArea(segment) {
			continue
		}
		droplets = append(droplets, droplet)
	}

	return droplets
}

func (c *char) getSourcewaterDropletsC6() []*sourcewaterdroplet.SourcewaterDroplet {
	player := c.Core.Combat.Player()

	circle := combat.NewCircleHitOnTarget(player, nil, 15)

	droplets := make([]*sourcewaterdroplet.SourcewaterDroplet, 0)
	for _, g := range c.Core.Combat.Gadgets() {
		droplet, ok := g.(*sourcewaterdroplet.SourcewaterDroplet)
		if !ok {
			continue
		}
		if !droplet.IsWithinArea(circle) {
			continue
		}
		droplets = append(droplets, droplet)
	}

	return droplets
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "droplets":
		return len(c.getSourcewaterDroplets()), nil
	case "droplets-c6":
		return len(c.getSourcewaterDropletsC6()), nil
	default:
		return c.Character.Condition(fields)
	}
}

// used for early CA cancel swap cd calculation
func (c *char) onSwap() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// do nothing if next char isn't neuvillette
		next := args[1].(int)
		if next != c.Index {
			return false
		}
		c.lastSwap = c.Core.F
		return false
	}, "neuvillette-swap")
}
