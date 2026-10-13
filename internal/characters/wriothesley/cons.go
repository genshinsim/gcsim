package wriothesley

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Status         = "wriothesley-c1"
	c1ICD            = 2.5 * 60
	c1ICDKey         = "wriothesley-c1-icd"
	c1SkillExtension = 4 * 60
	c4Status         = "wriothesley-c4-spd"
)

func (c *char) c1Ready() bool {
	return c.c1N5Proc || (c.CurrentHPRatio() < 0.6 && !c.StatusIsActive(c1ICDKey))
}

// The Gracious Rebuke from the Passive Talent "There Shall Be a Plea for Justice" is changed to this:
// When Wriothesley's HP is less than 60% or while he is in the Chilling Penalty state caused by Icefang Rush,
// when the fifth attack of Repelling Fists hits, it will create a Gracious Rebuke.
// 1 Gracious Rebuke effect can be obtained every 2.5s.
// Additionally, Rebuke: Vaulting Fist will obtain the following enhancement:
// - The DMG Bonus gained will be further increased to 200%.
// - When it hits while Wriothesley is in the Chilling Penalty state, that state's duration is extended by 4s. 1 such extension can occur per 1 Chilling Penalty duration.
// You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
func (c *char) c1(ai *info.AttackInfo, snap *info.Snapshot) (info.AttackCBFunc, bool) {
	if !c.c1Ready() {
		return nil, false
	}
	c.c1N5Proc = false

	// add status that is removed on consumption
	c.AddStatus(c1Status, -1, false)

	// adjust ai
	ai.Abil = "Rebuke: Vaulting Fist"
	ai.HitlagFactor = 0.03
	ai.HitlagHaltFrames = 0.12 * 60

	// 200% increased DMG
	dmg := 2.0
	snap.Stats[attributes.DmgP] += dmg
	c.Core.Log.NewEvent("adding c1", glog.LogCharacterEvent, c.Index()).Write("dmg%", dmg)

	// add c6 if applicable
	c.addC6Buff(snap)

	// return callback to heal, extend E, remove C1 and apply 2.5s cd
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		// do not proc if C1 not active
		if !c.StatusIsActive(c1Status) {
			return
		}
		// remove C1 and apply CD
		c.DeleteStatus(c1Status)
		c.AddStatus(c1ICDKey, c1ICD, true)

		// E extension
		c.c1SkillExtension()

		// heal
		c.Core.Player.Heal(info.HealInfo{
			Caller:  c.Index(),
			Target:  c.Index(),
			Message: a1HealName,
			Src:     c.rebukeHeal * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}, c.Base.Cons >= 6
}

func (c *char) makeC1N5CB() info.AttackCBFunc {
	if c.Base.Cons < 1 || c.NormalCounter != 4 {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		if c.a1StellarEnabled() {
			c.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBase(c1Status+"-revelation-ca", -1),
				Amount: func(ai info.AttackInfo) float64 {
					if slices.Contains(ai.AdditionalTags, attacks.AttackTagWrioSSCCA) {
						c.Core.Tasks.Add(func() {
							c.DeleteReactBonusMod(c1Status + "-revelation-ca")
						}, 0)
						return 0.5
					}
					return 0
				},
			})
			return
		}
		// check if E active
		if !c.StatusIsActive(skillKey) {
			return
		}
		// check CD
		if c.StatusIsActive(c1ICDKey) {
			return
		}
		c.c1N5Proc = true
		c.Core.Log.NewEvent("gained Gracious Rebuke from C1 N5", glog.LogCharacterEvent, c.Index())
	}
}

// Radiance: Stellar-Conduct:
// After a Luster: Vaulting Fist attack hits an opponent, the next fifth attack of Repelling Fist
// enhanced by Chilling Penalty that hits within the next 5s deals 50% increased Stellar-Conduct DMG.
// After the fifth attack of Repelling Fist enhanced by Chilling Penalty hits an opponent, Luster:
// Vaulting Fists that hit within the next 5s deal 50% increased DMG.
// If Luster: Vaulting Fist hits an opponent while Wriothesley is in the Chilling Penalty state,
// that state's duration is extended by 4s.
func (c *char) c1StellarCACB(a info.AttackCB) {
	if c.Base.Cons < 1 {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if !c.a1StellarEnabled() {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(c1Status+"-revelation-n5", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if slices.Contains(ai.AdditionalTags, attacks.AttackTagWrioSSCN5) {
				c.Core.Tasks.Add(func() {
					c.DeleteReactBonusMod(c1Status + "-revelation-n5")
				}, 0)
				return 0.5
			}
			return 0
		},
	})

	c.c1SkillExtension()
}

func (c *char) c1SkillExtension() {
	if !c.c1SkillExtensionProc && c.StatusIsActive(skillKey) {
		c.ExtendStatus(skillKey, c1SkillExtension)
		c.c1SkillExtensionProc = true
		c.Core.Log.NewEvent("c1: skill duration is extended", glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) resetC1SkillExtension() {
	if c.Base.Cons < 1 {
		return
	}
	c.c1SkillExtensionProc = false
}

// When using Darkgold Wolfbite, each Prosecution Edict stack from the Passive Talent
// "There Shall Be a Reckoning for Sin" will increase said ability's DMG dealt by 40%.
// You must first unlock the Passive Talent "There Shall Be a Reckoning for Sin."
func (c *char) c2(snap *info.Snapshot) {
	if c.Base.Ascension < 4 {
		return
	}
	if c.Base.Cons < 2 {
		return
	}
	if !c.StatusIsActive(skillKey) {
		return
	}

	dmg := 0.4 * float64(c.a4Stack)
	snap.Stats[attributes.DmgP] += dmg
	c.Core.Log.NewEvent("adding c2", glog.LogCharacterEvent, c.Index()).Write("dmg%", dmg)
}

func (c *char) c2RevelationSSCMult() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	if !c.revelation {
		return 0
	}

	if c.Base.Ascension < 4 {
		return 0
	}

	if c.a4Stack < 5 {
		return 0
	}

	if !c.skillBuffActive() {
		return 0
	}

	if !c.isStellarRadiance() {
		return 0
	}

	return 0.5
}

func (c *char) c2RevelationNAMult() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	if !c.revelation {
		return 0
	}

	if c.Base.Ascension < 4 {
		return 0
	}

	if !c.skillBuffActive() {
		return 0
	}

	if c.isStellarRadiance() {
		return 0
	}

	return 0.25
}

func (c *char) c2RevelationCAMult() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	if !c.revelation {
		return 0
	}

	if c.Base.Ascension < 4 {
		return 0
	}

	if !c.skillBuffActive() {
		return 0
	}

	if c.isStellarRadiance() {
		return 0
	}

	return 0.3
}

// The HP restored to Wriothesley through Rebuke: Vaulting Fist will be increased to 50%
// of his Max HP. You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
// Additionally, when Wriothesley is healed, if the amount of healing overflows, the following
// effects will occur depending on whether his is on the field or not. If he is on the field,
// his ATK SPD will be increased by 20% for 4s. If he is off-field, all party members' ATK SPD
// will be increased by 10% for 6s. These two ATK SPD increasing methods cannot stack.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		c.rebukeHeal = 0.3
		return
	}
	c.rebukeHeal = 0.5

	m := make([]float64, attributes.EndStatType)
	c.Core.Events.Subscribe(event.OnHeal, func(args ...any) {
		if c.revelation && c.isStellarRadiance() {
			return
		}
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if index != c.Index() {
			return
		}
		if amount <= 0 {
			return
		}
		if overheal <= 0 {
			return
		}

		chars := c.Core.Player.Chars()

		// remove old buffs
		for _, char := range chars {
			char.DeleteStatus(c4Status)
		}

		if c.Core.Player.Active() == c.Index() {
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Status, 4*60),
				AffectedStat: attributes.AtkSpd,
				Amount: func() []float64 {
					m[attributes.AtkSpd] = 0.2
					return m
				},
			})
		} else {
			for _, char := range chars {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(c4Status, 6*60),
					AffectedStat: attributes.AtkSpd,
					Amount: func() []float64 {
						m[attributes.AtkSpd] = 0.1
						return m
					},
				})
			}
		}
	}, "wriothesley-c4-heal")

	if !c.revelation {
		return
	}

	for _, ch := range c.Core.Player.Chars() {
		amt := 0.1
		if ch.Index() == c.Index() {
			amt = 0.2
		}

		ch.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c4Status+"-revelation", -1),
			AffectedStat: attributes.AtkSpd,
			Amount: func() []float64 {
				if c.StatusIsActive(c4Status) {
					return nil
				}
				if !c.isStellarRadiance() {
					return nil
				}
				m[attributes.AtkSpd] = amt
				return m
			},
		})
	}
}

// The CRIT Rate of Rebuke: Vaulting Fist will be increased by 10%, and its CRIT DMG by 80%. When
// unleashed, it will also create an additional icicle that deals 100% of Rebuke: Vaulting Fist's
// original DMG as Cryo DMG. DMG dealt this way is regarded as Charged Attack DMG.
//
// Radiance: Stellar-Conduct: The aforementioned effect is changed to:
// The Repelling Fist enhanced by Chilling Penalty, as well as Luster: Vaulting Fist, gain a 10%
// CRIT Rate and 80% CRIT DMG boost. In addition, the fifth attack of Repelling Fists enhanced by
// Chilling Penalty, as well as Luster: Vaulting Fist, will each generate an additional icicle that
// deals Cryo DMG at 60% of its original DMG. This DMG will be considered Stellar-Conduct DMG.
func (c *char) addC6Buff(snap *info.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}
	cr := 0.1
	cd := 0.8
	snap.Stats[attributes.CR] += cr
	snap.Stats[attributes.CD] += cd
	c.Core.Log.NewEvent("adding c6", glog.LogCharacterEvent, c.Index()).Write("cr", cr).Write("cd", cd)
}

func (c *char) c6OnN5(ai *info.AttackInfo, snap *info.Snapshot, ap *info.AttackPattern) {
	if c.Base.Cons < 6 {
		return
	}
	if !c.revelation {
		return
	}
	if !c.skillBuffActive() {
		return
	}
	if !c.a1StellarEnabled() {
		return
	}

	ai.Abil = "Normal 5 (Enhanced) (C6)" + stellarConductText
	ai.Mult *= 0.6
	ai.IsDeployable = true
	c.Core.QueueAttackWithSnap(*ai, *snap, *ap, 0, c.particleCB)
}

func (c *char) c6OnStellarCA(ai *info.AttackInfo, snap *info.Snapshot, ap *info.AttackPattern) {
	if c.Base.Cons < 6 {
		return
	}
	if !c.revelation {
		return
	}
	if !c.skillBuffActive() {
		return
	}
	if !c.a1StellarEnabled() {
		return
	}

	ai.Abil = "Luster: Vaulting Fist (C6)" + stellarConductText
	ai.Mult *= 0.6
	ai.IsDeployable = true
	c.Core.QueueAttackWithSnap(*ai, *snap, *ap, 0)
}
