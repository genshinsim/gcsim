package mika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// based on shenhe frames
// TODO: update frames, hitlags & hitboxes
var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressCDStart    = 2
	skillPressHitmark    = 4
	skillHoldCDStart     = 31
	skillHoldHitmark     = 33
	rimestarShardHitmark = 30

	skillBuffKey = "soulwind"
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
		Abil:               "Flowfrost Arrow",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillPress[c.TalentLvlSkill()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	var a1CB combat.AttackCBFunc
	if c.Base.Ascension >= 1 {
		gen := false
		a1CB = func(a combat.AttackCB) {
			if a.Target.Type() != targets.TargettableEnemy {
				return
			}
			if !gen { // ignore a first enemy
				gen = true
				return
			}
			c.addDetectorStack()
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.2}, 4, 8),
		skillPressHitmark,
		skillPressHitmark,
		c.makeParticleCB(),
		a1CB,
		c.c2(),
	)

	c.QueueCharTask(func() {
		c.SetTag(a1Stacks, 0)
		c.skillBuff()

		if c.Base.Ascension >= 1 {
			c.a1()
		}
		if c.Base.Ascension >= 4 {
			c.a4Stack = false
		}
	}, skillPressCDStart+1)
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillPressCDStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rimestar Flare",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	enemy := c.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), nil)
	if enemy != nil {
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(enemy.Key()),
			skillHoldHitmark,
			skillHoldHitmark,
			c.makeParticleCB(),
			c.makeRimestarShardsCB(),
			c.c2(),
		)
	}

	c.QueueCharTask(func() {
		c.SetTag(a1Stacks, 0)
		c.skillBuff()

		if c.Base.Ascension >= 1 {
			c.a1()
		}
		if c.Base.Ascension >= 4 {
			c.a4Stack = false
		}
	}, skillHoldCDStart+1)
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillHoldCDStart+1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Cryo, c.ParticleDelay)
	}
}

func (c *char) makeRimestarShardsCB() func(combat.AttackCB) {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Rimestar Shard",
			AttackTag:  attacks.AttackTagElementalArtHold,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillExplode[c.TalentLvlSkill()],
		}

		// TODO: radius? enemies should be sorted by distance?
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), func(t combat.Enemy) bool {
			return a.Target.Key() != t.Key()
		})
		for i := 0; i < 3; i++ {
			if i == len(enemies) {
				break
			}

			var a1CB combat.AttackCBFunc
			if c.Base.Ascension >= 1 {
				done := false
				a1CB = func(a combat.AttackCB) {
					if a.Target.Type() != targets.TargettableEnemy {
						return
					}
					if done {
						return
					}
					done = true
					c.addDetectorStack()
				}
			}

			c.Core.QueueAttack(
				ai,
				combat.NewSingleTargetHit(enemies[i].Key()),
				rimestarShardHitmark,
				rimestarShardHitmark,
				a1CB,
			)
		}
	}
}

func (c *char) skillBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(skillBuffKey, 12*60),
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				return c.skillbuff, true
			},
		})

		if c.Base.Cons >= 6 {
			c.c6(char)
		}
	}
}
