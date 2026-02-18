package lauma

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	lunarbloomBonusKey = "lauma-lunarbloom-bonus"
	a1Key              = "light-for-the-frosty-night"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	if c.ascendantGleam {
		c.a1Ascendant()
	} else {
		// we must have nascent because lauma contributes one moonsign herself
		c.a1Nascent()
	}
}

func (c *char) a1Nascent() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		if !c.StatusIsActive(a1Key) {
			return false
		}

		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*info.AttackEvent)

		switch ae.Info.AttackTag {
		case attacks.AttackTagBloom:
		case attacks.AttackTagHyperbloom:
		case attacks.AttackTagBurgeon:
		default:
			return false
		}

		// critrate stacks with nahida c2 while critdmg is overwritten
		ae.Snapshot.Stats[attributes.CR] += 0.15
		ae.Snapshot.Stats[attributes.CD] = 1

		c.Core.Log.NewEvent("lauma a1 buff", glog.LogCharacterEvent, ae.Info.ActorIndex).
			Write("final_crit", ae.Snapshot.Stats[attributes.CR])

		return false
	}, "lauma-a1-reaction-dmg-buff")
}

func (c *char) a1Ascendant() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		if !c.StatusIsActive(a1Key) {
			return false
		}

		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		ae := args[1].(*info.AttackEvent)

		switch ae.Info.AttackTag {
		case attacks.AttackTagDirectLunarBloom:
		default:
			return false
		}

		ae.Snapshot.Stats[attributes.CR] += 0.1
		ae.Snapshot.Stats[attributes.CD] += 0.2
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma a1 buff", glog.LogCharacterEvent, ae.Info.ActorIndex).
				Write("final_critrate", ae.Snapshot.Stats[attributes.CR]).
				Write("final_critdmg", ae.Snapshot.Stats[attributes.CD])
		}

		return false
	}, "lauma-a1-reaction-dmg-buff")
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	// increase skill dmg of self by EM * 0.4% up to 32%
	m := make([]float64, attributes.EndStatType)
	em := c.Stat(attributes.EM)
	m[attributes.DmgP] = min(0.004*em, 0.32)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("lauma-a4", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil
			}
			return m
		},
	})
}

func (c *char) lunarbloomInit() {
	c.Core.Flags.Custom[reactable.LunarBloomEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarBloom:
		default:
			return false
		}

		em := c.Stat(attributes.EM)
		bonus := min(em*0.000175, 0.14)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("lauma adding lunarbloom base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarbloomBonusKey)
}

func (c *char) a4SpiritEnvoyCooldownReduction() float64 {
	if c.Base.Ascension < 4 {
		return 1.0
	}

	em := c.Stat(attributes.EM)
	return 1.0 - min(0.2, em*0.0002)
}
