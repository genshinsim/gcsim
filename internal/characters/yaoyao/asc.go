package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

func (c *char) a1ticker() {
	c.QueueCharTask(func() {
		if !c.StatusIsActive(burstKey) {
			return
		}
		switch c.Core.Player.CurrentState() {
		case action.JumpState:
			fallthrough
		case action.DashState:
			c.Core.Log.NewEvent("yaoyao a1 triggered", glog.LogCharacterEvent, c.Index).
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

	ai := c.burstRadishAI
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
}

func (c *char) a4() {
	if c.Core.Player.ActiveChar().StatusIsActive("yaoyao-a4") {
		return
	}
	c.Core.Player.ActiveChar().AddStatus("yaoyao-a4", 5*60, true)
	for i := 0; i < 5; i++ {
		c.Core.Player.ActiveChar().QueueCharTask(func() {
			heal := 0.008 * c.MaxHP()
			hi := player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "yaoyao-a4-tick",
				Src:     heal,
				Bonus:   c.Stat(attributes.Heal),
			}
			c.Core.Player.Heal(hi)
		}, i*60)
	}
}
