package ganyu

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var aimedFrames []int

const aimedHitmark = 103

func (c *char) Aimed(p map[string]int) action.ActionInfo {
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
		ActorIndex:   c.Index,
		Abil:         "Frost Flake Arrow",
		AttackTag:    combat.AttackTagExtra,
		ICDTag:       combat.ICDTagNone,
		ICDGroup:     combat.ICDGroupDefault,
		StrikeType:   combat.StrikeTypePierce,
		Element:      attributes.Cryo,
		Durability:   25,
		Mult:         ffa[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}

	// TODO: not sure if this works as intended
	skip := 0
	if c.Core.Status.Duration("ganyuc6") > 0 {
		c.Core.Status.Delete("ganyuc6")
		c.Core.Log.NewEvent("ganyu c6 proc used", glog.LogCharacterEvent, c.Index, "char", c.Index)
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
			c.Core.Log.NewEvent("a1 adding crit rate", glog.LogCharacterEvent, c.Index, "old", old, "new", snap.Stats[attributes.CR], "expiry", c.a1Expiry)
		}

		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefSingleTarget(1, combat.TargettableEnemy), travel)

		ai.Abil = "Frost Flake Bloom"
		ai.Mult = ffb[c.TalentLvlAttack()]
		ai.HitWeakPoint = false
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(2, false, combat.TargettableEnemy), travel+bloom)

		// first shot/bloom do not benefit from a1
		c.a1Expiry = c.Core.F + 60*5
	}, aimedHitmark-skip)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return aimedFrames[next] - skip },
		AnimationLength: aimedFrames[action.InvalidAction] - skip,
		CanQueueAfter:   aimedHitmark - skip,
		Post:            aimedHitmark - skip,
		State:           action.AimState,
	}
}
