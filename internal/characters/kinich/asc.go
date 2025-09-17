package kinich

import (
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	desolationKey = "desolation"
	a1Icd         = "a1-icd"
	a4StackKey    = "hunters-experience"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	hook := func(args ...any) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagBurningDamage:
		case attacks.AttackTagBurgeon:
		default:
			return false
		}
		if !t.StatusIsActive(desolationKey) {
			return false
		}
		if c.StatusIsActive(a1Icd) {
			return false
		}
		c.nightsoulState.GeneratePoints(7)
		c.AddStatus(a1Icd, 0.8*60, false)
		return false
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, hook, "kinich-a1")
}

func (c *char) a1CB(a info.AttackCB) {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.StatusIsActive(nightsoul.NightsoulBlessingStatus) {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// TODO: add the modifier for a gadget
	e.AddStatus(desolationKey, 12*60, true)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...any) bool {
		stacks := min(c.Tag(a4StackKey)+1, 2)
		c.AddStatus(a4StackKey, 15*60, true)
		c.SetTag(a4StackKey, stacks)
		return false
	}, "kinich-a4")
}

func (c *char) a4Amount() float64 {
	if c.Base.Ascension < 4 {
		return 0.0
	}
	stacks := c.Tag(a4StackKey)
	c.SetTag(a4StackKey, 0)
	c.DeleteStatus(a4StackKey)
	return 3.2 * float64(stacks) * c.TotalAtk()
}
