package diluc

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
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, core.TargettableEnemy), f, f)
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	// reset counter
	if c.Core.F >= c.eWindow {
		c.eCounter = 0
	}

	f, a := c.ActionFrames(core.ActionSkill, p)

	orb := 1
	if c.Core.Rand.Float64() < 0.33 {
		orb = 2
	}
	c.QueueParticle("Diluc", orb, core.Pyro, f+60)

	//actual skill cd starts immediately on first cast
	//times out after 4 seconds of not using
	//every hit applies pyro
	//apply attack speed

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Searing Onslaught %v", c.eCounter),
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       skill[c.eCounter][c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f)

	//add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			c.Core.Status.AddStatus("dilucc4", 120) //effect lasts 2 seconds
		}, "dilucc4", f+120) // 2seconds after cast
	}

	// allow skill to be used again if 4s hasn't passed since last use
	c.eWindow = c.Core.F + 60*4

	c.eCounter++
	switch c.eCounter {
	case 1:
		// TODO: cd delay?
		// set cd on first use
		c.SetCD(core.ActionSkill, 10*60)
	case 3:
		// reset window since we're at 3rd use
		c.eWindow = -1
		c.eCounter = 0
	}

	//return animation cd
	//this also depends on which hit in the chain this is
	return f, a
}

func (c *char) Burst(p map[string]int) (int, int) {

	dot, ok := p["dot"]
	if !ok {
		dot = 2 //number of dot hits
	}
	if dot > 7 {
		dot = 7
	}
	explode, ok := p["explode"]
	if !ok {
		explode = 0 //if explode hits
	}

	c.Core.Status.AddStatus("dilucq", 720)
	f, a := c.ActionFrames(core.ActionBurst, p)

	//enhance weapon for 12 seconds
	// Infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "diluc-fire-weapon",
		Ele:    core.Pyro,
		Tags:   []core.AttackTag{core.AttackTagNormal, core.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.F + 720, //with a4
	})

	// add 20% pyro damage
	val := make([]float64, core.EndStatType)
	val[core.PyroP] = 0.2
	c.AddMod(core.CharStatMod{
		Key:    "diluc-fire-weapon",
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.F + 720,
	})

	// Snapshot occurs late in the animation when it is released from the claymore
	// For our purposes, snapshot upon damage proc
	c.AddTask(func() {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Dawn (Strike)",
			AttackTag:  core.AttackTagElementalBurst,
			ICDTag:     core.ICDTagElementalBurst,
			ICDGroup:   core.ICDGroupDiluc,
			StrikeType: core.StrikeTypeBlunt,
			Element:    core.Pyro,
			Durability: 50,
			Mult:       burstInitial[c.TalentLvlBurst()],
		}

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 1)

		//dot does damage every .2 seconds for 7 hits? so every 12 frames
		//dot does max 7 hits + explosion, roughly every 13 frame? blows up at 210 frames
		//first tick did 50 dur as well?
		ai.Abil = "Dawn (Tick)"
		ai.Mult = burstDOT[c.TalentLvlBurst()]
		for i := 1; i <= dot; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, i+12)
		}

		if explode > 0 {
			ai.Abil = "Dawn (Explode)"
			ai.Mult = burstExplode[c.TalentLvlBurst()]
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 110)
		}
	}, "diluc-burst", f-1)

	c.ConsumeEnergy(21)
	c.SetCDWithDelay(core.ActionBurst, 720, 14)
	return f, a
}

//Diluc dash is 19 frames
func (c *char) Dash(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionDash, p)
	return f, a
}
