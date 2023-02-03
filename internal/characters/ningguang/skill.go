package ningguang

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 17

func init() {
	skillFrames = frames.InitAbilSlice(62)
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 29
	skillFrames[action.ActionWalk] = 53
	skillFrames[action.ActionSwap] = 60
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Jade Screen",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Geo,
		Durability:         25,
		Mult:               skill[c.TalentLvlSkill()],
		HitlagHaltFrames:   0.05 * 60,
		HitlagFactor:       0.05,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	player := c.Core.Combat.Player()
	screenDir := player.Direction()
	screenPos := combat.CalcOffsetPoint(player.Pos(), combat.Point{Y: 3}, player.Direction())

	c.Core.Tasks.Add(func() {
		c.skillSnapshot = c.Snapshot(&ai)
		c.Core.QueueAttackWithSnap(
			ai,
			c.skillSnapshot,
			combat.NewCircleHitOnTarget(screenPos, nil, 5),
			0,
		)
	}, skillHitmark)

	//put skill on cd first then check for construct/c2
	c.SetCD(action.ActionSkill, 720)

	//create a construct
	c.Core.Constructs.New(c.newScreen(1800, screenDir, screenPos), true) //30 seconds

	c.lastScreen = c.Core.F

	noscreen, ok := p["noscreen"]
	if !ok && noscreen != 0 {
		c.Core.Tasks.Add(func() {
			c.Core.Constructs.Destroy(c.lastScreen)
		}, 1)
	}

	if !c.StatusIsActive(skillParticleICDKey) {
		//3 balls, 33% chance of a fourth
		var count float64 = 3
		if c.Core.Rand.Float64() < .33 {
			count = 4
		}
		c.Core.QueueParticle("ningguang", count, attributes.Geo, skillHitmark+c.ParticleDelay)
		c.AddStatus(skillParticleICDKey, 360, true)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}
}
