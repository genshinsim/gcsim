package shenhe

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-5+i, f-5+i)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupPole,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charged[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-1, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
func (c *char) Skill(p map[string]int) (int, int) {
	hold := p["hold"]
	var f, a, cd int
	if hold == 1 {
		f, a = c.skillHold(p)
		cd = 15 * 60
	} else {
		f, a = c.skillPress(p)
		cd = 10 * 60
	}
	//press, hold -> 10s, 15s
	//hold, press -> 15s, 10s
	//press, press -> 10s, 10s
	//hold, hold -> 15s, 15s

	c.SetCD(core.ActionSkill, cd)
	return f, a
}

func (c *char) skillPress(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Press)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	// c.AddTask(c.skillPressBuff, "shenhe (press) quill start", f-1)
	c.skillPressBuff()
	c.Core.Status.AddStatus(quillKey, 10*60)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f, f)

	// Skill actually moves you in game - actual catch is anywhere from 90-110 frames, take 100 as an average
	c.QueueParticle("shenhe", 3, core.Cryo, 100)

	return f, a
}

func (c *char) skillHold(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Hold)",
		AttackTag:  core.AttackTagElementalArtHold,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	// c.AddTask(c.skillHoldBuff, "shenhe (hold) quill start", f-1)
	c.skillHoldBuff()
	c.Core.Status.AddStatus(quillKey, 15*60)
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.5, false, core.TargettableEnemy), f, f)

	// Particle spawn timing is a bit later than press E
	c.QueueParticle("shenhe", 4, core.Cryo, 115)

	return f, a
}

// Burst attack damage queue generator
func (c *char) Burst(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionBurst, p)

	// TODO: Not 100% sure if this shares ICD with the DoT, currently coded that it does
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (Hit 1)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	x, y := c.Core.Targets[0].Shape().Pos()

	//duration is 12 second (extended by c2 by 6s)
	dur := 12 * 60
	count := 6
	if c.Base.Cons >= 2 {
		dur += 6 * 60
		count += 3

		// Active characters within the skill's field deals 15% increased Cryo CRIT DMG.
		// TODO: Exact mechanics of how this works is unknown. Not sure if it works like Gorou E/Bennett Q
		// For now, assume that it operates like Kazuha C2, and extends for 2s after burst ends like the res shred
		val := make([]float64, core.EndStatType)
		val[core.CD] = 0.15
		for _, char := range c.Core.Chars {
			this := char
			char.AddPreDamageMod(core.PreDamageMod{
				Key:    "shenhe-c2",
				Expiry: c.Core.F + dur + 2*60,
				Amount: func(ae *core.AttackEvent, t core.Target) ([]float64, bool) {
					if ae.Info.Element != core.Cryo {
						return nil, false
					}

					switch this.CharIndex() {
					case c.Core.ActiveChar, c.CharIndex():
						return val, true
					}
					return nil, false
				},
			})
		}
	}
	// Res shred persists for 2 seconds after burst ends
	cb := func(a core.AttackCB) {
		a.Target.AddResMod("Shenhe Burst Shred (Cryo)", core.ResistMod{
			Duration: dur + 2*60,
			Ele:      core.Cryo,
			Value:    -burstrespp[c.TalentLvlBurst()],
		})
	}
	cb2 := func(a core.AttackCB) {
		a.Target.AddResMod("Shenhe Burst Shred (Phys)", core.ResistMod{
			Duration: dur + 2*60,
			Ele:      core.Physical,
			Value:    -burstrespp[c.TalentLvlBurst()],
		})
	}
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, core.TargettableEnemy), 0, 15, cb, cb2)

	ai = core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Divine Maiden's Deliverance (DoT)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Cryo,
		Durability: 25,
		Mult:       burstdot[c.TalentLvlBurst()],
	}

	c.AddTask(func() {
		snap := c.Snapshot(&ai)
		c.Core.Status.AddStatus("shenheburst", dur)
		//TODO: check this accuracy? Siri's sheet has 137 per
		// dot every 2 second, double tick shortly after another
		for i := 0; i < count; i++ {
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewCircleHit(0, 0, 5, false, core.TargettableEnemy), i*120+50)
			c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewCircleHit(0, 0, 5, false, core.TargettableEnemy), i*120+80)
		}
	}, "shenhe-snapshot", f+2)

	c.SetCDWithDelay(core.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)

	return f, a
}

func (c *char) skillPressBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		c.quillcount[i] = 5
		c.updateBuffTags()
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-a1-press",
			Expiry: c.Core.F + 10*60,
			Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
				if a.Info.AttackTag != core.AttackTagElementalBurst && a.Info.AttackTag != core.AttackTagElementalArt && a.Info.AttackTag != core.AttackTagElementalArtHold {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) skillHoldBuff() {
	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15
	for i, char := range c.Core.Chars {
		c.quillcount[i] = 7
		c.updateBuffTags()
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-a1-hold",
			Expiry: c.Core.F + 15*60,
			Amount: func(a *core.AttackEvent, t core.Target) ([]float64, bool) {
				if a.Info.AttackTag != core.AttackTagNormal && a.Info.AttackTag != core.AttackTagExtra && a.Info.AttackTag != core.AttackTagPlunge {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) quillDamageMod() {

	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		consumeStack := true
		if atk.Info.Element != core.Cryo {
			return false
		}

		switch atk.Info.AttackTag {
		case core.AttackTagElementalBurst:
		case core.AttackTagElementalArt:
		case core.AttackTagElementalArtHold:
		case core.AttackTagNormal:
			consumeStack = c.Base.Cons < 6
		case core.AttackTagExtra:
			consumeStack = c.Base.Cons < 6
		case core.AttackTagPlunge:
		default:
			return false
		}

		if c.Core.Status.Duration(quillKey) == 0 {
			return false
		}

		if c.quillcount[atk.Info.ActorIndex] > 0 {
			// ai := core.AttackInfo{
			// 	Abil:      "Quills",
			// 	AttackTag: core.AttackTagNone,
			// }
			stats, _ := c.SnapshotStats()
			amt := skillpp[c.TalentLvlSkill()] * ((c.Base.Atk+c.Weapon.Atk)*(1+stats[core.ATKP]) + stats[core.ATK])
			if consumeStack { //c6
				c.quillcount[atk.Info.ActorIndex]--
				c.updateBuffTags()
			}
			c.Core.Log.NewEvent(
				"Shenhe Quill proc dmg add",
				core.LogPreDamageMod,
				atk.Info.ActorIndex,
				"before", atk.Info.FlatDmg,
				"addition", amt,
				"effect_ends_at", c.Core.Status.Duration(quillKey),
				"quills left", c.quillcount[atk.Info.ActorIndex],
			)
			atk.Info.FlatDmg += amt
			if c.Base.Cons >= 4 {
				if c.c4count < 50 {
					c.c4count++
				}
				c.c4expiry = c.Core.F + 60*60
			}
		}

		return false //not sure of correctness here
	}, "shenhe-quill")
}
