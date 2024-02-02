package layla

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var skillFrames []int

const (
	skillHitmark    = 32
	particleICD1Key = "layla-particle-icd-1"
	particleICD2Key = "layla-particle-icd-2"
)

func init() {
	skillFrames = frames.InitAbilSlice(43) // E -> Q/D/J
	skillFrames[action.ActionAttack] = 41  // E -> N1
	skillFrames[action.ActionSkill] = 42   // E -> E
	skillFrames[action.ActionWalk] = 42    // E -> W
	skillFrames[action.ActionSwap] = 41    // E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Nights of Formal Focus",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), 0, skillHitmark)

	c.QueueCharTask(func() {
		// add shield
		exist := c.Core.Player.Shields.Get(shield.LaylaSkill)
		if exist == nil {
			shield := shieldBase[c.TalentLvlSkill()] + shieldPer[c.TalentLvlSkill()]*c.MaxHP()
			if c.Base.Cons >= 1 {
				shield *= 1.2
			}
			c.Core.Player.Shields.Add(c.newShield(shield, 12*60))
		} else {
			shd, _ := exist.(*shd)
			shd.Expires = c.Core.F + 12*60
		}

		// apply cryo & run a task
		player, ok := c.Core.Combat.Player().(*avatar.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		player.ApplySelfInfusion(attributes.Cryo, 25, 0.1*60)

		c.starTickSrc = c.Core.F
		c.tickNightStar(c.starTickSrc, false)()
	}, 19)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 19)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionAttack], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) starsSkill() {
	c.Core.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		exist := c.Core.Player.Shields.Get(shield.LaylaSkill)
		if exist != nil {
			c.addNightStars(2, ICDNightStarSkill)
		}
		return false
	}, "stars-skill")
}
