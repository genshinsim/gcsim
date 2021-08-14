package raiden

import (
	"github.com/genshinsim/gsim/pkg/core"
)

/**
let style = document.createElement('style');
style.innerHTML = '*{ user-select: auto !important; }';
document.body.appendChild(style);
**/

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)

	if c.Sim.Status("raidenburst") > 0 {
		return c.swordAttack(f)
	}

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			//fmt.Sprintf("Normal %v", c.NormalCounter),
			"Normal",
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		c.QueueDmg(&d, f-2+i)
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) swordAttack(f int) int {

	for i, mult := range attackB[c.NormalCounter] {
		d := c.Snapshot(
			//fmt.Sprintf("Musou Isshin %v", c.NormalCounter),
			"Musou Isshin",
			core.AttackTagElementalBurst,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Electro,
			25,
			mult[c.TalentLvlAttack()],
		)
		//adds to talent %
		d.Mult += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		c.AddTask(func() {
			c.Sim.ApplyDamage(&d)
			//restore energy
			if c.Sim.Frame() > c.restoreICD && c.restoreCount < 5 {
				c.restoreCount++
				c.restoreICD = c.Sim.Frame() + 60 //once every 1 second
				energy := burstRestore[c.TalentLvlBurst()]
				//apply a4
				excess := int(d.Stats[core.ER] / 0.01)
				c.Log.Debugw("a4 energy restore stacks", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "char", c.Index, "stacks", excess, "increase", float64(excess)*0.006)
				energy = energy * (1 + float64(excess)*0.006)
				for _, char := range c.Sim.Characters() {
					char.AddEnergy(energy)
				}

			}
		}, "raiden-attack", f-2+i)
	}

	c.AdvanceNormalIndex()

	return f
}

/**
The Raiden Shogun unveils a shard of her Euthymia, dealing Electro DMG to nearby opponents, and granting nearby party members the Eye of Stormy Judgment.
Eye of Stormy Judgment
**/

func (c *char) Skill(p map[string]int) int {
	f := c.ActionFrames(core.ActionSkill, p)
	d := c.Snapshot(
		"Eye of Stormy Judgement",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll

	c.QueueDmg(&d, f)

	count := 3
	if c.Sim.Rand().Float64() < 0.5 {
		count = 4
	}
	c.QueueParticle("raiden", count, core.Electro, f+100)

	//activate eye
	c.Sim.AddStatus("raidenskill", 1500+f)

	c.SetCD(core.ActionSkill, 600)
	return f
}

/**
When characters with this buff attack and hit opponents, the Eye will unleash a coordinated attack, dealing AoE Electro DMG at the opponent's position.
The Eye can initiate one coordinated attack every 0.9s per party.
**/
func (c *char) eyeOnDamage() {
	c.Sim.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		//ignore if eye on icd
		if c.eyeICD > c.Sim.Frame() {
			return
		}
		//ignore if eye not active
		if c.Sim.Status("raidenskill") == 0 {
			return
		}
		//trigger a strike
		//does not snapshot, so this is fine
		d := c.Snapshot(
			"Eye of Stormy Judgement (Strike)",
			core.AttackTagElementalArt,
			core.ICDTagElementalArt,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Electro,
			25,
			skillTick[c.TalentLvlSkill()],
		)
		d.Targets = core.TargetAll
		//https://streamable.com/28at4f hit mark 857, eye land 862
		c.QueueDmg(&d, 5)
		c.eyeICD = c.Sim.Frame() + 54 //0.9 sec icd

	}, "raiden-eye")
}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(core.ActionBurst, p)

	//activate burst, reset stacks
	c.stacksConsumed = c.stacks
	c.stacks = 0
	c.Sim.AddStatus("raidenburst", 420+f) //7 seconds
	c.restoreCount = 0
	c.restoreICD = 0
	c.c6Count = 0
	c.c6ICD = 0

	if c.Base.Cons >= 4 {
		val := make([]float64, core.EndStatType)
		val[core.ATKP] = 0.3
		for i, char := range c.Sim.Characters() {
			if i == c.Index {
				continue
			}
			char.AddMod(core.CharStatMod{
				Key:    "raiden-c2",
				Expiry: 600, //10s
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return val, true
				},
			})
		}
	}

	c.Log.Debugw("resolve stacks", "frame", c.Sim.Frame(), "event", core.LogCharacterEvent, "char", c.Index, "stacks", c.stacksConsumed)

	d := c.Snapshot(
		"Musou Shinsetsu",
		core.AttackTagElementalBurst,
		core.ICDTagElementalBurst,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Electro,
		50,
		burstBase[c.TalentLvlBurst()],
	)
	d.Targets = core.TargetAll
	d.Mult += resolveBaseBonus[c.TalentLvlBurst()] * c.stacksConsumed

	c.QueueDmg(&d, f)

	if c.Base.Cons >= 2 {
		c.AddTask(func() {
			for _, t := range c.Sim.Targets() {
				t.AddDefMod("raiden-c4", -0.6, 420)
			}
		}, "c4", f+1)
	}

	c.SetCD(core.ActionBurst, 18*60) //20s cd
	c.Energy = 0
	return f
}

func (c *char) onBurstStackCount() {
	c.Sim.AddEventHook(func(s core.Sim) bool {
		if s.ActiveCharIndex() == c.Index {
			return false
		}
		char, _ := s.CharByPos(s.ActiveCharIndex())
		//add stacks based on char max energy
		stacks := resolveStackGain[c.TalentLvlBurst()] * char.MaxEnergy()
		if c.Base.Cons > 0 {
			if char.Ele() == core.Electro {
				stacks = stacks * 1.8
			} else {
				stacks = stacks * 1.2
			}
		}
		c.stacks += stacks
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-stacks", core.PostBurstHook)
	//a4 stack gain
	particleICD := 0
	c.Sim.AddEventHook(func(s core.Sim) bool {
		if particleICD > s.Frame() {
			return false
		}
		particleICD = s.Frame() + 180 // once every 3 seconds
		c.stacks += 2
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-particle-stack", core.PostParticleHook)
}
