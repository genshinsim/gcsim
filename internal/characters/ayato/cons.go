package ayato

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.4
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("ayato-c1", -1),
			Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if a.Info.AttackTag != attacks.AttackTagNormal || x.HP()/x.MaxHP() > 0.5 {
					return nil, false
				}
				return m, true
			},
		})
	}
}

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = 0.5
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("ayato-c2", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			if c.stacks >= 3 {
				return m, true
			}
			return nil, false
		},
	})
}

// After using Kamisato Art: Kyouka, Ayato's next Shunsuiken attack will create
// 2 extra Shunsuiken strikes when they hit opponents, each one dealing 450% of Ayato's ATK as DMG.
// Both these Shunsuiken attacks will not be affected by Namisen.
func (c *char) makeC6CB() combat.AttackCBFunc {
	if c.Base.Cons < 6 || !c.c6Ready {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.Core.Player.Active() != c.Index {
			return
		}
		if !c.c6Ready {
			return
		}
		c.c6Ready = false

		c.Core.Log.NewEvent("ayato c6 proc'd", glog.LogCharacterEvent, c.Index)
		ai := combat.AttackInfo{
			Abil:               "Ayato C6",
			ActorIndex:         c.Index,
			AttackTag:          attacks.AttackTagNormal,
			ICDTag:             combat.ICDTagNormalAttack,
			ICDGroup:           combat.ICDGroupDefault,
			StrikeType:         combat.StrikeTypeSlash,
			Element:            attributes.Hydro,
			Durability:         25,
			Mult:               4.5,
			HitlagFactor:       0.01,
			HitlagHaltFrames:   0.03 * 60,
			CanBeDefenseHalted: false,
			IsDeployable:       true,
		}
		for i := 0; i < 2; i++ {
			c.Core.QueueAttack(
				ai,
				combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 8, 7),
				20+i*2,
				20+i*2,
			)
		}
	}
}
