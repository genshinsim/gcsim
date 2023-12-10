package freminet

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key      = "freminet-c1"
	c2Key      = "freminet-c2"
	c4Key      = "freminet-c4"
	c6Key      = "freminet-c6"
	c4c6IcdKey = "freminet-c4-c6-icd"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.CR] = 0.15

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c1Key, -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if !strings.HasPrefix(atk.Info.Abil, pressureBaseName) {
				return nil, false
			}
			return buff, true
		},
	})
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	amt := 2.0
	if c.skillStacks == 4 {
		amt = 3
	}

	c.AddEnergy(c1Key, amt)
}

func (c *char) c4c6() {
	if c.Base.Cons < 4 {
		return
	}

	c4M := make([]float64, attributes.EndStatType)
	c6M := make([]float64, attributes.EndStatType)

	c4c6Buff := func(args ...interface{}) bool {
		if _, ok := args[0].(*gadget.Gadget); ok {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if c.StatusIsActive(c4c6IcdKey) {
			return false
		}

		c.AddStatus(c4c6IcdKey, 18, true)

		if !c.StatModIsActive(c4Key) {
			c.c4Stacks = 0
		}

		c.c4Stacks++

		if c.c4Stacks > 2 {
			c.c4Stacks = 2
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c4Key, 6*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				c4M[attributes.ATKP] = float64(c.c4Stacks) * 0.09
				return c4M, true
			},
		})

		c.Core.Log.NewEvent("freminet c4 proc", glog.LogCharacterEvent, c.Index)

		if c.Base.Cons < 6 {
			return false
		}

		// C6

		if !c.StatModIsActive(c6Key) {
			c.c6Stacks = 0
		}

		c.c6Stacks++

		if c.c6Stacks > 3 {
			c.c6Stacks = 3
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c6Key, 6*60),
			AffectedStat: attributes.CD,
			Amount: func() ([]float64, bool) {
				c6M[attributes.CD] = float64(c.c6Stacks) * 0.12
				return c6M, true
			},
		})

		c.Core.Log.NewEvent("freminet c6 proc", glog.LogCharacterEvent, c.Index)

		return false
	}

	c.Core.Events.Subscribe(event.OnShatter, c4c6Buff, "freminet-c4-c6-shatter")
	c.Core.Events.Subscribe(event.OnFrozen, c4c6Buff, "freminet-c4-c6-frozen")
	c.Core.Events.Subscribe(event.OnSuperconduct, c4c6Buff, "freminet-c4-c6-superconduct")
}
