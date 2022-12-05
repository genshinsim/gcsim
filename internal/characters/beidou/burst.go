package beidou

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 28
	burstKey     = "beidouburst"
	burstICDKey  = "beidou-burst-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(58)
	burstFrames[action.ActionAttack] = 55
	burstFrames[action.ActionDash] = 48
	burstFrames[action.ActionJump] = 48
	burstFrames[action.ActionSwap] = 46
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Stormbreaker (Q)",
		AttackTag:          combat.AttackTagElementalBurst,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeDefault,
		Element:            attributes.Electro,
		Durability:         100,
		Mult:               burstonhit[c.TalentLvlBurst()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   0.1 * 60,
		CanBeDefenseHalted: false,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 1), burstHitmark, burstHitmark)

	// beidou burst is not hitlag extendable
	c.AddStatus(burstKey, 900, false)

	procAI := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stormbreak Proc (Q)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
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
			for _, t := range c.Core.Combat.Enemies() {
				e, ok := t.(*enemy.Enemy)
				if !ok {
					continue
				}
				e.AddResistMod(enemy.ResistMod{
					Base:  modifier.NewBaseWithHitlag("beidouc6", 900-burstHitmark),
					Ele:   attributes.Electro,
					Value: -0.15,
				})
			}
		}, burstHitmark)
	}

	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 1200)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) burstProc() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		//make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(burstICDKey) {
			c.Core.Log.NewEvent("beidou Q (active) on icd", glog.LogCharacterEvent, c.Index)
			return false
		}

		//trigger a chain of attacks starting at the first target
		atk := *c.burstAtk
		atk.SourceFrame = c.Core.F
		atk.Pattern = combat.NewDefSingleTarget(t.Key())
		cb := c.chain(c.Core.F, 1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, 1)

		c.Core.Log.NewEvent("beidou Q proc'd", glog.LogCharacterEvent, c.Index).
			Write("char", ae.Info.ActorIndex).
			Write("attack tag", ae.Info.AttackTag)

		// this ICD is most likely tied to the deployable, so it's not hitlag extendable
		c.AddStatus(burstICDKey, 60, false)
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
		trgs := c.Core.Combat.EnemyExcl(a.Target.Key())
		if len(trgs) == 0 {
			//do nothing if no other target other than this one
			return
		}
		//otherwise pick a random one
		next := c.Core.Rand.Intn(len(trgs))
		//queue an attack vs next target
		atk := *c.burstAtk
		atk.SourceFrame = src
		atk.Pattern = combat.NewDefSingleTarget(c.Core.Combat.Enemy(trgs[next]).Key())
		cb := c.chain(src, count+1)
		if cb != nil {
			atk.Callbacks = append(atk.Callbacks, cb)
		}
		c.Core.QueueAttackEvent(&atk, 1)

	}
}
