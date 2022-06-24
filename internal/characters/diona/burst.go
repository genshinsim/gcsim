package diona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstStart = 49

func init() {
	burstFrames = frames.InitAbilSlice(burstStart)
}

func (c *char) Burst(p map[string]int) action.ActionInfo {

	//initial hit
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Signature Mix (Initial)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0, burstStart-10)

	ai.Abil = "Signature Mix (Tick)"
	ai.Mult = burstDot[c.TalentLvlBurst()]
	snap := c.Snapshot(&ai)
	hpplus := snap.Stats[attributes.Heal]
	maxhp := c.MaxHP()
	heal := burstHealPer[c.TalentLvlBurst()]*maxhp + burstHealFlat[c.TalentLvlBurst()]

	//ticks every 2s, first tick at t=1s, then t=3,5,7,9,11, lasts for 12.5
	for i := 0; i < 6; i++ {
		c.Core.Tasks.Add(func() {
			c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0)
			// c.Core.Log.NewEvent("diona healing", core.LogCharacterEvent, c.Index, "+heal", hpplus, "max hp", maxhp, "heal amount", heal)
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Drunken Mist",
				Src:     heal,
				Bonus:   hpplus,
			})
		}, 60+i*120)
	}

	//apparently lasts for 12.5
	c.Core.Status.Add("dionaburst", burstStart+750) //TODO not sure when field starts, is it at animation end? prob when it lands...

	//c1
	if c.Base.Cons >= 1 {
		//15 energy after ends, flat not affected by ER
		c.Core.Tasks.Add(func() {
			c.AddEnergy("diona-c1", 15)
		}, burstStart+750)
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.SetCDWithDelay(action.ActionBurst, 1200, burstStart)
	c.ConsumeEnergy(burstStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstStart,
		State:           action.BurstState,
	}
}
