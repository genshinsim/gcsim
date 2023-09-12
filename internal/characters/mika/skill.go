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

var skillPressFrames []int
var skillHoldFrames []int

const (
	skillPressCDStart = 16
	skillPressHitmark = 17
	skillPressTravel  = 1

	skillHoldCDStart    = 11
	skillHoldHitmark    = 12
	skillHoldTravel     = 3
	rimestarShardTravel = 46

	skillBuffKey      = "soulwind"
	skillBuffDuration = 12 * 60
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(39) // E -> N1/Q
	skillPressFrames[action.ActionDash] = 34
	skillPressFrames[action.ActionJump] = 35
	skillPressFrames[action.ActionWalk] = 19
	skillPressFrames[action.ActionSwap] = 37

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(46) // E -> Swap
	skillHoldFrames[action.ActionAttack] = 38
	skillHoldFrames[action.ActionBurst] = 37
	skillHoldFrames[action.ActionDash] = 30
	skillHoldFrames[action.ActionJump] = 30
	skillHoldFrames[action.ActionWalk] = 30
}

func (c *char) Skill(p map[string]int) action.Info {
	if p["hold"] != 0 {
		return c.skillHold()
	}
	return c.skillPress()
}

func (c *char) skillPress() action.Info {
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
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.4}, 2.5, 10),
		skillPressHitmark,
		skillPressHitmark+skillPressTravel,
		c.makeParticleCB(),
		a1CB,
		c.c2(),
	)

	c.QueueCharTask(c.applyBuffs, skillPressCDStart)
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillPressCDStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionWalk], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHold() action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Rimestar Flare",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
		skillHoldHitmark,
		skillHoldHitmark+skillHoldTravel,
		c.makeParticleCB(),
		c.makeRimestarShardsCB(),
		c.c2(),
	)

	c.QueueCharTask(c.applyBuffs, skillHoldCDStart+1)
	c.SetCDWithDelay(action.ActionSkill, 15*60, skillHoldCDStart+1)

	return action.Info{
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
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       skillExplode[c.TalentLvlSkill()],
		}

		enemies := c.Core.Combat.RandomEnemiesWithinArea(
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10),
			func(t combat.Enemy) bool { return a.Target.Key() != t.Key() },
			3,
		)
		for i := 0; i < len(enemies); i++ {
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
				0,
				rimestarShardTravel,
				a1CB,
			)
		}
	}
}

func (c *char) applyBuffs() {
	c.SetTag(a1Stacks, 0)
	c.skillBuff()

	if c.Base.Ascension >= 4 {
		c.a4Stack = false
	}
}

func (c *char) skillBuff() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(skillBuffKey, skillBuffDuration),
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				return c.skillbuff, true
			},
		})

		if c.Base.Ascension >= 1 {
			c.a1(char)
		}

		if c.Base.Cons >= 6 {
			c.c6(char)
		}
	}
}
