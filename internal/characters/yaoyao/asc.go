package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

const a4Status = "yaoyao-a4"

func (c *char) a1Ticker() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	switch c.Core.Player.CurrentState() {
	case action.DashState, action.JumpState:
		c.Core.Log.NewEvent("yaoyao a1 triggered", glog.LogCharacterEvent, c.Index).
			Write("state", c.Core.Player.CurrentState())
		c.a1Throw()
	}
	c.QueueCharTask(c.a1Ticker, 0.6*60)
}

func (c *char) a1Throw() {
	a1aoe := combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, skillTargetingRad)
	enemy := c.Core.Combat.RandomEnemyWithinArea(a1aoe, nil)
	if enemy == nil {
		return
	}
	target := enemy.Pos()

	radishExplodeAoE := combat.NewCircleHitOnTarget(target, nil, radishRad)

	c.QueueCharTask(func() {
		var hi player.HealInfo
		var ai combat.AttackInfo
		var snap combat.Snapshot

		if c.StatusIsActive(burstKey) {
			ai = c.burstRadishAI
			snap = c.Snapshot(&ai)
			hi = c.getBurstHealInfo(&snap)
		} else {
			ai = c.skillRadishAI
			snap = c.Snapshot(&ai)
			hi = c.getSkillHealInfo(&snap)
		}

		delay := 1
		c.Core.Tasks.Add(c.heal(radishExplodeAoE, hi), delay)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			radishExplodeAoE,
			delay,
			c.makeC2CB(),
		)
	}, travelDelay-1)
}

func (c *char) a4(index, src int) func() {
	return func() {
		if c.a4Srcs[index] != c.Core.F {
			return
		}

		char := c.Core.Player.ByIndex(index)
		if char.StatusIsActive(a4Status) {
			return
		}

		hi := player.HealInfo{
			Caller:  c.Index,
			Target:  index,
			Message: "Yaoyao A4",
			Src:     0.008 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		}
		c.Core.Player.Heal(hi)
		c.QueueCharTask(c.a4(index, src), 60)
	}
}
