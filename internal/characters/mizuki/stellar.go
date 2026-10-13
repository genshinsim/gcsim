package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	stellarSwirlText       = " (Stellar Swirl)"
	dreamDrifterBuffKey    = "dreamdrifter-buff"
	dreamDrifterBuffICDKey = "dreamdrifter-buff-icd"
	revelationShareEMKey   = "mizuki-revelation-share-em"
)

func (c *char) isRadianceSSw() bool {
	return c.StatusIsActive(reactable.RadianceSwirlKey)
}

func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	onSkillHit := func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagElementalArt {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if !c.StatusIsActive(dreamDrifterBuffKey) {
			return
		}
		c.QueueCharTask(func() {
			c.DeleteStatus(dreamDrifterBuffKey)
		}, 0.25*60)

		atk.Info.FlatDmg += c.Stat(attributes.EM) * 10

		if !c.isRadianceSSw() {
			return
		}

		t := args[0].(info.Target)

		c.stellarHit(t)
	}

	onSwirl := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if c.StatusIsActive(dreamDrifterBuffICDKey) {
			return
		}
		c.AddStatus(dreamDrifterBuffICDKey, 2.5*60, true)
		c.QueueCharTask(func() {
			c.AddStatus(dreamDrifterBuffKey, -1, false)
			c.Core.Log.NewEvent("Dream Drifter state enhanced", glog.LogCharacterEvent, c.Index())
		}, 0.25*60)
	}

	c.Core.Events.Subscribe(event.OnSwirlCryo, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlPyro, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlHydro, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnSwirlElectro, onSwirl, "mizuki-swirl")
	c.Core.Events.Subscribe(event.OnStellarSwirl, onSwirl, "mizuki-ssw")

	c.Core.Events.Subscribe(event.OnEnemyHit, onSkillHit, "mizuki-on-skill-hit")

	c.revelationBuff = make([]float64, attributes.EndStatType)
	c.revelationBuff[attributes.EM] = c.Stat(attributes.EM) * 0.1
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:  modifier.NewBase(revelationShareEMKey, -1),
			Extra: true,
			Amount: func() []float64 {
				if c.StatusIsActive(dreamDrifterStateKey) {
					return c.revelationBuff
				}
				return nil
			},
		})
	}
}

func (c *char) revelationShareEM(src int) {
	if !c.revelation {
		return
	}

	if c.cloudSrc != src {
		return
	}

	if !c.StatusIsActive(dreamDrifterStateKey) {
		return
	}

	c.QueueCharTask(func() {
		c.revelationShareEM(src)
	}, 0.5*60)
}

func (c *char) stellarHit(t info.Target) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Dreamdrifter" + stellarSwirlText,
		AttackTag:        attacks.AttackTagDirectStellarSwirl,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		PoiseDMG:         skillActivatePoise,
		Element:          attributes.Anemo,
		Durability:       0,
		UseEM:            true,
		Mult:             10,
		IgnoreDefPercent: 1,
	}

	ap := combat.NewSingleTargetHit(t.Key())

	c.Core.QueueAttack(ai, ap, 0.2*60, 0.2*60) // datamine says this should be 12 frames while practice says the difference in damage number rendering is 20 frames ¯\_(ツ)_/¯
}
