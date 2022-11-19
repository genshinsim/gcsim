package layla

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillFrames []int

const (
	skillEnergy = "skill-energy-icd"

	skillHitmark = 28
)

// TODO: FRAMES
func init() {
	skillFrames = frames.InitAbilSlice(53) // E -> N1
	skillFrames[action.ActionBurst] = 52   // E -> Q
	skillFrames[action.ActionDash] = 25    // E -> D
	skillFrames[action.ActionJump] = 26    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	travel, ok := p["star_travel"]
	if !ok {
		travel = 26
	}
	c.starTravel = travel

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Nights of Formal Focus",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 4.47), 0, skillHitmark)

	c.QueueCharTask(func() {
		// wet: apply self infusion for 0.1s
		player, ok := c.Core.Combat.Player().(*avatar.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		player.ApplySelfInfusion(attributes.Cryo, 25, 0.1*60)

		exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
		if exist == nil {
			shield := shieldBase[c.TalentLvlSkill()] + shieldPer[c.TalentLvlSkill()]*c.MaxHP()
			c.Core.Player.Shields.Add(c.newShield(shield, 12*60))
			c.TickNightStar(false)
		} else {
			shd, _ := exist.(*shield.Tmpl)
			shd.Expires = c.Core.F + 12*60
		}
	}, skillHitmark)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}
}

func (c *char) StarsSkill() {
	c.Core.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
		if exist != nil {
			c.AddNightStars(2, true)
		}
		return false
	}, "stars-skill")
}
