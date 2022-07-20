package shenhe

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
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
		ActorIndex:         c.Index,
		Abil:               "Spring Spirit Summoning (Press)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.skillPressBuff()
	c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), skillPressHitmark, skillPressHitmark)

	// Skill actually moves you in game - actual catch is anywhere from 90-110 frames, take 100 as an average
	c.Core.QueueParticle("shenhe", 3, attributes.Cryo, c.Core.Flags.ParticleDelay)

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
	c.Core.QueueParticle("shenhe", 4, attributes.Cryo, 15+c.Core.Flags.ParticleDelay)

	c.SetCD(action.ActionSkill, 15*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldHitmark,
		State:           action.SkillState,
	}
}

func (c *char) skillPressBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(quillKey, 600, true) //10 sec duration
		char.SetTag(quillKey, 5)            // 5 quill on press
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-press", 600),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if a.Info.AttackTag != combat.AttackTagElementalBurst && a.Info.AttackTag != combat.AttackTagElementalArt && a.Info.AttackTag != combat.AttackTagElementalArtHold {
					return nil, false
				}
				return c.skillBuff, true
			},
		})
	}
}

func (c *char) skillHoldBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(quillKey, 900, true) //15 sec duration
		char.SetTag(quillKey, 7)            // 5 quill on press
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-hold", 15*60),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if a.Info.AttackTag != combat.AttackTagNormal && a.Info.AttackTag != combat.AttackTagExtra && a.Info.AttackTag != combat.AttackTagPlunge {
					return nil, false
				}
				return c.skillBuff, true
			},
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

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if !char.StatusIsActive(quillKey) {
			return false
		}

		if char.Tags[quillKey] > 0 {
			stats, _ := c.Stats()
			amt := skillpp[c.TalentLvlSkill()] * ((c.Base.Atk+c.Weapon.Atk)*(1+stats[attributes.ATKP]) + stats[attributes.ATK])
			if consumeStack { //c6
				char.Tags[quillKey]--
			}

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent(
					"Shenhe Quill proc dmg add",
					glog.LogPreDamageMod,
					atk.Info.ActorIndex,
				).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("effect_ends_at", c.Core.Status.Duration(quillKey)).
					Write("quills left", c.Tags[quillKey])
			}

			atk.Info.FlatDmg += amt
			if c.Base.Cons >= 4 {
				//reset stacks to zero if all expired
				if !c.StatusIsActive(c4BuffKey) {
					c.c4count = 0
				}
				if c.c4count < 50 {
					c.c4count++
				}
				c.AddStatus(c4BuffKey, 3600, true) // 60 s
			}
		}

		return false
	}, "shenhe-quill")
}
