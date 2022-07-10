package jean

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstStart = 40

func init() {
	burstFrames = frames.InitAbilSlice(84)
	burstFrames[action.ActionAttack] = 83
	burstFrames[action.ActionSkill] = 83
	burstFrames[action.ActionDash] = 70
	burstFrames[action.ActionJump] = 70
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//p is the number of times enemy enters or exits the field
	enter := p["enter"]
	if enter < 1 {
		enter = 1
	}
	delay, ok := p["enter_delay"]
	if !ok {
		delay = 600 / enter
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dandelion Breeze",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&ai)

	//initial hit at 40f
	c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 40)

	ai.Abil = "Dandelion Breeze (In/Out)"
	ai.Mult = burstEnter[c.TalentLvlBurst()]
	//first enter is at frame 55
	for i := 0; i < enter; i++ {
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 55+i*delay)
	}

	c.Core.Status.Add("jeanq", 600+burstStart)

	if c.Base.Cons >= 4 {
		//add debuff to all target for ??? duration
		for _, t := range c.Core.Combat.Targets() {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				continue
			}
			//10 seconds + animation
			e.AddResistMod(enemy.ResistMod{
				Base:  modifier.NewBase("jeanc4", 600+burstStart),
				Ele:   attributes.Anemo,
				Value: -0.4,
			})
		}
	}

	//heal on cast
	hpplus := snap.Stats[attributes.Heal]
	atk := snap.BaseAtk*(1+snap.Stats[attributes.ATKP]) + snap.Stats[attributes.ATK]
	heal := burstInitialHealFlat[c.TalentLvlBurst()] + atk*burstInitialHealPer[c.TalentLvlBurst()]
	healDot := burstDotHealFlat[c.TalentLvlBurst()] + atk*burstDotHealPer[c.TalentLvlBurst()]

	c.Core.Tasks.Add(func() {
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Dandelion Breeze",
			Src:     heal,
			Bonus:   hpplus,
		})
	}, burstStart)

	self, ok := c.Core.Combat.Target(0).(*avatar.Player)
	if !ok {
		panic("target 0 should be Player but is not!!")
	}

	//attack self
	selfSwirl := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dandelion Breeze (Self Swirl)",
		Element:    attributes.Anemo,
		Durability: 25,
	}

	//duration is 10.5s, first tick start at frame 100, + 60 each
	for i := 100; i < 100+630; i += 60 {
		c.Core.Tasks.Add(func() {
			// c.Core.Log.NewEvent("jean q healing", glog.LogCharacterEvent, c.Index, "+heal", hpplus, "atk", atk, "heal amount", healDot)
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Dandelion Field",
				Src:     healDot,
				Bonus:   hpplus,
			})

			ae := combat.AttackEvent{
				Info:        selfSwirl,
				Pattern:     combat.NewDefSingleTarget(0, combat.TargettablePlayer),
				SourceFrame: c.Core.F,
			}
			c.Core.Log.NewEvent("jean self swirling", glog.LogCharacterEvent, c.Index)
			self.ReactWithSelf(&ae)
		}, i)
	}

	c.SetCDWithDelay(action.ActionBurst, 1200, 38)
	// handle energy delay and a4
	c.Core.Tasks.Add(func() {
		c.Energy = 16 //jean a4
	}, 41)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}
}
