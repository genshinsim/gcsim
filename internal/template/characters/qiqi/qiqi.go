package qiqi

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.Qiqi, NewChar)
}

type char struct {
	*character.Tmpl
	// talismanExpiry    []int
	// talismanICDExpiry []int
	c4ICDExpiry       int
	skillLastUsed     int
	skillHealSnapshot core.Snapshot // Required as both on hit procs and continuous healing need to use this
}

const (
	talismanKey    = "qiqi-talisman"
	talismanICDKey = "qiqi-talisman-icd"
)

// TODO: Not implemented - C6 (revival mechanic, not suitable for sim)
// C4 - Enemy Atk reduction, not useful in this sim version
func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSword
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	c.skillLastUsed = 0

	return &c, nil
}

// Ensures the set of targets are initialized properly
func (c *char) Init() {
	c.Tmpl.Init()

	c.talismanHealHook()
	c.onNACAHitHook()
	c.a1()

	if c.Base.Cons >= 2 {
		c.c2()
	}
}

//Qiqi's Normal and Charge Attack DMG against opponents affected by Cryo is increased by 15%.
func (c *char) c2() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = .15
	c.AddPreDamageMod(core.PreDamageMod{
		Key:    "qiqi-c2",
		Expiry: -1,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
				return nil, false
			}
			if !t.AuraContains(core.Cryo, core.Frozen) {
				return nil, false
			}
			return val, true
		},
	})
}

func (c *char) talismanHealHook() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		//do nothing if talisman expired
		if t.GetTag(talismanKey) < c.Core.F {
			return false
		}
		//do nothing if talisman still on icd
		if t.GetTag(talismanICDKey) >= c.Core.F {
			return false
		}

		atk := args[1].(*core.AttackEvent)

		healAmt := c.healDynamic(burstHealPer, burstHealFlat, c.TalentLvlBurst())
		c.Core.Health.Heal(core.HealInfo{
			Caller:  c.Index,
			Target:  atk.Info.ActorIndex,
			Message: "Fortune-Preserving Talisman",
			Src:     healAmt,
			Bonus:   c.Stat(core.Heal),
		})
		t.SetTag(talismanICDKey, c.Core.F+60)

		// c.Core.Log.NewEvent(
		// 	"Qiqi Talisman Healing",
		// 	core.LogCharacterEvent,
		// 	c.Index,
		// 	"target", t.Index(),
		// 	"healed_char", atk.Info.ActorIndex,
		// 	"talisman_expiry", t.GetTag(talismanKey),
		// 	"talisman_healing_icd", t.GetTag(talismanICDKey),
		// 	"healed_amt", healAmt,
		// )

		return false
	}, "talisman-heal-hook")
}

// Handles C2, A4, and skill NA/CA on hit hooks
// Additionally handles burst Talisman hook - can't be done another way since Talisman is applied before the burst damage is dealt
func (c *char) onNACAHitHook() {
	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		t := args[0].(core.Target)
		atk := args[1].(*core.AttackEvent)

		if atk.Info.ActorIndex != c.Index {
			return false
		}

		// Talisman is applied before the damage is dealt
		if atk.Info.Abil == "Fortune-Preserving Talisman" {
			// c.talismanExpiry[t.Index()] = c.Core.F + 15*60
			t.SetTag(talismanKey, c.Core.F+15*60)
		}

		// All of the below only occur on Qiqi NA/CA hits
		if !((atk.Info.AttackTag == core.AttackTagNormal) || (atk.Info.AttackTag == core.AttackTagExtra)) {
			return false
		}

		// A4
		// When Qiqi hits opponents with her Normal and Charged Attacks, she has a 50% chance to apply a Fortune-Preserving Talisman to them for 6s. This effect can only occur once every 30s.
		if (c.c4ICDExpiry <= c.Core.F) && (c.Rand.Float64() < 0.5) {
			// Don't want to overwrite a longer burst duration talisman with a shorter duration one
			// TODO: Unclear how the interaction works if there is already a talisman on enemy
			// TODO: Being generous for now and not putting it on CD if there is a conflict
			if t.GetTag(talismanKey) < c.Core.F+360 {
				t.SetTag(talismanKey, c.Core.F+360)
				c.c4ICDExpiry = c.Core.F + 30*60
				c.Core.Log.NewEvent(
					"Qiqi A4 Adding Talisman",
					core.LogCharacterEvent,
					c.Index,
					"target", t.Index(),
					"talisman_expiry", t.GetTag(talismanKey),
					"c4_icd_expiry", c.c4ICDExpiry,
				)
			}
		}

		// Qiqi NA/CA healing proc in skill duration
		if c.Core.Status.Duration("qiqiskill") > 0 {
			c.Core.Health.Heal(core.HealInfo{
				Caller:  c.Index,
				Target:  -1,
				Message: "Herald of Frost (Attack)",
				Src:     c.healSnapshot(&c.skillHealSnapshot, skillHealOnHitPer, skillHealOnHitFlat, c.TalentLvlSkill()),
				Bonus:   c.skillHealSnapshot.Stats[core.Heal],
			})
		}

		return false
	}, "qiqi-onhit-naca-hook")
}

// Implements event hook and incoming healing bonus function for A1
// TODO: Could possibly change this so the AddIncHealBonus occurs at start, then event subscription occurs upon using Qiqi skill?
// TODO: Likely more efficient to not maintain event subscription always, but grouping the two for clarity currently
// When a character under the effects of Adeptus Art: Herald of Frost triggers an Elemental Reaction, their Incoming Healing Bonus is increased by 20% for 8s.
func (c *char) a1() {

	// Add to incoming healing bonus array
	c.Core.Health.AddIncHealBonus(func(healedCharIndex int) float64 {
		healedCharName := c.Core.Chars[healedCharIndex].Name()

		if c.Core.Status.Duration("qiqia1"+healedCharName) > 0 {
			return .2
		}
		return 0
	})

	a1Hook := func(args ...interface{}) bool {
		if c.Core.Status.Duration("qiqiskill") == 0 {
			return false
		}
		atk := args[1].(*core.AttackEvent)

		// Active char is the only one under the effects of Qiqi skill
		if atk.Info.ActorIndex != c.Core.ActiveChar {
			return false
		}

		c.Core.Status.AddStatus("qiqia1"+c.Core.Chars[c.Core.ActiveChar].Name(), 8*60)

		return false
	}
	for i := core.EventType(core.ReactionEventStartDelim + 1); i < core.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, a1Hook, "qiqi-a1")
	}
	// c.Core.Events.Subscribe(core.OnTransReaction, a1Hook, "qiqi-a1")
	// c.Core.Events.Subscribe(core.OnAmpReaction, a1Hook, "qiqi-a1")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", core.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
