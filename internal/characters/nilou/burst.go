package nilou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const (
	lingeringAeonStatus = "lingeringaeon"

	burstHitmark     = 91
	burstAeonHitmark = 121
)

func init() {
	burstFrames = frames.InitAbilSlice(110) // Q -> Dash
	burstFrames[action.ActionAttack] = 108
	burstFrames[action.ActionSkill] = 108
	burstFrames[action.ActionJump] = 109
	burstFrames[action.ActionSwap] = 107
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dance of Abzendegi: Distant Dreams, Listening Spring",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeSlash,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		burstHitmark,
		burstHitmark,
		c.LingeringAeon,
	)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 18*60)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) LingeringAeon(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	t.AddStatus(lingeringAeonStatus, burstAeonHitmark, false)

	t.QueueEnemyTask(func() {
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Lingering Aeon",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			FlatDmg:    c.MaxHP() * burstAeon[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewSingleTargetHit(t.Key()),
			0,
			0,
		)
	}, burstAeonHitmark)
}
