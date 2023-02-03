package sayu

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstHitmark = 12
const tickTaskDelay = 20

func init() {
	burstFrames = frames.InitAbilSlice(65) // Q -> N1/E/J
	burstFrames[action.ActionDash] = 64    // Q -> D
	burstFrames[action.ActionSwap] = 64    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// dmg
	ai := combat.AttackInfo{
		ActorIndex:       c.Index,
		Abil:             "Yoohoo Art: Mujina Flurry",
		AttackTag:        combat.AttackTagElementalBurst,
		ICDTag:           combat.ICDTagNone,
		ICDGroup:         combat.ICDGroupDefault,
		StrikeType:       combat.StrikeTypeDefault,
		Element:          attributes.Anemo,
		Durability:       25,
		Mult:             burst[c.TalentLvlBurst()],
		HitlagFactor:     0.05,
		HitlagHaltFrames: 0.02 * 60,
	}
	snap := c.Snapshot(&ai)
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 1.5}, 10)
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		combat.NewCircleHitOnTarget(burstArea.Shape.Pos(), nil, 4.5),
		burstHitmark,
	)

	// heal
	atk := snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK]
	heal := initHealFlat[c.TalentLvlBurst()] + atk*initHealPP[c.TalentLvlBurst()]
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Yoohoo Art: Mujina Flurry",
		Src:     heal,
		Bonus:   snap.Stats[attributes.Heal],
	})

	// ticks
	d := c.createBurstSnapshot()
	atk = d.Snapshot.BaseAtk*(1+d.Snapshot.Stats[attributes.ATKP]) + d.Snapshot.Stats[attributes.ATK]
	heal = burstHealFlat[c.TalentLvlBurst()] + atk*burstHealPP[c.TalentLvlBurst()]

	if c.Base.Cons >= 6 {
		// TODO: is it snapshots?
		d.Info.FlatDmg += atk * math.Min(d.Snapshot.Stats[attributes.EM]*0.002, 4.0)
		heal += math.Min(d.Snapshot.Stats[attributes.EM]*3, 6000)
	}

	// make sure that this task gets executed:
	// - after q initial hit hitlag happened
	// - before sayu can get affected by any more hitlag
	c.QueueCharTask(func() {
		// first tick is at 145
		for i := 145 - tickTaskDelay; i < 145+540-tickTaskDelay; i += 90 {
			c.Core.Tasks.Add(func() {
				// check for player
				// only check HP of active character
				char := c.Core.Player.ActiveChar()
				hasC1 := c.Base.Cons >= 1
				// C1 ignores HP limit
				needHeal := c.Core.Combat.Player().IsWithinArea(burstArea) && (char.HPCurrent/char.MaxHP() <= 0.7 || hasC1)

				// check for enemy
				enemy := c.Core.Combat.ClosestEnemyWithinArea(burstArea, nil)

				// determine whether to attack or heal
				// - C1 makes burst able to attack both an enemy and heal the player at the same time
				needAttack := enemy != nil && (!needHeal || hasC1)
				if needAttack {
					d.Pattern = combat.NewCircleHitOnTarget(enemy, nil, 3.5) // including A4
					c.Core.QueueAttackEvent(d, 0)
				}
				if needHeal {
					c.Core.Player.Heal(player.HealInfo{
						Caller:  c.Index,
						Target:  char.Index,
						Message: "Muji-Muji Daruma",
						Src:     heal,
						Bonus:   d.Snapshot.Stats[attributes.Heal],
					})
				}
			}, i)
		}
	}, tickTaskDelay)

	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}

// TODO: is this helper function needed?
func (c *char) createBurstSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Muji-Muji Daruma",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstSkill[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	return (&combat.AttackEvent{
		Info:        ai,
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
}
