package varesa

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/stacks"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Status = "rainbow-crash"

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(a1Status, 5*60, true)
}

func (c *char) a1PlungeBuff() float64 {
	if c.Base.Ascension < 1 {
		return 0.0
	}
	if !c.StatusIsActive(a1Status) {
		return 0.0
	}
	if c.Base.Cons >= 1 || c.nightsoulState.HasBlessing() {
		return 1.8
	}
	return 0.5
}

func (c *char) a1Cancel(_ combat.AttackCB) {
	if c.Base.Ascension < 1 {
		return
	}
	c.DeleteAttackMod(a1Status)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBase("varesa-a4", -1),
		Amount: func() ([]float64, bool) {
			m[attributes.ATKP] = 0.35 * float64(c.a4Stacks.Count())
			return m, true
		},
	})

	c.a4Stacks = stacks.NewMultipleRefreshNoRemove(2, c.QueueCharTask, &c.Core.F)
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.a4Stacks.Add(12 * 60)
		return false
	}, "varesa-a4")
}
