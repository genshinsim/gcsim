package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 24
	burstBuffKey = "yanfei-q"
)

func init() {
	burstFrames = frames.InitAbilSlice(58) // Q -> N1
	burstFrames[action.ActionCharge] = 47  // Q -> CA
	burstFrames[action.ActionSkill] = 55   // Q -> E
	burstFrames[action.ActionDash] = 33    // Q -> D
	burstFrames[action.ActionJump] = 32    // Q -> J
	burstFrames[action.ActionSwap] = 46    // Q -> Swap
}

// Burst - Deals burst damage and adds status for charge attack bonus
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// +1 is to make sure the scarlet seal grant works correctly on the last frame
	// TODO: Not 100% sure whether this adds a seal at the exact moment the burst ends or not
	c.AddStatus(burstBuffKey, 15*60+1, true)

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(burstBuffKey, 15*60),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagExtra {
				return c.burstBuff, true
			}
			return nil, false
		},
	})

	done := false
	addSeal := func(_ combat.AttackCB) {
		if done {
			return
		}
		// Create max seals on hit
		if c.sealCount < c.maxTags {
			c.sealCount = c.maxTags
		}
		c.AddStatus(sealBuffKey, 600, true)
		c.Core.Log.NewEvent("yanfei gained max seals", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.sealCount)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Done Deal",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeBlunt,
		Element:            attributes.Pyro,
		Durability:         50,
		Mult:               burst[c.TalentLvlBurst()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5),
		0,
		burstHitmark,
		addSeal,
	)

	c.Core.Tasks.Add(c.burstAddSealHook(), 60)

	c.c4()

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(5)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel
		State:           action.BurstState,
	}
}

// Recurring task to add seals every second while burst is up
func (c *char) burstAddSealHook() func() {
	return func() {
		if !c.StatusIsActive(burstBuffKey) {
			return
		}
		if c.sealCount < c.maxTags {
			c.sealCount++
		}
		c.AddStatus(sealBuffKey, 600, true)

		c.Core.Log.NewEvent("yanfei gained seal from burst", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.sealCount)

		c.Core.Tasks.Add(c.burstAddSealHook(), 60)
	}
}
