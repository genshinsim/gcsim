package ifa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1EnergyKey    = "ifa-c1-energy"
	c1EnergyICDKey = "ifa-c1-energy-icd"
)

func (c *char) c1CB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}
	if c.StatusIsActive(c1EnergyICDKey) {
		return
	}
	c.AddEnergy(c1EnergyKey, 6.0)
	c.AddStatus(c1EnergyICDKey, 8*60, true)
}

func (c *char) c2Mult() float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	if c.Base.Ascension < 1 {
		return 0
	}
	return 4
}

func (c *char) c2CapIncrease() float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	if c.Base.Ascension < 1 {
		return 0
	}
	return 50
}

func (c *char) c4OnBurst() {
	if c.Base.Cons < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 100

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("ifa-c2", 15*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return m
		},
	})
}

func (c *char) c6OnHoldAttackSkill() {
	if c.Base.Cons < 6 {
		return
	}

	if c.Core.Rand.Float64() > 0.5 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Tonic Shot C6",
		AttackTag:      attacks.AttackTagNormal,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagIfaSkill,
		ICDGroup:       attacks.ICDGroupIfaSkillHit,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           1.2,
	}

	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.PrimaryTarget(),
		nil,
		3,
	)

	c.QueueCharTask(func() {
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.Core.QueueAttack(
			ai,
			ap,
			0,
			0,
		)
	}, 1)
}
