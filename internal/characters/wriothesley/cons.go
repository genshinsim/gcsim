package wriothesley

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
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
func (c *char) c1(ai *combat.AttackInfo, snap *combat.Snapshot) (combat.AttackCBFunc, bool) {
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
	c.Core.Log.NewEvent("adding c1", glog.LogCharacterEvent, c.Index).Write("dmg%", dmg)

	// add c6 if applicable
	c.addC6Buff(snap)

	// return callback to heal, extend E, remove C1 and apply 2.5s cd
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
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
		if !c.c1SkillExtensionProc && c.StatusIsActive(skillKey) {
			c.ExtendStatus(skillKey, c1SkillExtension)
			c.c1SkillExtensionProc = true
			c.Core.Log.NewEvent("c1: skill duration is extended", glog.LogCharacterEvent, c.Index)
		}

		// heal
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Index,
			Message: "There Shall Be a Plea for Justice",
			Src:     c.caHeal * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})
	}, c.Base.Cons >= 6
}

func (c *char) makeC1N5CB() combat.AttackCBFunc {
	if c.Base.Cons < 1 || c.NormalCounter != 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
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
		c.Core.Log.NewEvent("gained Gracious Rebuke from C1 N5", glog.LogCharacterEvent, c.Index)
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
func (c *char) c2(snap *combat.Snapshot) {
	if c.Base.Cons < 2 {
		return
	}
	if !c.StatusIsActive(skillKey) {
		return
	}

	dmg := 0.4 * float64(c.a4Stack)
	snap.Stats[attributes.DmgP] += dmg
	c.Core.Log.NewEvent("adding c2", glog.LogCharacterEvent, c.Index).Write("dmg%", dmg)
}

// The HP restored to Wriothesley through Rebuke: Vaulting Fist will be increased to 50%
// of his Max HP. You must first unlock the Passive Talent "There Shall Be a Plea for Justice."
// Additionally, when Wriothesley is healed, if the amount of healing overflows, the following
// effects will occur depending on whether his is on the field or not. If he is on the field,
// his ATK SPD will be increased by 20% for 4s. If he is off-field, all party members' ATK SPD
// will be increased by 10% for 6s. These two ATK SPD increasing methods cannot stack.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		c.caHeal = 0.3
		return
	}
	c.caHeal = 0.5

	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if index != c.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		if overheal <= 0 {
			return false
		}

		chars := c.Core.Player.Chars()
		m := make([]float64, attributes.EndStatType)

		// remove old buffs
		for _, char := range chars {
			char.DeleteStatus(c4Status)
		}

		if c.Core.Player.Active() == c.Index {
			m[attributes.AtkSpd] = 0.2
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Status, 4*60),
				AffectedStat: attributes.AtkSpd,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		} else {
			m[attributes.AtkSpd] = 0.1
			for _, char := range chars {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(c4Status, 6*60),
					AffectedStat: attributes.AtkSpd,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
		}

		return false
	}, "wriothesley-c4-heal")
}

// The CRIT Rate of Rebuke: Vaulting Fist will be increased by 10%, and its CRIT DMG by 80%.
func (c *char) addC6Buff(snap *combat.Snapshot) {
	cr := 0.1
	cd := 0.8
	snap.Stats[attributes.CR] += cr
	snap.Stats[attributes.CD] += cd
	c.Core.Log.NewEvent("adding c6", glog.LogCharacterEvent, c.Index).Write("cr", cr).Write("cd", cd)
}
