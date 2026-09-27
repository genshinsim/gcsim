package kachina

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	playercharacter "github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1EnergyKey = "kachina-c1-energy"
	c4DefKey    = "kachina-c4-def"
	c6IcdKey    = "kachina-c6-icd"
)

func (c *char) initC1() {
	if c.Base.Cons < 1 {
		return
	}

	restore := func() {
		if !c.StatusIsActive(twirlyKey) || c.StatusIsActive(c1EnergyKey) {
			return
		}
		c.AddStatus(c1EnergyKey, 5*60, true)
		c.AddEnergy(c1EnergyKey, 3)
	}
	crystallizeHook := func(_ ...any) {
		restore()
	}
	c.Core.Events.Subscribe(event.OnCrystallizeCryo, crystallizeHook, "kachina-c1-crystallize-cryo")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, crystallizeHook, "kachina-c1-crystallize-electro")
	c.Core.Events.Subscribe(event.OnCrystallizeHydro, crystallizeHook, "kachina-c1-crystallize-hydro")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, crystallizeHook, "kachina-c1-crystallize-pyro")
	c.Core.Events.Subscribe(event.OnLunarCrystallize, crystallizeHook, "kachina-c1-lunar-crystallize")
}

func (c *char) initC4() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		if len(args) < 2 {
			return
		}
		if prev, ok := args[0].(int); ok {
			c.Core.Player.ByIndex(prev).DeleteStatMod(c4DefKey)
		}
		if next, ok := args[1].(int); ok {
			c.applyC4(next)
		}
	}, "kachina-c4-swap")
}

func (c *char) applyC4(index int) {
	if c.Base.Cons < 4 || !c.StatusIsActive(fieldKey) {
		return
	}
	if c.Core.Combat.Player().Pos().Distance(c.fieldCenter) > 5.2 {
		return
	}
	enemies := len(c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.fieldCenter, nil, 5.2), nil))
	bonus := []float64{0, 0.08, 0.12, 0.16, 0.20}[min(enemies, 4)]
	mod := make([]float64, attributes.EndStatType)
	mod[attributes.DEFP] = bonus
	c.Core.Player.ByIndex(index).AddStatMod(playercharacter.StatMod{
		Base:         modifier.NewBaseWithHitlag(c4DefKey, -1),
		AffectedStat: attributes.DEFP,
		Amount: func() []float64 {
			if !c.StatusIsActive(fieldKey) {
				return nil
			}
			return mod
		},
	})
}

func (c *char) initC6() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Events.Subscribe(event.OnShielded, func(args ...any) {
		if len(args) == 0 {
			return
		}
		sh, ok := args[0].(shield.Shield)
		if !ok {
			return
		}
		key := c6ShieldKey(sh)
		if c.c6ShieldKeys[key] {
			c.triggerC6(sh.ShieldTarget())
		}
		c.c6ShieldKeys[key] = true
	}, "kachina-c6-shielded")
	c.Core.Events.Subscribe(event.OnShieldBreak, func(args ...any) {
		if len(args) == 0 {
			return
		}
		sh, ok := args[0].(shield.Shield)
		if !ok {
			return
		}
		delete(c.c6ShieldKeys, c6ShieldKey(sh))
		c.triggerC6(sh.ShieldTarget())
	}, "kachina-c6-shield-break")
}

func c6ShieldKey(sh shield.Shield) string {
	return fmt.Sprintf("%d:%d", sh.Type(), sh.ShieldTarget())
}

func (c *char) triggerC6(target int) {
	if c.StatusIsActive(c6IcdKey) {
		return
	}
	active := c.Core.Player.Active()
	if target != -1 && target != active {
		return
	}
	c.AddStatus(c6IcdKey, 5*60, true)
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "This Time, I've Gotta Win",
		AttackTag:      attacks.AttackTagExtra,
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       100,
		Element:        attributes.Geo,
		Durability:     25,
		Mult:           2.0,
		UseDef:         true,
		IgnoreInfusion: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0)
}
