package ineffa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	hitmark = 111
)

func init() {
	burstFrames = frames.InitAbilSlice(127)
	burstFrames[action.ActionSkill] = 127
	burstFrames[action.ActionDash] = 105
	burstFrames[action.ActionJump] = 106
	burstFrames[action.ActionWalk] = 130
	burstFrames[action.ActionSwap] = 126
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Supreme Instruction: Cyclonic Exterminator",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), 0, 0, c.c2MakeCB())
	}, hitmark)

	c.a4OnBurst()
	c.c2OnBurst()
	c.SetCD(action.ActionBurst, 15*60)

	c.QueueCharTask(c.summonBirgitta, hitmark)

	c.ConsumeEnergy(4)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}
