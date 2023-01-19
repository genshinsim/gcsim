package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const a1ICDKey = "yaoyao-a1-icd"

func (c *char) a1hook() {
	c.Core.Events.Subscribe(event.OnStateChange, func(args ...interface{}) bool {
		//check if buff is up
		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(a1ICDKey) {
			return false
		}
		next := args[1].(action.AnimationState)
		switch next {
		case action.DashState:
			fallthrough
		case action.JumpState:
			c.Core.Log.NewEvent("yaoyao a1 triggered from state change", glog.LogCharacterEvent, c.Index).
				Write("state", next)
			c.a1Throw()
		}
		return false
	}, "yaoayo-a1")
}

func (c *char) a1ticker(src int) {
	c.a1src = src
	c.QueueCharTask(func() {
		if !c.StatusIsActive(burstKey) {
			return
		}
		if c.a1src != src {
			return
		}
		switch c.Core.Player.CurrentState() {
		case action.JumpState:
			fallthrough
		case action.DashState:
			c.Core.Log.NewEvent("yaoyao a1 triggered from ticker", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("state", c.Core.Player.CurrentState())
			c.a1Throw()
		}
	}, 0.6*60)

}
func (c *char) a1Throw() {

	a1aoe := combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 7)
	enemy := c.Core.Combat.RandomEnemyWithinArea(a1aoe, nil)
	if enemy == nil {
		return
	}
	target := enemy.Pos()

	radishExplodeAoE := combat.NewCircleHitOnTarget(target, nil, radishRad)

	ai := c.burstAI
	hi := c.getBurstHealInfo()

	c.Core.QueueAttack(
		ai,
		radishExplodeAoE,
		travelDelay,
		travelDelay,
	)
	if c.Core.Combat.Player().IsWithinArea(radishExplodeAoE) {
		c.radishHeal(hi)
	}
	c.AddStatus(a1ICDKey, 0.6*60, true)
	c.a1ticker(c.Core.F)
}

func (c *char) a4() {
	// fuck this shit I'll do it later it's just some healing anyways
}
