package mona

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 107

func init() {
	burstFrames = frames.InitAbilSlice(127) // Q -> Swap
	burstFrames[action.ActionAttack] = 121  // Q -> N1
	burstFrames[action.ActionCharge] = 118  // Q -> CA
	burstFrames[action.ActionSkill] = 115   // Q -> E
	burstFrames[action.ActionDash] = 115    // Q -> D
	burstFrames[action.ActionJump] = 104    // Q -> J
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	//bubble deal 0 dmg hydro app
	//add bubble status, when bubble status disappears trigger omen dmg the frame after
	//bubble status bursts either -> takes dmg no freeze OR freeze and freeze disappears

	//apply first non damage after 1.7 seconds
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Initial)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
	}
	cb := func(a combat.AttackCB) {
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		//bubble is applied to each target on a per target basis
		//lasts 8 seconds if not popped normally
		t.AddStatus(bubbleKey, 481, true) //1 frame extra so we don't run into problems breaking
		c.Core.Log.NewEvent("mona bubble on target", glog.LogCharacterEvent, c.Index).
			Write("char", c.Index)
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), -1, burstHitmark, cb)

	//queue a 0 damage attack to break bubble after 8 sec if bubble not broken yet
	aiBreak := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Break)",
		AttackTag:  attacks.AttackTagMonaBubbleBreak,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
		Durability: 0,
		Mult:       0,
	}
	c.Core.QueueAttack(aiBreak, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), -1, burstHitmark+480)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(5)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionJump], // earliest cancel is before burstHitmark
		State:           action.BurstState,
	}
}

func (c *char) burstDamageBonus() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = dmgBonus[c.TalentLvlBurst()]
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("mona-omen", -1),
			Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				//ok only if either bubble or omen is present
				if x.StatusIsActive(bubbleKey) || x.StatusIsActive(omenKey) {
					return m, true
				}
				return nil, false
			},
		})
	}
}

// bubble bursts when hit by an attack either while not frozen, or when the attack breaks freeze
// i.e. impulse > 0
func (c *char) burstHook() {
	//hook on to OnDamage; leave this always active
	//since freeze will trigger an attack, this should be ok
	//TODO: this implementation would currently cause bubble to break immediately on the first EC tick.
	//According to: https://docs.google.com/document/d/1pXlgCaYEpoizMIP9-QKlSkQbmRicWfrEoxb9USWD1Ro/edit#
	//only 2nd ec tick should break
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		//ignore if target doesn't have debuff
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.StatusIsActive(bubbleKey) {
			return false
		}
		//always break if it's due to time up
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagMonaBubbleBreak {
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

func (c *char) triggerBubbleBurst(t *enemy.Enemy) {
	//remove bubble tag
	t.DeleteStatus(bubbleKey)
	//add omen debuff
	t.AddStatus(omenKey, omenDuration[c.TalentLvlBurst()], true)
	//trigger dmg
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Explosion)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       explosion[c.TalentLvlBurst()],
	}
	c.Core.QueueAttack(ai, combat.NewSingleTargetHit(t.Key()), 1, 1)
}
