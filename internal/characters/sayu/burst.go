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

func init() {
	burstFrames = frames.InitAbilSlice(65)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	// dmg
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Yoohoo Art: Mujina Flurry",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 16)

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

	for i := 150; i < 150+540; i += 90 {
		c.Core.Tasks.Add(func() {
			active := c.Core.Player.ActiveChar()
			//this is going to be a bit slow..
			enemies := c.Core.Combat.EnemyByDistance(0, 0, 7) //TODO: no idea what the range of this check is
			needHeal := len(enemies) == 0 || active.HPCurrent/active.MaxHP() <= .7
			needAttack := !needHeal
			if c.Base.Cons >= 1 {
				needHeal = true
				needAttack = true
			}
			if needHeal {
				c.Core.Player.Heal(player.HealInfo{
					Caller:  c.Index,
					Target:  c.Core.Player.Active(),
					Message: "Muji-Muji Daruma",
					Src:     heal,
					Bonus:   d.Snapshot.Stats[attributes.Heal],
				})
			}
			if needAttack {
				c.Core.QueueAttackEvent(d, 0)
			}
		}, i)
	}

	c.SetCDWithDelay(action.ActionBurst, 20*60, 11)
	c.ConsumeEnergy(11)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.InvalidAction],
		State:           action.BurstState,
	}
}

func (c *char) createBurstSnapshot() *combat.AttackEvent {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Muji-Muji Daruma",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstSkill[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	return (&combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewDefCircHit(5, false, combat.TargettableEnemy), // including A4
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})
}
