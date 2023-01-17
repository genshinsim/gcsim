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
	withinBurstKey = "nahida-q-within"
	burstKey       = "nahida-q"
)

func (c *char) Burst(p map[string]int) action.ActionInfo {
	var dur float64 = 15
	if c.hydroCount > 0 {
		dur += burstTriKarmaDurationExtend[c.hydroCount-1][c.TalentLvlBurst()]
	}
	f := int(dur * 60)

	withinTimer := 60
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1}, 20)
	c.Core.Tasks.Add(func() {
		c.burstSrc = c.Core.F
		src := c.Core.F
		c.Core.Status.Add(burstKey, f)
		// a1 buff is calculated at the start of burst
		c.calcA1Buff()
		for i := 30; i <= f; i += 30 {
			c.Core.Tasks.Add(func() {
				// don't tick if another burst has already started
				if src != c.burstSrc {
					return
				}
				// don't apply anything if outside of burst area
				if !c.Core.Combat.Player().IsWithinArea(burstArea) {
					return
				}

				c.AddStatus(withinBurstKey, withinTimer, true)
				if c.pyroCount > 0 {
					c.AddAttackMod(character.AttackMod{
						Base: modifier.NewBaseWithHitlag(burstKey, 60),
						Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
							return c.pyroBurstBuff, atk.Info.Abil == "Tri-Karma Purification"
						},
					})
				}
				c.applyA1(withinTimer)
			}, i)
		}
		if c.Base.Cons >= 6 {
			//lasts 10s
			//TODO: should this be delayed until animation end?
			c.AddStatus(c6ActiveKey, 600, true)
			c.c6Count = 0
			c.DeleteStatus(c6ICDKey) //TODO: check if this resets icd?
		}
	}, 66)

	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 810)
	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}
