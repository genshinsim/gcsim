package yanfei

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 65

func init() {
	burstFrames = frames.InitAbilSlice(65)
}

// Burst - Deals burst damage and adds status for charge attack bonus
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// +1 is to make sure the scarlet seal grant works correctly on the last frame
	// TODO: Not 100% sure whether this adds a seal at the exact moment the burst ends or not
	c.Core.Status.Add("yanfeiburst", 15*60+1)

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = burstBonus[c.TalentLvlBurst()]
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("yanfei-burst", 15*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == combat.AttackTagExtra {
				return m, true
			}
			return nil, false
		},
	})

	done := false
	addSeal := func(a combat.AttackCB) {
		if done {
			return
		}
		// Create max seals on hit
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"] = c.maxTags
		}
		c.sealExpiry = c.Core.F + 600
		c.Core.Log.NewEvent("yanfei gained max seals", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.Tags["seal"]).
			Write("expiry", c.sealExpiry)
		done = true
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Done Deal",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, burstHitmark, addSeal)

	c.Core.Tasks.Add(c.burstAddSealHook(), 60)

	c.c4()

	c.SetCDWithDelay(action.ActionBurst, 20*60, 8)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

// Recurring task to add seals every second while burst is up
func (c *char) burstAddSealHook() func() {
	return func() {
		if c.Core.Status.Duration("yanfeiburst") == 0 {
			return
		}
		if c.Tags["seal"] < c.maxTags {
			c.Tags["seal"]++
		}
		c.sealExpiry = c.Core.F + 600

		c.Core.Log.NewEvent("yanfei gained seal from burst", glog.LogCharacterEvent, c.Index).
			Write("current_seals", c.Tags["seal"]).
			Write("expiry", c.sealExpiry)

		c.Core.Tasks.Add(c.burstAddSealHook(), 60)
	}
}
