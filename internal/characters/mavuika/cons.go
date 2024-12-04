package mavuika

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	c.nightsoulState.MaxPoints = 120
	c.fightingSpiritMult = 1.25
}

func (c *char) c1Atk() {
	if c.Base.Cons < 1 {
		return
	}
	buff := make([]float64, attributes.EndStatType)
	buff[attributes.ATKP] = 0.4
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("mavuika-c1-atk", 8*60),
		Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			return buff, true
		},
	})
}

func (c *char) c2BaseIncrease() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.BaseATK] = 200
	c.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag("mavuika-c2-base-atk", -1),
		Amount: func() ([]float64, bool) {
			if c.nightsoulState.HasBlessing() {
				return m, true
			}
			return nil, false
		},
	})
}

func (c *char) c2FlatIncrease(tag attacks.AttackTag) float64 {
	if c.Base.Cons < 2 {
		return 0
	}
	switch tag {
	case attacks.AttackTagNormal:
		return 0.6 * c.TotalAtk()
	case attacks.AttackTagExtra:
		return 0.9 * c.TotalAtk()
	case attacks.AttackTagElementalBurst:
		return 1.2 * c.TotalAtk()
	default:
		return 0
	}
}

func (c *char) c2AddDefMod() {
	if c.Base.Cons < 2 {
		return
	}
	// TODO: the actual detection range
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7), nil)
	for _, t := range enemies {
		// if you set duration=-1, it won't work
		t.AddDefMod(combat.DefMod{
			Base:  modifier.NewBaseWithHitlag("mavuika-c2-def-shred", 6000),
			Value: -0.2,
		})
	}
}

func (c *char) c2DeleteDefMod() {
	if c.Base.Cons < 2 {
		return
	}
	// all enemies
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 100), nil)
	for _, t := range enemies {
		if t.DefModIsActive("mavuika-c2-def-shred") {
			t.DeleteDefMod("mavuika-c2-def-shred")
		}
	}
}

func (c *char) c6RSRModeHit() {
	if c.Base.Cons < 6 {
		return
	}
	if c.flamestriderModeActive {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Flamestrider Hit (C6)",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     0,
		Mult:           2.,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6)
	// TODO: the actual frames
	c.Core.QueueAttack(ai, ap, 40, 40)
}

func (c *char) c6FlamestriderModeHit(src int) func() {
	return func() {
		if c.Base.Cons < 6 {
			return
		}
		if src != c.nightsoulSrc {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		if c.flamestriderModeActive {
			ai := combat.AttackInfo{
				ActorIndex:     c.Index,
				Abil:           "Rings of Searing Radiance DMG (C6)",
				AttackTag:      attacks.AttackTagElementalArt,
				AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
				ICDTag:         attacks.ICDTagNone,
				ICDGroup:       attacks.ICDGroupDefault,
				StrikeType:     attacks.StrikeTypeDefault,
				Element:        attributes.Pyro,
				Durability:     0,
				Mult:           4.,
			}
			// TODO: change hurt box
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5), 0, 0)
		}
		c.QueueCharTask(c.ringsOfSearchingRadianceHit(src), 3*60)
	}
}
