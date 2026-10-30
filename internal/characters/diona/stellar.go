package diona

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
	revelationSkillKey = "diona-revelation-skill"
	revelationKey      = "diona-revelation"
	revelationICDKey   = "diona-revelation-icd"
)

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *char) getRadiance() radianceState {
	if !c.revelation {
		return radianceNone
	}

	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(reactable.StellarSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

// Within 20s after using the Elemental Skill Icy Paws, Diona will also fire off
// three Icy Paws when party members trigger Superconduct, Stellar-Conduct, Cryo
// Swirl, or Stellar Swirl reactions. Icy Paws fired this way cannot generate
// Elemental Particles, nor will they generate a shield upon hitting opponents.
// This effect can occur once every 3.5s.
func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !c.StatusIsActive(revelationSkillKey) {
			return
		}

		if c.StatusIsActive(revelationICDKey) {
			return
		}

		c.AddStatus(revelationICDKey, 3.5*60, true)

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Icy Paw (Revelation)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       paw[c.TalentLvlSkill()],
		}

		for i := range 3 {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					0.5,
				),
				0,
				5+i,
			)
		}
	}

	c.Core.Events.Subscribe(event.OnStellarSwirl, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnSwirlCryo, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnSuperconduct, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, revelationKey)
}

func (c *char) revelationOnSkill() {
	if !c.revelation {
		return
	}
	c.AddStatus(revelationSkillKey, 20*60, true)
}
