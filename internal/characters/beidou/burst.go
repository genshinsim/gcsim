package beidou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

var burstFrames []int

const burstHitmark = 28

func init() {
	burstFrames = frames.InitAbilSlice(58)
	burstFrames[action.ActionAttack] = 55
	burstFrames[action.ActionDash] = 48
	burstFrames[action.ActionJump] = 48
	burstFrames[action.ActionSwap] = 46
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stormbreaker (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 100,
		Mult:       burstonhit[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.Core.Status.Add("beidouburst", 900)

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stormbreak Proc (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstproc[c.TalentLvlBurst()],
	}
	snap := c.Snapshot(&procAI)
	c.burstAtk = &combat.AttackEvent{
		Info:     procAI,
		Snapshot: snap,
	}

	if c.Base.Cons >= 1 {
		//create a shield
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.ShieldBeidouThunderShield,
			Name:       "Beidou C1",
			HP:         .16 * c.MaxHP(),
			Ele:        attributes.Electro,
			Expires:    c.Core.F + 900, //15 sec
		})
	}

	// apply after hitmark
	if c.Base.Cons >= 6 {
		c.Core.Tasks.Add(func() {
			for _, t := range c.Core.Combat.Targets() {
				e, ok := t.(core.Enemy)
				if !ok {
					continue
				}
				e.AddResistMod("beidouc6", 900-burstHitmark, attributes.Electro, -0.15)
			}
		}, burstHitmark)
	}

	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 1200)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		Post:            burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstProc() {
	icd := 0
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		//make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if c.Core.Status.Duration("beidouburst") == 0 {
			return false
		}
		if icd > c.Core.F {
			c.Core.Log.NewEvent("beidou Q (active) on icd", glog.LogCharacterEvent, c.Index)
			return false
		}

		//trigger a chain of attacks starting at the first target
		atk := *c.burstAtk
		atk.SourceFrame = c.Core.F
		atk.Pattern = combat.NewDefSingleTarget(t.Index(), combat.TargettableEnemy)
		cb := c.chain(c.Core.F, 1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, 1)

		c.Core.Log.NewEvent("beidou Q proc'd", glog.LogCharacterEvent, c.Index, "char", ae.Info.ActorIndex, "attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 60 // once per second
		return false
	}, "beidou-burst")
}

func (c *char) chain(src int, count int) combat.AttackCBFunc {
	if c.Base.Cons >= 2 && count == 5 {
		return nil
	}
	if c.Base.Cons <= 1 && count == 3 {
		return nil
	}
	return func(a combat.AttackCB) {
		//on hit figure out the next target
		trgs := c.Core.Combat.EnemyExcl(a.Target.Index())
		if len(trgs) == 0 {
			//do nothing if no other target other than this one
			return
		}
		//otherwise pick a random one
		next := c.Core.Rand.Intn(len(trgs))
		//queue an attack vs next target
		atk := *c.burstAtk
		atk.SourceFrame = src
		atk.Pattern = combat.NewDefSingleTarget(trgs[next], combat.TargettableEnemy)
		cb := c.chain(src, count+1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, 1)

	}
}
