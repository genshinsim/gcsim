package shenhe

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var skillPressFrames []int
var skillHoldFrames []int

const skillPressHitmark = 31
const skillHoldHitmark = 44
const quillKey = "shenhequill"

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(31)

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(44)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	if p["hold"] != 0 {
		return c.skillHold(p)
	}
	return c.skillPress(p)
}

func (c *char) skillPress(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Press)",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	c.skillPressBuff()
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), skillPressHitmark, skillPressHitmark)

	// Skill actually moves you in game - actual catch is anywhere from 90-110 frames, take 100 as an average
	c.Core.QueueParticle("shenhe", 3, attributes.Cryo, 100)

	c.SetCD(action.ActionSkill, 10*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressHitmark,
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Hold)",
		AttackTag:  combat.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.skillHoldBuff()
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.5, false, combat.TargettableEnemy), skillHoldHitmark, skillHoldHitmark)

	// Particle spawn timing is a bit later than press E
	c.Core.QueueParticle("shenhe", 4, attributes.Cryo, 115)

	c.SetCD(action.ActionSkill, 15*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldHitmark,
		State:           action.SkillState,
	}
}

// Helper function to update tags that can be used in configs
// Should be run whenever c.quillcount is updated
func (c *char) updateBuffTags() {
	for _, char := range c.Core.Player.Chars() {
		c.Tags["quills_"+char.Base.Name] = c.quillcount[char.Index]
		c.Tags[fmt.Sprintf("quills_%v", char.Index)] = c.quillcount[char.Index]
	}
}

func (c *char) skillPressBuff() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.15
	for i := range c.Core.Player.Chars() {
		c.quillcount[i] = 5
	}
	c.updateBuffTags()

	c.Core.Status.Add(quillKey, 10*60)

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod("shenhe-a4-press", 10*60, func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if a.Info.AttackTag != combat.AttackTagElementalBurst && a.Info.AttackTag != combat.AttackTagElementalArt && a.Info.AttackTag != combat.AttackTagElementalArtHold {
				return nil, false
			}
			return m, true
		})
	}
}

func (c *char) skillHoldBuff() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.15
	for i := range c.Core.Player.Chars() {
		c.quillcount[i] = 7
	}
	c.updateBuffTags()

	c.Core.Status.Add(quillKey, 15*60)

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod("shenhe-a4-hold", 15*60, func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if a.Info.AttackTag != combat.AttackTagNormal && a.Info.AttackTag != combat.AttackTagExtra && a.Info.AttackTag != combat.AttackTagPlunge {
				return nil, false
			}
			return m, true
		})
	}
}

func (c *char) quillDamageMod() {
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		consumeStack := true
		if atk.Info.Element != attributes.Cryo {
			return false
		}

		switch atk.Info.AttackTag {
		case combat.AttackTagElementalBurst:
		case combat.AttackTagElementalArt:
		case combat.AttackTagElementalArtHold:
		case combat.AttackTagNormal:
			consumeStack = c.Base.Cons < 6
		case combat.AttackTagExtra:
			consumeStack = c.Base.Cons < 6
		case combat.AttackTagPlunge:
		default:
			return false
		}

		if c.Core.Status.Duration(quillKey) == 0 {
			return false
		}

		if c.quillcount[atk.Info.ActorIndex] > 0 {
			stats, _ := c.Stats()
			amt := skillpp[c.TalentLvlSkill()] * ((c.Base.Atk+c.Weapon.Atk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK])
			if consumeStack { //c6
				c.quillcount[atk.Info.ActorIndex]--
				c.updateBuffTags()
			}
			c.Core.Log.NewEvent(
				"Shenhe Quill proc dmg add",
				glog.LogPreDamageMod,
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

		return false
	}, "shenhe-quill")
}
