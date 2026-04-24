package nefer

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1BaseShadeScaleBonus = 0.6
	c4VerdantDewRateKey = "nefer-c4-verdant-dew-rate"
	c4ResShredKey = "nefer-c4-res-shred"
	c4NearbyRadius = 10
	c4ResShredLinger = 270
)

func c1ShadeScaleBonus(cons int) float64 {
	if cons < 1 {
		return 0
	}
	return c1BaseShadeScaleBonus
}

func (c *char) c1ShadeScaleBonus() float64 {
	return c1ShadeScaleBonus(c.Base.Cons)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.ascendantGleam {
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
				Base: modifier.NewBaseWithHitlag(c4ResShredKey, c4ResShredLinger),
				Ele: attributes.Dendro,
				Value: -0.2,
			})
		}
	}, c4ResShredKey)
}
