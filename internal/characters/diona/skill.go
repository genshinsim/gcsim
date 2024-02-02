package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillPressFrames []int
	skillHoldFrames  []int
)

const (
	skillPressHitmark = 5  // release
	skillHoldHitmark  = 29 // release
)

func init() {
	skillPressFrames = frames.InitAbilSlice(34) // Tap E -> E
	skillPressFrames[action.ActionAttack] = 33  // Tap E -> N1
	skillPressFrames[action.ActionBurst] = 33   // Tap E -> Q
	skillPressFrames[action.ActionDash] = 11    // Tap E -> D
	skillPressFrames[action.ActionJump] = 11    // Tap E -> J
	skillPressFrames[action.ActionSwap] = 16    // Tap E -> Swap

	skillHoldFrames = frames.InitAbilSlice(49) // Hold E -> E
	skillHoldFrames[action.ActionAttack] = 36  // Hold E -> N1
	skillHoldFrames[action.ActionBurst] = 37   // Hold E -> Q
	skillHoldFrames[action.ActionDash] = 31    // Hold E -> D
	skillHoldFrames[action.ActionJump] = 31    // Hold E -> J
	skillHoldFrames[action.ActionSwap] = 23    // Hold E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	if p["hold"] == 1 {
		return c.skillHold(travel)
	}
	return c.skillPress(travel)
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
		if c.Core.Rand.Float64() < 0.8 {
			c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
		}
	}
}

func (c *char) skillPress(travel int) (action.Info, error) {
	c.pawsPewPew(skillPressHitmark, travel, 2)
	c.SetCDWithDelay(action.ActionSkill, 360, skillPressHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillHold(travel int) (action.Info, error) {
	c.pawsPewPew(skillHoldHitmark, travel, 5)
	c.SetCDWithDelay(action.ActionSkill, 900, skillHoldHitmark)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) pawsPewPew(f, travel, pawCount int) {
	bonus := 1.0
	if pawCount == 5 {
		bonus = 1.75 // bonus if firing off 5
	}
	shdHp := (pawShieldPer[c.TalentLvlSkill()]*c.MaxHP() + pawShieldFlat[c.TalentLvlSkill()]) * bonus
	if c.Base.Cons >= 2 {
		shdHp *= 1.15
	}
	// call back to generate shield on hit
	// note that each paw should only be able to trigger callback once (if hit multi target)
	// and that subsequent shield generation should increase duation only
	//TODO: need to look into maybe additional paw hits actually create "new" shields?
	pawCB := func(done bool) combat.AttackCBFunc {
		return func(_ combat.AttackCB) {
			if done {
				return
			}
			// make sure this is only triggered once
			done = true

			// check if shield already exists, if so then just update duration
			exist := c.Core.Player.Shields.Get(shield.DionaSkill)
			var shd *shield.Tmpl
			if exist != nil {
				// update
				shd, _ = exist.(*shield.Tmpl)
				shd.Expires += pawDur[c.TalentLvlSkill()]
			} else {
				shd = &shield.Tmpl{
					ActorIndex: c.Index,
					Target:     -1,
					Src:        c.Core.F,
					ShieldType: shield.DionaSkill,
					Name:       "Diona Skill",
					HP:         shdHp,
					Ele:        attributes.Cryo,
					Expires:    c.Core.F + pawDur[c.TalentLvlSkill()], // 15 sec
				}
			}
			//TODO: check that this is actually properly extending duration
			c.Core.Player.Shields.Add(shd)
		}
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Icy Paw",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       paw[c.TalentLvlSkill()],
	}

	for i := 0; i < pawCount; i++ {
		done := false
		cb := pawCB(done)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.5,
			),
			0,
			travel+f-5+i,
			cb,
			c.makeParticleCB(),
		)
	}
}
