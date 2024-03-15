package gaming

import (
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2Key = "gaming-c2"
const c4Key = "gaming-c4"
const c4IcdKey = "gaming-c4"

// When the Suanni Man Chai from Suanni's Gilded Dance meets back up with Gaming,
// it will heal 15% of Gaming's HP.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Index,
		Message: "Bringer of Blessing (C1)",
		Type:    player.HealTypePercent,
		Src:     0.15,
		Bonus:   c.Stat(attributes.Heal),
	})
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		hi := args[0].(*player.HealInfo)
		overheal := args[3].(float64)

		if overheal <= 0 {
			return false
		}

		if hi.Target != c.Index {
			return false
		}

		c.AddStatus(c2Key, 5*60, true)

		return false
	}, c2Key+"-onheal")

	c2buff := make([]float64, attributes.EndStatType)
	c2buff[attributes.ATKP] = 0.2
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2Key+"-atk", 5*60),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			if !c.StatusIsActive(c2Key) {
				return nil, false
			}
			return c2buff, true
		},
	})
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4IcdKey) {
		return
	}
	c.AddStatus(c4IcdKey, 0.2*60, true)
	c.AddEnergy(c4Key, 2)
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c6Buff := make([]float64, attributes.EndStatType)
	c6Buff[attributes.CR] = 0.2
	c6Buff[attributes.CD] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("gaming-c6", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagPlunge {
				return nil, false
			}
			if !strings.Contains(atk.Info.Abil, ePlungeKey) {
				return nil, false
			}
			return c6Buff, true
		},
	})
}
