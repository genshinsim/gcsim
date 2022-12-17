package wanderer

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

const a4Hitmark = 30

func (c *char) Dash(p map[string]int) action.ActionInfo {
	delay := c.checkForSkillEnd()

	f := 21

	ai := action.ActionInfo{
		Frames:          func(action.Action) int { return delay + f },
		AnimationLength: delay + f,
		CanQueueAfter:   delay + f,
		State:           action.DashState,
	}

	if !c.StatusIsActive(skillKey) {
		c.Core.Tasks.Add(func() {
			req := c.Core.Player.AbilStamCost(c.Index, action.ActionDash, p)
			c.Core.Player.Stam -= req
			//this really shouldn't happen??
			if c.Core.Player.Stam < 0 {
				c.Core.Player.Stam = 0
			}
			c.Core.Player.LastStamUse = delay + c.Core.F
			c.Core.Events.Emit(event.OnStamUse, action.DashState)
		}, delay+f-1)

		return ai
	}

	if c.a4Active {
		c.a4Active = false

		a4Mult := 0.35

		// TODO: Write as attack mod?
		if c.StatusIsActive("wanderer-c1-atkspd") {
			a4Mult = 0.6
		}

		a4Info := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Gales of Reverie",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       a4Mult,
		}

		for i := 0; i < 4; i++ {
			c.Core.QueueAttack(a4Info, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.5),
				delay+a4Hitmark, delay+a4Hitmark)
		}
	} else {
		// TODO: Check Point consumption
		c.skydwellerPoints -= 20
	}

	return ai
}
