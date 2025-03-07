package collei

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	sproutKey        = "collei-a1"
	sproutHitmark    = 86
	sproutTickPeriod = 89
	a4Key            = "collei-a4-modcheck"
)

// If one of your party members has triggered Burning, Quicken, Aggravate, Spread, Bloom, Hyperbloom, or Burgeon reactions
// before the Floral Ring returns, it will grant the character the Sprout effect upon return, which will continuously deal
// Dendro DMG equivalent to 40% of Collei's ATK to nearby opponents for 3s.
// If another Sprout effect is triggered during its initial duration, the initial effect will be removed.
// DMG dealt by Sprout is considered Elemental Skill DMG.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	f := func(...interface{}) bool {
		if c.sproutShouldProc {
			return false
		}
		if !c.StatusIsActive(skillKey) {
			return false
		}
		c.sproutShouldProc = true
		c.Core.Log.NewEvent("collei a1 proc", glog.LogCharacterEvent, c.Index)
		return false
	}

	for _, evt := range dendroEvents {
		switch evt {
		case event.OnHyperbloom, event.OnBurgeon:
			c.Core.Events.Subscribe(evt, f, "collei-a1")
		default:
			c.Core.Events.Subscribe(evt, func(args ...interface{}) bool {
				if _, ok := args[0].(*enemy.Enemy); !ok {
					return false
				}
				return f(args...)
			}, "collei-a1")
		}
	}
}

func (c *char) a1AttackInfo() combat.AttackInfo {
	return combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Floral Sidewinder (A1)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagColleiSprout,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       0.4,
	}
}

// When a character within the Cuilein-Anbar Zone triggers Burning, Quicken, Aggravate, Spread, Bloom, Hyperbloom, or Burgeon reactions,
// the Zone's duration will be increased by 1s.
// A single Trump-Card Kitty can be extended this way by up to 3s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	f := func(args ...interface{}) bool {
		if !c.StatusIsActive(burstKey) {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !char.StatusIsActive(a4Key) {
			return false
		}
		if c.burstExtendCount >= 3 {
			return false
		}
		c.ExtendStatus(burstKey, 60)
		c.burstExtendCount++
		c.Core.Log.NewEvent("collei a4 proc", glog.LogCharacterEvent, c.Index).
			Write("extend_count", c.burstExtendCount)
		return false
	}

	for _, evt := range dendroEvents {
		switch evt {
		case event.OnHyperbloom, event.OnBurgeon:
			c.Core.Events.Subscribe(evt, f, "collei-a4")
		default:
			c.Core.Events.Subscribe(evt, func(args ...interface{}) bool {
				if _, ok := args[0].(*enemy.Enemy); !ok {
					return false
				}
				return f(args...)
			}, "collei-a4")
		}
	}
}

func (c *char) a1Ticks(startFrame int, snap combat.Snapshot) {
	if !c.StatusIsActive(sproutKey) {
		return
	}
	if startFrame != c.sproutSrc {
		c.Core.Log.NewEvent("collei a1 tick ignored, src diff", glog.LogCharacterEvent, c.Index).
			Write("src", startFrame).
			Write("new src", c.sproutSrc)
		return
	}
	c.Core.QueueAttackWithSnap(
		c.a1AttackInfo(),
		snap,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2),
		0,
	)
	c.Core.Tasks.Add(func() {
		c.a1Ticks(startFrame, snap)
	}, sproutTickPeriod)
}
