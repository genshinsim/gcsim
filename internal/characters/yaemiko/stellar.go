package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarConductText = " (Stellar-Conduct)"
	revelationKey      = "yae-revelation"
	revelationICDKey   = "yae-revelation-icd"
)

// Extends the duration of the Sesshou Sakura by 10s. When a nearby party member triggers a
// Superconduct or Stellar-Conduct reaction, the next instance of Sesshou Sakura lightning is
// enhanced as follows: DMG dealt is increased at 80% of Yae Miko's ATK. This effect can trigger
// once every 2.5s.
//
// Radiance: Stellar-Conduct: An enhanced Sesshou Sakura lightning bolt hit on an opponent will also
// cause an additional instance of Electro DMG at 200% of Yae Miko's ATK. This DMG is considered
// Stellar-Conduct DMG.
//
// Additionally, when inside a Polestar Field, Yae Miko will enter the Radiance: Stellar-Conduct state.
func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(revelationKey, 8*60, true)
	}

	c.Core.Events.Subscribe(event.OnSuperconduct, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, revelationKey)
}

func (c *char) revelationBonusSkillDur() int {
	if !c.revelation {
		return 0
	}

	return 10 * 60
}

func (c *char) revelationEnhanceDMG() (float64, info.AttackCBFunc) {
	if !c.revelation {
		return 0, nil
	}

	if !c.StatusIsActive(revelationKey) {
		return 0, nil
	}

	if c.StatusIsActive(revelationICDKey) {
		return 0, nil
	}

	c.AddStatus(revelationICDKey, 2.5*60, true)

	c.DeleteStatus(revelationKey)

	dmg := 0.8 * c.TotalAtk()
	if !c.isRadianceSSC() {
		return dmg, nil
	}

	done := false
	cb := func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}

		if done {
			return
		}

		done = true

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Sesshou Sakura" + stellarConductText,
			AttackTag:        attacks.AttackTagDirectStellarConduct,
			ICDTag:           attacks.ICDTagNone,
			ICDGroup:         attacks.ICDGroupDefault,
			StrikeType:       attacks.StrikeTypePierce,
			Element:          attributes.Electro,
			Mult:             2,
			IgnoreDefPercent: 1,
		}

		ap := combat.NewSingleTargetHit(a.Target.Key())

		c.Core.QueueAttack(ai, ap, 12, 12)
	}

	return dmg, cb
}

func (c *char) isRadianceSSC() bool {
	if !c.revelation {
		return false
	}
	return c.StatusIsActive(reactable.PolestarFieldKey)
}
