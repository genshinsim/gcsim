package travelerelectro

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeSlash,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f-1, f-1)

	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionSkill, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lightning Blade",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	hits, ok := p["hits"]
	if !ok {
		hits = 1
	} else if hits > 3 {
		hits = 3
	}

	maxAmulets := 2
	if c.Base.Cons >= 1 {
		maxAmulets = 3
	}

	// clear existing amulets
	c.abundanceAmulets = 0

	// accept param to limit the amount of amulets generated
	pMaxAmulets, ok := p["max_amulets"]
	if ok && pMaxAmulets < maxAmulets {
		maxAmulets = pMaxAmulets
	}

	// Counting from the frame E is pressed, it takes an average of 1.79 seconds for a character to be able to pick one up
	// https://library.keqingmains.com/evidence/characters/electro/traveler-electro#amulets-delay
	amuletDelay := p["amulet_delay"]
	//make it so that it can't be faster than 1.79s
	if amuletDelay < 107 {
		amuletDelay = 107 // ~1.79s
	}

	//particles appear to be generated if the blades lands but capped at 1
	partCount := 0
	particlesCB := func(atk core.AttackCB) {
		if partCount > 0 {
			return
		}
		partCount++
		c.QueueParticle(c.Name(), 1, core.Electro, 100) //this way we're future proof if for whatever reason this misses
	}

	amuletCB := func(atk core.AttackCB) {
		// generate amulet if generated amulets < limit
		if c.abundanceAmulets >= maxAmulets {
			return
		}

		// 1 amulet per attack
		c.abundanceAmulets++
		c.AddTag("generated", c.abundanceAmulets)

		c.Core.Log.NewEvent("travelerelectro abundance amulet generated", core.LogCharacterEvent, c.Index, "amulets", c.abundanceAmulets)
	}

	for i := 0; i < hits; i++ {
		c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(0.3, false, core.TargettableEnemy), f, particlesCB, amuletCB)
	}

	// try to pick up amulets
	c.AddTask(func() {
		activeChar := c.Core.Chars[c.Core.ActiveChar]
		c.collectAmulets(c.Core.F, activeChar)
	}, "pickup-abundance-amulets", amuletDelay)

	c.SetCDWithDelay(core.ActionSkill, 810, 21) //13.5s, starts 21 frames in

	return f, a
}

/**
[12:01 PM] pai: never tried to measure it but emc burst looks like it has roughly 1~1.5 abyss tile of range, skill goes a bit further i think
[12:01 PM] pai: the 3 hits from the skill also like split out and kind of auto target if that's useful information
**/
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Bellowing Thunder",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 0, f)

	//1573 start, 1610 cd starts, 1612 energy drained, 1633 first swapable
	c.ConsumeEnergy(42)
	c.SetCD(core.ActionBurst, 1200+37)

	c.Core.Status.AddStatus("travelerelectroburst", 720) // 12s

	procAI := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Falling Thunder Proc (Q)",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagElementalBurst,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}
	c.burstSnap = c.Snapshot(&procAI)
	c.burstAtk = &core.AttackEvent{
		Info:     procAI,
		Snapshot: c.burstSnap,
	}
	c.burstSrc = c.Core.F

	return f, a
}

func (c *char) burstProc() {
	icd := 0

	// Lightning Shroud
	//  When your active character's Normal or Charged Attacks hit opponents, they will call Falling Thunder forth, dealing Electro DMG.
	//  When Falling Thunder hits opponents, it will regenerate Energy for that character.
	//  One instance of Falling Thunder can be generated every 0.5s.
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		t := args[0].(core.Target)

		// only apply on na/ca
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		// make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.ActiveChar {
			return false
		}
		// only apply if burst is active
		if c.Core.Status.Duration("travelerelectroburst") == 0 {
			return false
		}
		// One instance of Falling Thunder can be generated every 0.5s.
		if icd > c.Core.F {
			c.Core.Log.NewEvent("travelerelectro Q (active) on icd", core.LogCharacterEvent, c.Index)
			return false
		}

		// Use burst snapshot, update target & source frame
		atk := *c.burstAtk
		atk.SourceFrame = c.Core.F
		//attack is 2 (or 2.5 for enhanced) aoe centered on target
		x, y := t.Shape().Pos()
		atk.Pattern = core.NewCircleHit(x, y, 2, false, core.TargettableEnemy)

		// C2 - Violet Vehemence
		// When Falling Thunder created by Bellowing Thunder hits an opponent, it will decrease their Electro RES by 15% for 8s.
		// c6 - World-Shaker
		//  Every 2 Falling Thunder attacks triggered by Bellowing Thunder will greatly increase the DMG
		//  dealt by the next Falling Thunder, which will deal 200% of its original DMG and will restore
		//  an additional 1 Energy to the current character.
		c.c6Damage(&atk)
		atk.Callbacks = append(atk.Callbacks, c.fallingThunderEnergy(), c.c2(t), c.c6Energy())

		c.Core.Combat.QueueAttackEvent(&atk, 1)

		c.Core.Log.NewEvent("travelerelectro Q proc'd", core.LogCharacterEvent, c.Index, "char", ae.Info.ActorIndex, "attack tag", ae.Info.AttackTag)

		icd = c.Core.F + 30 // 0.5s
		return false
	}, "travelerelectro-bellowingthunder")
}

func (c *char) fallingThunderEnergy() core.AttackCBFunc {
	return func(a core.AttackCB) {
		// Regenerate 1 flat energy for the active character
		activeChar := c.Core.Chars[c.Core.ActiveChar]
		activeChar.AddEnergy("travelerelectro-fallingthunder", burstRegen[c.TalentLvlBurst()])
	}
}
