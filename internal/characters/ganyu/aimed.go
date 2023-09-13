package ganyu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
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

func (c *char) Aimed(p map[string]int) action.Info {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Frost Flake Arrow",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Cryo,
		Durability:           25,
		Mult:                 ffa[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     .12 * 60,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c1cb := c.c1()
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
	c.Core.Tasks.Add(func() {
		snap := c.Snapshot(&ai)
		// A1:
		// After firing a Frostflake Arrow, the CRIT Rate of subsequent Frostflake Arrows
		// and their resulting bloom effects is increased by 20% for 5s.
		// - doesn't apply to the first aimed shot
		if c.Base.Ascension >= 1 && c.Core.F < c.a1Expiry {
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
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			travel,
			c1cb,
		)

		ai.Abil = "Frost Flake Bloom"
		ai.Mult = ffb[c.TalentLvlAttack()]
		ai.HitWeakPoint = false
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5),
			travel+18, // bloom always hits 18f after the arrow
			c1cb,
		)

		// first shot/bloom do not benefit from a1
		c.a1Expiry = c.Core.F + 60*5
	}, aimedHitmark-skip)

	return action.Info{
		Frames:          func(next action.Action) int { return aimedFrames[next] - skip },
		AnimationLength: aimedFrames[action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmark - skip,
		State:           action.AimState,
	}
}
