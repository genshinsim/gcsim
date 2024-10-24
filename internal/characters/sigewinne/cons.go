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
			if atk.Info.AttackTag != attacks.AttackTagElementalArt ||
				atk.Info.AttackTag != attacks.AttackTagElementalArtHold ||
				atk.Info.AttackTag != attacks.AttackTagElementalBurst {
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

func (c *char) addC2Shield() {
	shieldHP := 2.5 * c.MaxHP()

	c.Core.Player.Shields.Add(c.newShield(shieldHP))
	c.Tags["shielded"] = 1
}

func (c *char) removeC2Shield() {
	c.Tags["shielded"] = 0
	c.Tags["a1"] = 0
}

func (c *char) newShield(base float64) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.ActorIndex = c.Index
	n.Tmpl.Target = -1
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = shield.SigewinneC2
	n.Tmpl.Ele = attributes.Hydro
	n.Tmpl.HP = base
	n.Tmpl.Name = "Sigewinee C2"
	n.Tmpl.Expires = -1
	n.c = c
	return n
}

type shd struct {
	*shield.Tmpl
	c *char
}

func (s *shd) OnExpire() {
	s.c.removeC2Shield()
}

func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok {
		s.c.removeC2Shield()
	}
	return taken, ok
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
			crit_amt := make([]float64, attributes.EndStatType)
			crit_amt[attributes.CD] = min(C6CDmgCap, c.MaxHP()*C6CDmgHpRatio)
			crit_amt[attributes.CR] = min(C6CRateCap, c.MaxHP()*C6CRateCap)
			return crit_amt, true
		},
	})
}
