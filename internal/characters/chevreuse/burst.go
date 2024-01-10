package chevreuse

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var (
	burstFrames []int
)

func init() {
	burstFrames = frames.InitAbilSlice(59) // Q -> N1/Dash/Walk
	burstFrames[action.ActionSkill] = 60
	burstFrames[action.ActionJump] = 60
	burstFrames[action.ActionSwap] = 59
}

const (
	snapshotDelay         = 10
	damageDelay           = 10
	mineExplosionInterval = 60
)

func (c *char) Burst(p map[string]int) (action.Info, error) {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Explosive Grenade",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		PoiseDMG:   100,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	mineAi := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Secondary Explosive Shell",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagChevreuseBurstMines,
		ICDGroup:   attacks.ICDGroupChevreuseBurstMines,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstSecondary[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
		snapshotDelay,
		damageDelay,
	)

	shellNum := 8
	c.QueueCharTask(func() {
		for i := 0; i < shellNum; i++ {
			c.Core.QueueAttack(
				mineAi,
				// TODO: How to tackle mines?
				combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 6),
				damageDelay,   // random number
				snapshotDelay) // random number
		}
	}, mineExplosionInterval)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
