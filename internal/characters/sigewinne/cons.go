package sigewinne

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	C1flatconvalescenceIncrease = 100
	C1flatconvalescenceCap      = 3500

	C6CDmgHpRatio  = 0.022 / 1000
	C6CRateHpRatio = 0.004 / 1000
	C6CDmgCap      = 1.1
	C6CRateCap     = 0.2
)

func (c *char) c2() {
	c2func := func() func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			t, ok := args[0].(*enemy.Enemy)
			if !ok {
				return false
			}
			if atk.Info.ActorIndex != c.Index {
				return false
			}
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return false
			}

			t.AddResistMod(combat.ResistMod{
				Base:  modifier.NewBaseWithHitlag("sigewinne-c2-hydro-res-shred", 8*60),
				Ele:   attributes.Hydro,
				Value: -0.35,
			})
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "Sigewinne C2 proc").Write("char", c.Index).Write("target", t.Key())

			return false
		}
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, c2func(), "sigewinne-c2")
}

func (c *char) addC2Shield(duration int) func() {
	return func() {
		shieldHP := 0.3 * c.MaxHP()
		c.Core.Player.Shields.Add(&shield.Tmpl{
			ActorIndex: c.Index,
			Target:     c.Index,
			Src:        c.Core.F,
			Name:       "Sigewinne C2 shield",
			ShieldType: shield.SigewinneC2,
			HP:         shieldHP,
			Ele:        attributes.Hydro,
			Expires:    duration,
		})
	}
}

func (c *char) c6() {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		c.c6CritMode()
		return false
	}, "sigewinne-c6-activation")
}

func (c *char) c6CritMode() {
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sigewinne-c6-crit-buff", 15*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			critAmt := make([]float64, attributes.EndStatType)
			critAmt[attributes.CD] = min(C6CDmgCap, c.MaxHP()*C6CDmgHpRatio)
			critAmt[attributes.CR] = min(C6CRateCap, c.MaxHP()*C6CRateHpRatio)
			return critAmt, true
		},
	})
}
