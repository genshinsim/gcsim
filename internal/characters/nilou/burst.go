package nilou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const (
	lingeringAeonStatus = "lingering_aeon"

	burstHitmark     = 48
	burstAeonHitmark = 3 * 60
)

// TODO: cancel frames & hitlags
func init() {
	burstFrames = frames.InitAbilSlice(70)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dance of Abzendegi: Distant Dreams, Listening Spring",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy),
		burstHitmark,
		burstHitmark,
		c.LingeringAeon,
	)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 900) // 15s * 60

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) LingeringAeon(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddStatus(lingeringAeonStatus, burstAeonHitmark, false)

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lingering Aeon",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    c.MaxHP() * burstAeon[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), 0.1, false, combat.TargettableEnemy),
		burstAeonHitmark,
		burstAeonHitmark,
	)
}
