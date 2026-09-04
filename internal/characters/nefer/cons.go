package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4VerdantDewRateKey = "nefer-c4-verdant-dew-rate"
	c4ResShredKey       = "nefer-c4-res-shred"
	c4NearbyRadius      = 10
	c4ResShredLinger    = 270
)

func (c *char) c1ShadeScaleBonus() float64 {
	if c.Base.Cons < 1 {
		return 0
	}
	return 0.6
}

func (c *char) c2VeilStackInit() {
	c.maxVeilStacks = 3
	if c.Base.Cons < 2 {
		return
	}
	c.maxVeilStacks = 5
}

func (c *char) c2VeilDurationBonus() int {
	if c.Base.Cons < 2 {
		return 0
	}
	return 5 * 60
}

func (c *char) c2VeilStackFull(count int) bool {
	if c.Base.Cons < 2 {
		return false
	}

	if count >= c.maxVeilStacks {
		return true
	}

	return false
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Player.AddVerdantDewRateMod(c4VerdantDewRateKey, -1, func() (float64, bool) {
		if c.Core.Player.Active() != c.Index() || !c.StatusIsActive(shadowDanceKey) {
			return 0, false
		}
		return 0.25, false
	})

	c.Core.Events.Subscribe(event.OnTick, func(args ...any) {
		if c.Core.Player.Active() != c.Index() || !c.StatusIsActive(shadowDanceKey) {
			return
		}

		playerPos := c.Core.Combat.Player().Pos()
		for _, target := range c.Core.Combat.Enemies() {
			enemy, ok := target.(info.Enemy)
			if !ok || !target.IsAlive() {
				continue
			}
			radius := 0.0
			if circle, ok := target.Shape().(*info.Circle); ok {
				radius = circle.Radius()
			}
			if playerPos.Distance(target.Pos()) > c4NearbyRadius+radius {
				continue
			}
			enemy.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag(c4ResShredKey, c4ResShredLinger),
				Ele:   attributes.Dendro,
				Value: -0.2,
			})
		}
	}, c4ResShredKey)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Core.Player.GetMoonsignLevel() < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagDirectLunarBloom {
			return
		}
		atk.Info.Elevation += 0.15
	}, "nefer-c6-elevation")
}

func (c *char) c6Charge(ai *info.AttackInfo) {
	if c.Base.Cons < 6 {
		return
	}

	c6EndHit := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Phantasm Performance (Shade C6)",
		AttackTag:        attacks.AttackTagDirectLunarBloom,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Dendro,
		UseEM:            true,
		IgnoreDefPercent: 1,
		HitlagFactor:     0.01,
		Mult:             c6PhantasmEndEM,
		FlatDmg:          c.Stat(attributes.EM) * c.c1ShadeScaleBonus(),
	}

	hook := func() {
		c6EndHit.Mult *= c.phantasmVeilBonus()
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 3.15}, 4.5)
		c.Core.QueueAttack(c6EndHit, ap, 1, 1)
	}

	c.Core.Tasks.Add(hook, phantasmAnimationLength)

	// change 2nd nefer hit to lunarbloom
	ai.AttackTag = attacks.AttackTagDirectLunarBloom
	ai.Durability = 0
	ai.UseEM = true
	ai.IgnoreDefPercent = 1
	ai.Mult = c6PhantasmHit2EM
	ai.FlatDmg = c.Stat(attributes.EM) * c.c1ShadeScaleBonus()
}
