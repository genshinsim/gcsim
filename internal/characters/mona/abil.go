package mona

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) Attack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Normal",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagMonaWaterDamage,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	var cb core.AttackCBFunc

	if c.Base.Cons > 1 {
		cb = c.c2cb
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, f-1, cb)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), 0, f-1)

	return f, a
}

func (c *char) Dash(p map[string]int) (int, int) {
	f, ok := p["f"]
	if !ok {
		f = 36
	}
	//no dmg attack at end of dash
	ai := core.AttackInfo{
		Abil:       "Dash",
		ActorIndex: c.Index,
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f)

	//After she has used Illusory Torrent for 2s, if there are any opponents nearby,
	//Mona will automatically create a Phantom.
	//A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.

	//TODO: a4 not implemented. needs to know if this can be created while already on one field
	//and if it overrides

	return f, f
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Tick)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       skillDot[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	//5.22 seconds duration after cast
	//tick every 1 sec
	for i := 60; i < 313; i += 60 {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), f+i)
	}

	aiExplode := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Mirror Reflection of Doom (Explode)",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(aiExplode, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, f+313)

	count := 3
	if c.Core.Rand.Float64() < .33 {
		count = 4
	}
	c.QueueParticle("mona", count, core.Hydro, f+313+100)

	c.SetCD(core.ActionSkill, 12*60)

	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {
	//bubble deal 0 dmg hydro app
	//add bubble status, when bubble status disappears trigger omen dmg the frame after
	//bubble status bursts either -> takes dmg no freeze OR freeze and freeze disappears
	f, a := c.ActionFrames(core.ActionBurst, p)

	//apply first non damage after 1.7 seconds
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Initial)",
		AttackTag:  core.AttackTagNone,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       0,
	}
	cb := func(a core.AttackCB) {
		//bubble is applied to each target on a per target basis
		//lasts 8 seconds if not popped normally
		a.Target.SetTag(bubbleKey, c.Core.F+481) //1 frame extra so we don't run into problems breaking
		c.Core.Log.Debugw("mona bubble on target", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index)
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(4, false, core.TargettableEnemy), -1, 102, cb)

	//queue a 0 damage attack to break bubble after 8 sec if bubble not broken yet
	aiBreak := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Break)",
		AttackTag:  core.AttackTagMonaBubbleBreak,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 0,
		Mult:       0,
	}
	c.Core.Combat.QueueAttack(aiBreak, core.NewDefCircHit(4, false, core.TargettableEnemy), -1, 102+480)

	c.SetCD(core.ActionBurst, 15*60)
	c.ConsumeEnergy(13)
	return f, a
}

//bubble bursts when hit by an attack either while not frozen, or when the attack breaks freeze
//i.e. impulse > 0
func (c *char) burstHook() {
	//hook on to OnDamage; leave this always active
	//since freeze will trigger an attack, this should be ok
	//TODO: this implementation would currently cause bubble to break immediately on the first EC tick.
	//According to: https://docs.google.com/document/d/1pXlgCaYEpoizMIP9-QKlSkQbmRicWfrEoxb9USWD1Ro/edit#
	//only 2nd ec tick should break
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		//ignore if target doesn't have debuff
		t := args[0].(core.Target)
		if t.GetTag(bubbleKey) < c.Core.F {
			return false
		}
		//always break if it's due to time up
		atk := args[1].(*core.AttackEvent)
		if atk.Info.AttackTag == core.AttackTagMonaBubbleBreak {
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

func (c *char) triggerBubbleBurst(t core.Target) {
	//remove bubble tag
	t.RemoveTag(bubbleKey)
	//add omen debuff
	t.SetTag(omenKey, c.Core.F+omenDuration[c.TalentLvlBurst()])
	//trigger dmg
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Illusory Bubble (Explosion)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Hydro,
		Durability: 50,
		Mult:       explosion[c.TalentLvlBurst()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(t.Index(), t.Type()), 1, 1)
}
