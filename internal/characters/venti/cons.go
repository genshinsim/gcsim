package venti

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1HexICDKey       = "venti-c1-hex-icd"
	c2HexSkillBuffKey = "venti-c2-hex-skill-buff"
)

// C1:
// Fires 2 additional arrows per Aimed Shot, each dealing 33% of the original arrow's DMG.
func (c *char) c1(ai info.AttackInfo, hitmark, travel int) {
	if c.Base.Cons < 1 {
		return
	}

	ai.Abil += " (C1)"
	switch ai.AttackTag {
	case attacks.AttackTagExtra:
		ai.Mult /= 3.0
	case attacks.AttackTagNormal:
		if !c.IsHexerei {
			return
		}
		if c.StatusIsActive(c1HexICDKey) {
			return
		}

		c.AddStatus(c1HexICDKey, 0.25*60, true)

		ai.Mult /= 5.0
	default:
		return
	}
	ai.Mult /= 3.0
	for range 2 {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				info.Point{Y: -0.5},
				0.1,
				1,
			),
			hitmark,
			hitmark+travel,
		)
	}
}

// C2:
// Skyward Sonnet decreases opponents' Anemo RES and Physical RES by 12% for 10s.
// Opponents launched by Skyward Sonnet suffer an additional 12% Anemo RES and Physical RES decrease while airborne.
// TODO: the airborne part isn't implemented
func (c *char) c2(a info.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	m := -0.12

	if c.IsHexerei {
		m = -0.24
	}

	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-anemo", 600),
		Ele:   attributes.Anemo,
		Value: m,
	})
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-phys", 600),
		Ele:   attributes.Physical,
		Value: m,
	})
}

func (c *char) c2SkillBuffInit() {
	if c.Base.Cons < 2 {
		return
	}
	if !c.IsHexerei {
		return
	}
	c.AddStatus(c2HexSkillBuffKey, 15*60, true)
	c.ResetActionCooldown(action.ActionSkill)
}

func (c *char) c2OnSkill(ai info.AttackInfo) {
	if c.Base.Cons < 2 {
		return
	}
	if !c.IsHexerei {
		return
	}
	if !c.StatusIsActive(c2HexSkillBuffKey) {
		return
	}
	c.DeleteStatus(c2HexSkillBuffKey)
	ai.FlatDmg += ai.Mult * 2 * c.TotalAtk()
}

// C4:
// When Venti picks up an Elemental Orb or Particle, he receives a 25% Anemo DMG Bonus for 10s.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.c4bonus = make([]float64, attributes.EndStatType)
	c.c4bonus[attributes.AnemoP] = 0.25
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...any) {
		// only trigger if Venti catches the particle
		if c.Core.Player.Active() != c.Index() {
			return
		}
		if c.IsHexerei {
			return
		}
		// apply C4 to Venti
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
			AffectedStat: attributes.AnemoP,
			Amount: func() []float64 {
				return c.c4bonus
			},
		})
	}, "venti-c4")
}

func (c *char) c4Hexerei() {
	if c.Base.Cons < 4 {
		return
	}
	if !c.IsHexerei {
		return
	}
	if c.StatModIsActive("venti-c4-hexerei") {
		return
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("venti-c4-hexerei", 10*60),
		AffectedStat: attributes.AnemoP,
		Amount: func() []float64 {
			return c.c4bonus
		},
	})

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("venti-c4-hexerei", 10*60),
			Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
				if atk.Info.ActorIndex != c.Core.Player.Active() {
					return nil
				}

				return c.c4bonus
			},
		})
	}
}

// C6:
// Targets who take DMG from Wind's Grand Ode have their Anemo RES decreased by 20%.
// If an Elemental Absorption occurred, then their RES towards the corresponding Element is also decreased by 20%.
func (c *char) c6(ele attributes.Element) func(a info.AttackCB) {
	return func(a info.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("venti-c6-"+ele.String(), 600),
			Ele:   ele,
			Value: -0.20,
		})
	}
}

func (c *char) c6HexereiInit() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.IsHexerei {
		return
	}

	c.c6HexBonus = make([]float64, attributes.EndStatType)
	c.c6HexBonus[attributes.CD] = 1
	for _, ele := range []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro, attributes.Anemo} {
		c.c6ResistModRange = append(c.c6ResistModRange, "venti-c6-"+ele.String())
	}

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("venti-c6-hexerei-cd", -1),
		Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil
			}

			isBuffed := false
			for _, resMod := range c.c6ResistModRange {
				if e.ResistModIsActive(resMod) {
					isBuffed = true
				}
			}

			if isBuffed {
				return c.c6HexBonus
			}
			return nil
		},
	})
}
