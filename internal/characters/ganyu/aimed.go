package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames [][]int

var aimedHitmarks = []int{15 - 12, 15, 74, 103}

func init() {
	aimedFrames = make([][]int, 4)

	// Aimed Shot (ARCC)
	aimedFrames[0] = frames.InitAbilSlice(25 - 12)
	aimedFrames[0][action.ActionDash] = aimedHitmarks[0]
	aimedFrames[0][action.ActionJump] = aimedHitmarks[0]

	// Aimed Shot
	aimedFrames[1] = frames.InitAbilSlice(25)
	aimedFrames[1][action.ActionDash] = aimedHitmarks[1]
	aimedFrames[1][action.ActionJump] = aimedHitmarks[1]

	// Fully-Charged Aimed Shot Lv. 1 (Fully-Charged Aimed Shot)
	aimedFrames[2] = frames.InitAbilSlice(85)
	aimedFrames[2][action.ActionDash] = aimedHitmarks[2]
	aimedFrames[2][action.ActionJump] = aimedHitmarks[2]

	// Fully-Charged Aimed Shot Lv. 2 (Frostflake Arrow + Frostflake Arrow Bloom)
	aimedFrames[3] = frames.InitAbilSlice(113)
	aimedFrames[3][action.ActionDash] = aimedHitmarks[3]
	aimedFrames[3][action.ActionJump] = aimedHitmarks[3]
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	hold, ok := p["hold"]
	if !ok || hold < 0 {
		hold = 3
	}
	if hold > 3 {
		hold = 3
	}
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	bloom, ok := p["bloom"]
	if !ok {
		bloom = 24
	}
	weakspot, ok := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Fully-Charged Aimed Shot",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Cryo,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	if hold < 2 {
		ai.Abil = "Aimed Shot"
		if hold == 0 {
			ai.Abil += " (ARCC)"
		}
		ai.Element = attributes.Physical
		ai.Mult = aim[c.TalentLvlAttack()]
	}

	// TODO: not sure if this works as intended
	skip := 0
	if c.Core.Status.Duration(c6Key) > 0 && hold == 3 {
		c.Core.Status.Delete(c6Key)
		c.Core.Log.NewEvent(c6Key+" proc used", glog.LogCharacterEvent, c.Index).
			Write("char", c.Index)
		// skip aimed charge time
		skip = 83
	}

	if hold == 3 {
		// snapshot delay and handles A1
		// A1 doesn't apply to the first aimed shot
		c.Core.Tasks.Add(func() {
			// make sure Frostflake Arrow and Bloom have the correct values
			ai.Abil = "Frostflake Arrow"
			ai.Element = attributes.Cryo
			ai.Mult = ffa[c.TalentLvlAttack()]

			snap := c.Snapshot(&ai)
			if c.Core.F < c.a1Expiry {
				old := snap.Stats[attributes.CR]
				snap.Stats[attributes.CR] += .20
				c.Core.Log.NewEvent("a1 adding crit rate", glog.LogCharacterEvent, c.Index).
					Write("old", old).
					Write("new", snap.Stats[attributes.CR]).
					Write("expiry", c.a1Expiry)
			}

			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), travel)

			ai.Abil = "Frostflake Arrow Bloom"
			ai.Mult = ffb[c.TalentLvlAttack()]
			ai.HitWeakPoint = false
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), travel+bloom)

			// first shot/bloom do not benefit from a1
			c.a1Expiry = c.Core.F + 60*5
		}, aimedHitmarks[hold]-skip)
	} else {
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), aimedHitmarks[hold], aimedHitmarks[hold]+travel)
	}

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[hold][next] - skip },
		AnimationLength: aimedFrames[hold][action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmarks[hold] - skip,
		State:           action.AimState,
	}
}
