package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int

const burstHitmark = 102

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//bubble deal 0 dmg hydro app
	//add bubble status, when bubble status disappears trigger omen dmg the frame after
	//bubble status bursts either -> takes dmg no freeze OR freeze and freeze disappears

	//apply first non damage after 1.7 seconds
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Initial)",
		AttackTag:  combat.AttackTagNone,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
	}
	cb := func(a combat.AttackCB) {
		//bubble is applied to each target on a per target basis
		//lasts 8 seconds if not popped normally
		a.Target.SetTag(bubbleKey, c.Core.F+481) //1 frame extra so we don't run into problems breaking
		c.Core.Log.NewEvent("mona bubble on target", glog.LogCharacterEvent, c.Index, "char", c.Index)
	}
	c.Core.QueueAttack(ai, combat.NewDefCircHit(4, false, combat.TargettableEnemy), -1, burstHitmark, cb)

	//queue a 0 damage attack to break bubble after 8 sec if bubble not broken yet
	aiBreak := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Break)",
		AttackTag:  combat.AttackTagMonaBubbleBreak,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       0,
	}
	c.Core.QueueAttack(aiBreak, combat.NewDefCircHit(4, false, combat.TargettableEnemy), -1, burstHitmark+480)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 13)
	c.ConsumeEnergy(13)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		Post:            burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) burstDamageBonus() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = dmgBonus[c.TalentLvlBurst()]
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod("mona-omen", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			//ignore if omen or bubble not present
			if t.GetTag(bubbleKey) < c.Core.F && t.GetTag(omenKey) < c.Core.F {
				return nil, false
			}
			return m, true
		})
	}
}

//bubble bursts when hit by an attack either while not frozen, or when the attack breaks freeze
//i.e. impulse > 0
func (c *char) burstHook() {
	//hook on to OnDamage; leave this always active
	//since freeze will trigger an attack, this should be ok
	//TODO: this implementation would currently cause bubble to break immediately on the first EC tick.
	//According to: https://docs.google.com/document/d/1pXlgCaYEpoizMIP9-QKlSkQbmRicWfrEoxb9USWD1Ro/edit#
	//only 2nd ec tick should break
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		//ignore if target doesn't have debuff
		t := args[0].(combat.Target)
		if t.GetTag(bubbleKey) < c.Core.F {
			return false
		}
		//always break if it's due to time up
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag == combat.AttackTagMonaBubbleBreak {
			c.triggerBubbleBurst(t)
			return false
		}
		//dont break if no impulse
		if atk.Info.NoImpulse {
			return false
		}
		//otherwise break on damage
		c.triggerBubbleBurst(t)

		return false
	}, "mona-bubble-check")
}

func (c *char) triggerBubbleBurst(t combat.Target) {
	//remove bubble tag
	t.RemoveTag(bubbleKey)
	//add omen debuff
	t.SetTag(omenKey, c.Core.F+omenDuration[c.TalentLvlBurst()])
	//trigger dmg
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Explosion)",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       explosion[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewDefSingleTarget(t.Index(), t.Type()), 1, 1)
}
