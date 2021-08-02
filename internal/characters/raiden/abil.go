package raiden

import (
	"github.com/genshinsim/gsim/pkg/def"
)

/**
let style = document.createElement('style');
style.innerHTML = '*{ user-select: auto !important; }';
document.body.appendChild(style);
**/

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(def.ActionAttack, p)

	if c.Sim.Status("raidenburst") > 0 {
		return c.swordAttack(f)
	}

	for i, mult := range attack[c.NormalCounter] {
		d := c.Snapshot(
			//fmt.Sprintf("Normal %v", c.NormalCounter),
			"Normal",
			def.AttackTagNormal,
			def.ICDTagNormalAttack,
			def.ICDGroupDefault,
			def.StrikeTypeSlash,
			def.Physical,
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
			def.AttackTagNormal,
			def.ICDTagNormalAttack,
			def.ICDGroupDefault,
			def.StrikeTypeSlash,
			def.Physical,
			25,
			mult[c.TalentLvlAttack()],
		)
		//assume adding to bonus dmg (not talent %)
		d.Stats[def.DmgP] += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed
		c.AddTask(func() {
			c.Sim.ApplyDamage(&d)
			//restore energy
			if c.Sim.Frame() > c.restoreICD && c.restoreCount < 5 {
				c.restoreCount++
				c.restoreICD = c.Sim.Frame() + 60 //once every 1 second
				energy := burstRestore[c.TalentLvlBurst()]
				//apply a4
				excess := int(d.Stats[def.ER] / 0.01)
				c.Log.Debugw("a4 energy restore stacks", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index, "stacks", excess, "increase", float64(excess)*0.006)
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
	f := c.ActionFrames(def.ActionSkill, p)
	d := c.Snapshot(
		"Eye of Stormy Judgement",
		def.AttackTagElementalArt,
		def.ICDTagElementalArt,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		25,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = def.TargetAll

	c.QueueDmg(&d, f)

	count := 3
	if c.Sim.Rand().Float64() < 0.5 {
		count = 4
	}
	c.QueueParticle("raiden", count, def.Electro, f+100)

	//activate eye
	//does eye snapshot?
	c.Sim.AddStatus("raidenskill", 1500+f)

	c.SetCD(def.ActionSkill, 600)
	return f
}

/**
When characters with this buff attack and hit opponents, the Eye will unleash a coordinated attack, dealing AoE Electro DMG at the opponent's position.
The Eye can initiate one coordinated attack every 0.9s per party.
**/
func (c *char) eyeOnDamage() {
	c.Sim.AddOnAttackLanded(func(t def.Target, ds *def.Snapshot, dmg float64, crit bool) {
		//ignore if eye on icd
		if c.eyeICD > c.Sim.Frame() {
			return
		}
		//ignore if eye not active
		if c.Sim.Status("raidenskill") == 0 {
			return
		}
		//trigger a strike
		d := c.Snapshot(
			"Eye of Stormy Judgement (Strike)",
			def.AttackTagElementalArt,
			def.ICDTagElementalArt,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Electro,
			25,
			skillTick[c.TalentLvlSkill()],
		)
		d.Targets = def.TargetAll
		c.QueueDmg(&d, 1)
		c.eyeICD = c.Sim.Frame() + 54 //0.9 sec icd

	}, "raiden-eye")
}

func (c *char) Burst(p map[string]int) int {

	f := c.ActionFrames(def.ActionBurst, p)

	//activate burst, reset stacks
	c.stacksConsumed = c.stacks
	c.stacks = 0
	c.Sim.AddStatus("raidenburst", 420+f) //7 seconds
	c.restoreCount = 0
	c.restoreICD = 0

	c.Log.Debugw("resolve stacks", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "char", c.Index, "stacks", c.stacksConsumed)

	d := c.Snapshot(
		"Musou Shinsetsu",
		def.AttackTagElementalBurst,
		def.ICDTagElementalBurst,
		def.ICDGroupDefault,
		def.StrikeTypeDefault,
		def.Electro,
		50,
		burstBase[c.TalentLvlBurst()],
	)
	d.Targets = def.TargetAll
	d.Stats[def.DmgP] += resolveBonus[c.TalentLvlBurst()] * c.stacksConsumed

	c.QueueDmg(&d, f)

	c.SetCD(def.ActionBurst, 20*60) //20s cd
	c.Energy = 0
	return f
}

func (c *char) onBurstStackCount() {
	c.Sim.AddEventHook(func(s def.Sim) bool {
		if s.ActiveCharIndex() == c.Index {
			return false
		}
		char, _ := s.CharByPos(s.ActiveCharIndex())
		//add stacks based on char max energy
		c.stacks += resolveStackGain[c.TalentLvlBurst()] * char.MaxEnergy()
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-stacks", def.PostBurstHook)
	//a4 stack gain
	particleICD := 0
	c.Sim.AddEventHook(func(s def.Sim) bool {
		if particleICD > s.Frame() {
			return false
		}
		particleICD = s.Frame() + 180 // once every 3 seconds
		c.stacks += 2
		if c.stacks > 60 {
			c.stacks = 60
		}
		return false
	}, "raiden-particle-stack", def.PostParticleHook)
}
