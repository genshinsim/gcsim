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
	hitmark = 92
)

func init() {
	burstFrames = frames.InitAbilSlice(110)
	burstFrames[action.ActionSkill] = 109
	burstFrames[action.ActionSwap] = 108
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
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 7), 0, 0)
	}, hitmark)

	c.a4OnBurst()
	c.SetCD(action.ActionBurst, 15*60)

	c.QueueCharTask(c.summonBirgitta, 110)

	c.ConsumeEnergy(5)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
