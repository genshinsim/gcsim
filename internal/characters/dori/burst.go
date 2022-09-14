package dori

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const burstStart = 55

func init() {
	burstFrames = frames.InitAbilSlice(63) // Q -> D/J
	burstFrames[action.ActionAttack] = 62  // Q -> N1
	burstFrames[action.ActionSkill] = 62   // Q -> E
	burstFrames[action.ActionSwap] = 62    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Alcazarzaray's Exactitude: Connector DMG",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    c.MaxHP() * burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai) //TODO: does it snaps?

	//damage ticks
	for i := 1; i < 30; i++ { //12/0.4s=30 ticks
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefBoxHit(1, -2, false, combat.TargettableEnemy), 24*i+burstStart) //TODO: put an actual accurate hitbox that depends on both the player and lamps position....yeah it wont be me lmao

	}
	interval := 12 * 60 / 6 //interval between heals

	for i := burstStart; i < 6*interval+burstStart; i += interval { //TODO: it does regen and heal at the same ticks?
		c.Core.Tasks.Add(func() {
			//Heals
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Alcazarzaray's Exactitude: Healing",
				Src:     bursthealpp[c.TalentLvlBurst()]*c.MaxHP() + bursthealflat[c.TalentLvlBurst()],
				Bonus:   snap.Stats[attributes.Heal],
			})
			//Energy regen to active char
			active := c.Core.Player.ActiveChar()
			active.AddEnergy("Alcazarzaray's Exactitude: Energy regen", burstenergy[c.TalentLvlBurst()])
		}, i)
	}
	c.Core.Tasks.Add(func() {
		// C4
		if c.Base.Cons >= 4 {
			c.c4()
		}

	}, burstStart)

	c.ConsumeEnergy(4)
	c.SetCD(action.ActionBurst, 1200) // 20s * 60

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}
