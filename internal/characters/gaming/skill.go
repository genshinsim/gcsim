package gaming

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int
var lionWalkBack int

const (
	skillHitmark  = 20
	gamingICD     = 3   // seconds
	gamingSkillCD = 360 // frames
)

func init() {
	// E copied from kazuha
	skillFrames = frames.InitAbilSlice(77)
	skillFrames[action.ActionHighPlunge] = 24
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if p["ticks"] > 0 {
		lionWalkBack = p["ticks"]
	} else {
		lionWalkBack = 120
	}

	ai := combat.AttackInfo{
		Abil:       "Bestial Ascent (E)",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Physical,
		Durability: 25,
		Mult:       0, // initial skill hits but deals 0 damage
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)

	c.Core.QueueAttack(ai, ap, 0, skillHitmark, c.onPounceHit)
	c.SetCDWithDelay(action.ActionSkill, gamingSkillCD, 33) // delay unknown

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionHighPlunge], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) onPounceHit(a combat.AttackCB) {
	c.spawnManchai(a)
}

func (c *char) spawnManchai(_ combat.AttackCB) {
	if c.StatusIsActive(burstKey) && c.CurrentHPRatio() > 0.5 {
		if !c.StatusIsActive(lionKey) {
			c.AddStatus(lionKey, lionWalkBack, false)
			c.QueueCharTask(func() {
				c.ResetActionCooldown(action.ActionSkill)
				c.DeleteStatus(lionKey)
				c.c1()
			}, lionWalkBack)
		}
	}
}
