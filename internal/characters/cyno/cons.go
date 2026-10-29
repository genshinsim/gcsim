package cyno

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
	c1Key       = "cyno-c1"
	c1RevKey    = "cyno-together-we-rise"
	c2Key       = "cyno-c2"
	c2ICD       = c2Key + "-icd"
	c2RevICDKey = c2Key + "-revelation-icd"
	c4RevICDKey = "cyno-c4-revelation-icd"
	c6Key       = "cyno-c6"
	c6ICDKey    = c6Key + "-icd"
)

// After using Sacred Rite: Wolf's Swiftness, Cyno's Normal Attack SPD will be increased by 20% for 10s.
// If the Judication effect of his Passive Talent Featherfall Judgment is triggered during Secret Rite: Chasmic Soulfarer,
// the duration of this increase will be refreshed.
//
// You need to unlock the Passive Talent "Featherfall Judgment."
func (c *char) c1() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c1Key, 600), // 10s
		AffectedStat: attributes.AtkSpd,
		Amount: func() []float64 {
			if c.Core.Player.CurrentState() != action.NormalAttackState {
				return nil
			}
			return c.c1AtkSpd
		},
	})
}

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.c1AtkSpd = make([]float64, attributes.EndStatType)
	c.c1AtkSpd[attributes.AtkSpd] = 0.2

	if !c.revelation {
		return
	}
	c.c1EM = make([]float64, attributes.EndStatType)
	c.c1EM[attributes.EM] = 200

	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if prev != c.Index() {
			return
		}

		dur := c.StatusDuration(c1RevKey)
		if dur <= 0 {
			return
		}
		c.DeleteStatus(c1RevKey)
		nextChar := c.Core.Player.ByIndex(next)
		c.addC1RevelationBuff(nextChar, dur)
	}, "cyno-c1-revelation")
}

// Radiance: Stellar-Conduct: When Cyno enters the Pactsworn Pathclearer state,
// he will also gain the Together We Rise effect: Elemental Mastery is increased
// by 200. During this time, when switching characters from Cyno on the field,
// the incoming active character will take on the Together We Rise effect for
// the remaining duration.
func (c *char) c1RevelationOnPactsworn(dur int) {
	if c.Base.Cons < 1 {
		return
	}
	if !c.isStellarRadiance() {
		return
	}
	c.addC1RevelationBuff(c.CharWrapper, dur)
}

func (c *char) addC1RevelationBuff(ch *character.CharWrapper, dur int) {
	ch.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c1RevKey, dur),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return c.c1EM
		},
	})
}

func (c *char) c2RevelationOnPactsworn() {
	if c.Base.Cons < 1 {
		return
	}
	if !c.revelation {
		return
	}
	c.c2Stacks = 0
}

// When Cyno's Normal Attacks hit opponents, his Electro DMG Bonus will
// increase by 10% for 4s. This effect can be triggered once every 0.1s. Max 5
// stacks.
func (c *char) makeC2CB() info.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c2ICD) {
			return
		}
		c.AddStatus(c2ICD, 0.1*60, true)

		if !c.StatModIsActive(c2Key) {
			c.c2Stacks = 0
		}
		c.c2Stacks++
		if c.c2Stacks > 5 {
			c.c2Stacks = 5
		}

		m := make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2Key, 4*60),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				m[attributes.ElectroP] = 0.1 * float64(c.c2Stacks)
				return m
			},
		})
	}
}

// Additionally, when characters affected by Together We Rise hit opponents with
// their Normal Attacks or Charged Attacks, their Stellar-Conduct DMG dealt is
// increased by 16%. This effect, which can trigger once every 0.1s, will remain
// active for the remaining duration of Together We Rise. Max 5 stacks.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	if !c.revelation {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		default:
			return
		}

		ch := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		if !ch.StatModIsActive(c1RevKey) {
			return
		}

		// TODO: is the ICD on the triggering char or on Cyno? Probably doesn't
		// matter since the ICD is 0.1s
		if ch.StatusIsActive(c2RevICDKey) {
			return
		}
		ch.AddStatus(c2RevICDKey, 0.1*60, true)

		c.c2RevStacks = min(c.c2RevStacks+1, 5)
	}, "cyno-c2-revelation")

	for _, ch := range c.Core.Player.Chars() {
		ch.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBaseWithHitlag("cyno-c2-revelation", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag != attacks.AttackTagDirectStellarConduct {
					return 0
				}
				if !ch.StatusIsActive(c1RevKey) {
					return 0
				}

				return 0.16 * float64(c.c2RevStacks)
			},
		})
	}
}

// When Cyno is in the Pactsworn Pathclearer state triggered by Sacred Rite:
// Wolf's Swiftness, after he triggers Electro-Charged, Superconduct,
// Overloaded, Quicken, Aggravate, Hyperbloom, or an Electro Swirl reaction, he
// will restore 3 Elemental Energy for all nearby party members (except
// himself.)
// This effect can occur 5 times within one use of Sacred Rite: Wolf’s Swiftness.
//
// When other characters affected by Together We Rise trigger the aforementioned
// reaction, they will also regenerate 20 Elemental Energy for Cyno. This effect
// can trigger once every 18s.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	restore := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			if !c.revelation {
				return
			}

			ch := c.Core.Player.ByIndex(atk.Info.ActorIndex)

			if !ch.StatusIsActive(c1RevKey) {
				return
			}

			if c.StatusIsActive(c4RevICDKey) {
				return
			}
			c.AddStatus(c4RevICDKey, 18*60, true)
			c.AddEnergy("cyno-c4-revelation", 20)
			return
		}
		if c.c4Counter > 4 { // counting from 0 to 4, 5 instances max
			return
		}
		c.c4Counter++
		for _, this := range c.Core.Player.Chars() {
			// not for cyno
			if this.Index() != c.Index() {
				this.AddEnergy("cyno-c4", 3)
			}
		}
	}
	restoreNoGadget := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); ok {
			restore(args...)
		}
	}
	c.Core.Events.Subscribe(event.OnOverload, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnElectroCharged, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnLunarCharged, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSuperconduct, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnStellarConduct, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnQuicken, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnAggravate, restoreNoGadget, "cyno-c4")
	c.Core.Events.Subscribe(event.OnHyperbloom, restore, "cyno-c4")
	c.Core.Events.Subscribe(event.OnSwirlElectro, restoreNoGadget, "cyno-c4")
}

// After using Sacred Rite: Wolf's Swiftness or triggering the Judication effect of the Passive Talent "Featherfall Judgment,"
// Cyno will gain 4 stacks of the "Day of the Jackal" effect.
func (c *char) c6OnA1Skill() {
	if c.Base.Cons < 6 {
		return
	}
	c.AddStatus(c6Key, 8*60, true)
	c.c6Stacks += 4
	if c.c6Stacks > 8 {
		c.c6Stacks = 8
	}
}

// When he hits opponents with Normal Attacks, he will consume 1 stack of "Day of the Jackal" to fire off one Duststalker Bolt.
// "Day of the Jackal" lasts for 8s. Max 8 stacks. It will be canceled once Pactsworn Pathclearer ends.
// If Cyno is affected by Together We Rise, he will fire one Duststalker Bolt - Starsame instead.
// A maximum of 1 Duststalker Bolt can be unleashed this way every 0.4s.
// You must first unlock the Passive Talent "Featherfall Judgment."
func (c *char) makeC6CB() info.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}

	if c.c6Stacks == 0 {
		return nil
	}

	if !c.StatusIsActive(c6Key) {
		return nil
	}

	if !c.StatusIsActive(burstKey) {
		return nil
	}

	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		if c.c6Stacks == 0 {
			return
		}
		if !c.StatusIsActive(c6Key) {
			return
		}
		if c.StatusIsActive(c6ICDKey) {
			return
		}
		c.AddStatus(c6ICDKey, 0.4*60, true)
		c.c6Stacks--

		var ai info.AttackInfo

		if c.revelation && c.StatusIsActive(c1RevKey) {
			ai = info.AttackInfo{
				ActorIndex:       c.Index(),
				Abil:             "Raiment: Just Scales (C6)" + stellarConductText,
				AttackTag:        attacks.AttackTagDirectStellarConduct,
				StrikeType:       attacks.StrikeTypeSlash,
				Element:          attributes.Electro,
				Mult:             2.0,
				FlatDmg:          c.a4BoltSSC(),
				IgnoreDefPercent: 1,
			}
		} else {
			// technically should use ICDGroupCynoC6, but it's just reskinned standard ICD
			ai = info.AttackInfo{
				ActorIndex:   c.Index(),
				Abil:         "Raiment: Just Scales (C6)",
				AttackTag:    attacks.AttackTagElementalArtHold,
				ICDTag:       attacks.ICDTagElementalArt,
				ICDGroup:     attacks.ICDGroupDefault,
				StrikeType:   attacks.StrikeTypeSlash,
				Element:      attributes.Electro,
				Durability:   25,
				IsDeployable: true,
				Mult:         1.0,
				FlatDmg:      c.a4Bolt(),
			}
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.3,
			),
			0,
			0,
		)
	}
}
