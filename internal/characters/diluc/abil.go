package diluc

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(1, false, coretype.TargettableEnemy), f-1, f-1)
	c.AdvanceNormalIndex()

	return f, a
}

func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	if c.eCounter == 0 {
		c.eStarted = true
		c.eStartFrame = c.Core.Frame
	}
	c.eLastUse = c.Core.Frame

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
		Abil:       "Searing Onslaught",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Pyro,
		Durability: 25,
		Mult:       skill[c.eCounter][c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), f-5, f-5)

	//add a timer to activate c4
	if c.Base.Cons >= 4 {
		c.AddTask(func() {
			c.Core.AddStatus("dilucc4", 120) //effect lasts 2 seconds
		}, "dilucc4", f+120) // 2seconds after cast
	}

	//skill only goes on cd once all 3 charges have been used
	//or if 4 second passed since last use, skill will also go on cd

	c.eCounter++
	if c.eCounter == 3 {
		//ability can go on cd now
		cd := 600 - (c.Core.Frame - c.eStartFrame)
		c.coretype.Log.NewEvent("diluc skill going on cd", coretype.LogCharacterEvent, c.Index, "duration", cd)
		c.SetCD(core.ActionSkill, cd)
		c.eStarted = false
		c.eStartFrame = -1
		c.eLastUse = -1
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

	c.Core.AddStatus("dilucq", 720)
	f, a := c.ActionFrames(core.ActionBurst, p)

	//enhance weapon for 12 seconds
	// Infusion starts when burst starts and ends when burst comes off CD - check any diluc video
	c.AddWeaponInfuse(core.WeaponInfusion{
		Key:    "diluc-fire-weapon",
		Ele:    core.Pyro,
		Tags:   []core.AttackTag{coretype.AttackTagNormal, coretype.AttackTagExtra, core.AttackTagPlunge},
		Expiry: c.Core.Frame + 720, //with a4
	})

	// add 20% pyro damage
	val := make([]float64, core.EndStatType)
	val[core.PyroP] = 0.2
	c.AddMod(coretype.CharStatMod{
		Key:    "diluc-fire-weapon",
		Amount: func() ([]float64, bool) { return val, true },
		Expiry: c.Core.Frame + 720,
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

		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, 1)

		//dot does damage every .2 seconds for 7 hits? so every 12 frames
		//dot does max 7 hits + explosion, roughly every 13 frame? blows up at 210 frames
		//first tick did 50 dur as well?
		ai.Abil = "Dawn (Tick)"
		ai.Mult = burstDOT[c.TalentLvlBurst()]
		for i := 1; i <= dot; i++ {
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, i+12)
		}

		if explode > 0 {
			ai.Abil = "Dawn (Explode)"
			ai.Mult = burstExplode[c.TalentLvlBurst()]
			c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, coretype.TargettableEnemy), 0, 110)
		}
	}, "diluc-burst", 100)

	c.ConsumeEnergy(24)
	c.SetCDWithDelay(core.ActionBurst, 720, 24)
	return f, a
}
