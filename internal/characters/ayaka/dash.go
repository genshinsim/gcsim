package ayaka

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var dashFrames []int

const dashHitmark = 20

func init() {
	dashFrames = frames.InitAbilSlice(35)
	dashFrames[action.ActionDash] = 30
	dashFrames[action.ActionSwap] = 34
}

// TODO: move this into PostDash event instead
func (c *char) Dash(p map[string]int) action.ActionInfo {
	f, ok := p["f"]
	if !ok {
		f = 0
	}

	//no dmg attack at end of dash
	ai := combat.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagDash,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Cryo,
		Durability: 25,
	}

	//restore on hit, once per attack
	//a4 increase cryo dmg by 18% for 10s
	m := make([]float64, attributes.EndStatType)
	m[attributes.CryoP] = 0.18
	once := false
	cb := func(a combat.AttackCB) {
		if once {
			return
		}
		once = true

		c.Core.Player.RestoreStam(10)
		c.AddStatMod("ayaka-a4", 600, attributes.CryoP, func() ([]float64, bool) {
			return m, true
		})
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), dashHitmark+f, dashHitmark+f, cb)

	//add cryo infuse
	// TODO: weapon infuse happen at dash end
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"ayaka-dash",
		attributes.Cryo,
		300,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)

	// call default implementation to handle stamina
	c.Character.Dash(p)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return dashFrames[next] + f },
		AnimationLength: dashFrames[action.InvalidAction] + f,
		CanQueueAfter:   dashHitmark + f,
		State:           action.DashState,
	}
}
