package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames []int

const aimedHitmark = 103

func init() {
	// TODO: get separate counts for each cancel, currently using generic frames for all of them
	aimedFrames = frames.InitAbilSlice(113)
	aimedFrames[action.ActionDash] = aimedHitmark
	aimedFrames[action.ActionJump] = aimedHitmark
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	bloom, ok := p["bloom"]
	if !ok {
		bloom = 24
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Frost Flake Arrow",
		AttackTag:            combat.AttackTagExtra,
		ICDTag:               combat.ICDTagNone,
		ICDGroup:             combat.ICDGroupDefault,
		StrikeType:           combat.StrikeTypePierce,
		Element:              attributes.Cryo,
		Durability:           25,
		Mult:                 ffa[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	// TODO: not sure if this works as intended
	skip := 0
	if c.Core.Status.Duration(c6Key) > 0 {
		c.Core.Status.Delete(c6Key)
		c.Core.Log.NewEvent(c6Key+" proc used", glog.LogCharacterEvent, c.Index).
			Write("char", c.Index)
		// skip aimed charge time
		skip = 83
	}

	// snapshot delay and handles A1
	// A1 doesn't apply to the first aimed shot
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		if c.Core.F < c.a1Expiry {
			old := snap.Stats[attributes.CR]
			snap.Stats[attributes.CR] += .20
			c.Core.Log.NewEvent("a1 adding crit rate", glog.LogCharacterEvent, c.Index).
				Write("old", old).
				Write("new", snap.Stats[attributes.CR]).
				Write("expiry", c.a1Expiry)

		}

		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				combat.Point{Y: -0.5},
				0.1,
				1,
			),
			travel,
			c.c1(),
		)

		ai.Abil = "Frost Flake Bloom"
		ai.Mult = ffb[c.TalentLvlAttack()]
		ai.HitWeakPoint = false
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5),
			travel+bloom,
			c.c1(),
		)
		
		// first shot/bloom do not benefit from a1
		c.a1Expiry = c.Core.F + 60*5
	}, aimedHitmark-skip)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[next] - skip },
		AnimationLength: aimedFrames[action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmark - skip,
		State:           action.AimState,
	}
}
