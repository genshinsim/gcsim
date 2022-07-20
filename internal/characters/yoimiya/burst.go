package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

var burstFrames []int

const burstHitmark = 75

func init() {
	burstFrames = frames.InitAbilSlice(114)
	burstFrames[action.ActionSkill] = 110
	burstFrames[action.ActionDash] = 111
	burstFrames[action.ActionJump] = 113
	burstFrames[action.ActionSwap] = 109
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//assume it does skill dmg at end of it's animation
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aurous Blaze",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(5, false, combat.TargettableEnemy), 0, burstHitmark)

	//marker an opponent after first hit
	//ignore the bouncing around for now (just assume it's always target 0)
	//icd of 2s, removed if down
	duration := 600
	if c.Base.Cons >= 1 {
		duration = 840
	}
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add("aurous", duration)
		//attack buff
		c.a4()
	}, burstHitmark)

	//add cooldown to sim
	c.SetCD(action.ActionBurst, 15*60)
	//use up energy
	c.ConsumeEnergy(5)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstHook() {
	//check on attack landed for target 0
	//if aurous active then trigger dmg if not on cd
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		if c.Core.Status.Duration("aurous") == 0 {
			return false
		}
		if ae.Info.ActorIndex == c.Index {
			//ignore for self
			return false
		}
		//ignore if on icd
		if c.Core.Status.Duration("aurousicd") > 0 {
			return false
		}
		//ignore if wrong tags
		switch ae.Info.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
		case combat.AttackTagPlunge:
		case combat.AttackTagElementalArt:
		case combat.AttackTagElementalBurst:
		default:
			return false
		}
		//do explosion, set icd
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Aurous Blaze (Explode)",
			AttackTag:  combat.AttackTagElementalBurst,
			ICDTag:     combat.ICDTagElementalBurst,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       burstExplode[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(ai, combat.NewDefCircHit(3, false, combat.TargettableEnemy), 0, 1)

		c.Core.Status.Add("aurousicd", 120) //2 sec icd

		//check for c4

		if c.Base.Cons >= 4 {
			c.ReduceActionCooldown(action.ActionSkill, 72)
		}

		return false

	}, "yoimiya-burst-check")

	if c.Core.Flags.DamageMode {
		//add check for if yoimiya dies
		c.Core.Events.Subscribe(event.OnCharacterHurt, func(_ ...interface{}) bool {
			if c.HPCurrent <= 0 {
				c.Core.Status.Delete("aurous")
			}
			return false
		}, "yoimiya-died")
	}
}
