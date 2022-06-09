package noelle

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

// TODO: not sure about this
const burstStart = 80

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// TODO: Assume snapshot happens immediately upon cast since the conversion buffs the two burst hits
	// Generate a "fake" snapshot in order to show a listing of the applied mods in the debug
	aiSnapshot := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Stat Snapshot)",
	}
	snapshot := c.Snapshot(&aiSnapshot)
	burstDefSnapshot := snapshot.BaseDef*(1+snapshot.Stats[attributes.DEFP]) + snapshot.Stats[attributes.DEF]
	mult := defconv[c.TalentLvlBurst()]
	if c.Base.Cons >= 6 {
		mult += 0.5
	}
	// Add mod for def to attack burst conversion
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATK] = mult * burstDefSnapshot

	// TODO: Confirm exact timing of buff - for now matched to status duration previously set, which is 900 + animation frames
	c.AddStatMod("noelle-burst", 900+burstStart, attributes.ATK, func() ([]float64, bool) {
		return m, true
	})
	c.Core.Log.NewEvent("noelle burst", glog.LogSnapshotEvent, c.Index, "total def", burstDefSnapshot, "atk added", m[attributes.ATK], "mult", mult)

	c.burstInfusion(900 + burstStart)
	c.Core.Status.Add("noelleq", 900+burstStart)
	// Queue up task for Noelle burst extension
	// https://library.keqingmains.com/evidence/characters/geo/noelle#noelle-c6-burst-extension
	if c.Base.Cons >= 6 {
		c.Core.Tasks.Add(func() {
			if c.Core.Player.Active() == c.Index {
				return
			}
			// Adding the mod again with the same key replaces it
			c.AddStatMod("noelle-burst", 600, attributes.ATK, func() ([]float64, bool) {
				return m, true
			})
			c.Core.Log.NewEvent("noelle max burst extension activated", glog.LogCharacterEvent, c.Index, "new_expiry", c.Core.F+600)
			// check if this work as intended
			c.burstInfusion(600)
			c.Core.Status.Add("noelleq", 600)
		}, 900+burstStart)
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sweeping Time (Burst)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(6.5, false, combat.TargettableEnemy), 24, 24)

	ai.Abil = "Sweeping Time (Skill)"
	ai.Mult = burstskill[c.TalentLvlBurst()]
	c.Core.QueueAttack(ai, combat.NewDefCircHit(4.5, false, combat.TargettableEnemy), 65, 65)

	c.SetCD(action.ActionBurst, 900)
	c.ConsumeEnergy(8)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		Post:            burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstInfusion(dur int) {
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"noelle-burst",
		attributes.Geo,
		dur,
		false,
		combat.AttackTagNormal, combat.AttackTagPlunge, combat.AttackTagExtra,
	)
}
