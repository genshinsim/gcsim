package nahida

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(112)
	burstFrames[action.ActionSwap] = 111
}

const (
	burstKey = "nahida-q"
)

func (c *char) Burst(p map[string]int) action.ActionInfo {
	var dur float64 = 15
	if c.hydroCount > 0 {
		dur += burstTriKarmaDurationExtend[c.hydroCount-1][c.TalentLvlBurst()]
	}
	f := int(dur * 60)

	//TODO: gadget shouldn't be affected by hitlag
	//TODO: consider using an actual gadget here and use collision to detect if "in range"
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add(burstKey, f)
	}, 66)

	if c.pyroCount > 0 {
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(burstKey, f),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				return c.pyroBurstBuff, atk.Info.Abil == "Tri-Karma Purification"
			},
		})
	}

	c.a1(f)

	if c.Base.Cons > 5 {
		//lasts 10s
		//TODO: should this be delayed until animation end?
		c.Core.Tasks.Add(func() {
			c.AddStatus(c6ActiveKey, 600, true)
			c.c6count = 0
			c.DeleteStatus(c6ICDKey) //TODO: check if this resets icd?
		}, 66)
	}

	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 810)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
