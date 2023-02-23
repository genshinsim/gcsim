package shenhe

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	skillPressFrames []int
	skillHoldFrames  []int
)

const (
	skillPressCDStart  = 2
	skillPressHitmark  = 4
	skillHoldCDStart   = 31
	skillHoldHitmark   = 33
	holdParticleICDKey = "shenhe-hold-particle-icd"
	quillKey           = "shenhe-quill"
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(38) // walk
	skillPressFrames[action.ActionAttack] = 27
	skillPressFrames[action.ActionSkill] = 27
	skillPressFrames[action.ActionBurst] = 27
	skillPressFrames[action.ActionDash] = 21
	skillPressFrames[action.ActionJump] = 21
	skillPressFrames[action.ActionSwap] = 27

	// skill (hold) -> x
	// TODO: skill (hold) -> skill (hold) is 52 frames.
	skillHoldFrames = frames.InitAbilSlice(78) // walk
	skillHoldFrames[action.ActionAttack] = 45
	skillHoldFrames[action.ActionSkill] = 45 // assume skill (press)
	skillHoldFrames[action.ActionBurst] = 45
	skillHoldFrames[action.ActionDash] = 38
	skillHoldFrames[action.ActionJump] = 39
	skillHoldFrames[action.ActionSwap] = 44
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
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		skillPressHitmark,
		skillPressHitmark,
		c.makePressParticleCB(),
		c.makeC4ResetCB(),
	)

	if c.Base.Ascension >= 4 {
		c.Core.Tasks.Add(c.skillPressBuff, skillPressCDStart+1)
	}
	c.SetCDWithDelay(action.ActionSkill, 10*60, skillPressCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) makePressParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if done {
			return
		}
		// Skill actually moves you in game - actual catch is anywhere from 90-110 frames, take 100 as an average
		c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Cryo, c.ParticleDelay)
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Spring Spirit Summoning (Hold)",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Cryo,
		Durability: 50,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.5}, 4),
		skillHoldHitmark,
		skillHoldHitmark,
		c.holdParticleCB,
		c.makeC4ResetCB(),
	)

	if c.Base.Ascension >= 4 {
		c.Core.Tasks.Add(c.skillHoldBuff, skillHoldCDStart+1)
	}
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillHoldCDStart+1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) holdParticleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(holdParticleICDKey) {
		return
	}
	c.AddStatus(holdParticleICDKey, 0.5*60, true)
	// Particle spawn timing is a bit later than press E
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
}

// A4:
// After Shenhe uses Spring Spirit Summoning, she will grant all nearby party members the following effects:
//
// - Press: Elemental Skill and Elemental Burst DMG increased by 15% for 10s.
func (c *char) skillPressBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(quillKey, 10*60, true) // 10 sec duration
		char.SetTag(quillKey, 5)              // 5 quill on press
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-press", 10*60),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				switch a.Info.AttackTag {
				case attacks.AttackTagElementalArt:
				case attacks.AttackTagElementalArtHold:
				case attacks.AttackTagElementalBurst:
				default:
					return nil, false
				}
				return c.skillBuff, true
			},
		})
	}
}

// A4:
// After Shenhe uses Spring Spirit Summoning, she will grant all nearby party members the following effects:
//
// - Hold: Normal, Charged, and Plunging Attack DMG increased by 15% for 15s.
func (c *char) skillHoldBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(quillKey, 15*60, true) // 15 sec duration
		char.SetTag(quillKey, 7)              // 5 quill on hold
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("shenhe-a4-hold", 15*60),
			Amount: func(a *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				switch a.Info.AttackTag {
				case attacks.AttackTagNormal:
				case attacks.AttackTagExtra:
				case attacks.AttackTagPlunge:
				default:
					return nil, false
				}
				return c.skillBuff, true
			},
		})
	}
}

func (c *char) quillDamageMod() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		consumeStack := true
		if atk.Info.Element != attributes.Cryo {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagNormal:
			consumeStack = c.Base.Cons < 6
		case attacks.AttackTagExtra:
			consumeStack = c.Base.Cons < 6
		case attacks.AttackTagPlunge:
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
				c.Core.Log.NewEvent("Shenhe Quill proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("effect_ends_at", c.StatusExpiry(quillKey)).
					Write("quill_left", c.Tags[quillKey])
			}

			atk.Info.FlatDmg += amt
			if c.Base.Cons >= 4 {
				atk.Callbacks = append(atk.Callbacks, c.c4CB)
			}
		}

		return false
	}, "shenhe-quill-hook")
}
