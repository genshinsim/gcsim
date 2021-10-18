package childe

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

// Normal attack
// Perform up to 6 consecutive shots with a bow.
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	f, a := c.ActionFrames(core.ActionAttack, p)

	if c.Core.Status.Duration("childemelee") > 0 {
		return c.meleeAttack(f, a)
	}

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f+travel)

	c.AdvanceNormalIndex()

	return f, a
}

// random delayed numbers
var meleeDelayOffset = [][]int{
	{1},
	{1},
	{1},
	{1},
	{1},
	{1, 1},
}

// Melee stance attack.
// Perform up to 6 consecutive Hydro strikes.
func (c *char) meleeAttack(f, a int) (int, int) {
	for i, mult := range eAttack[c.NormalCounter] {
		d := c.Snapshot(
			fmt.Sprintf("Normal %v", c.NormalCounter),
			//"Normal",
			core.AttackTagNormal,
			core.ICDTagNormalAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Hydro,
			25,
			mult[c.TalentLvlSkill()],
		)

		c.AddTask(func() {
			c.Core.Combat.ApplyDamage(&d)
		}, "childe-attack", f-meleeDelayOffset[c.NormalCounter][i])

	}
	c.AdvanceNormalIndex()

	return f, a
}

//Perform a more precise Aimed Shot. Once fully charged, deal Hydro DMG and apply the Riptide status.
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 20
	}

	d := c.Snapshot(
		"Aim (Charged)",
		core.AttackTagExtra,
		core.ICDTagNone,
		core.ICDGroupDefault,
		core.StrikeTypePierce,
		core.Hydro,
		25,
		aim[c.TalentLvlAttack()],
	)

	d.HitWeakPoint = true
	d.AnimationFrames = f

	c.QueueDmg(&d, travel+f)

	return f, a
}

//Charged Attack: Consume 20 Stamina to unleash a cross slash, dealing Hydro DMG.
// hitWeakPoint: childe can proc Prototype Cresent's Passive on Geovishap's weakspots.
// Evidence: https://youtu.be/oOfeu5pW0oE
func (c *char) ChargeAttack(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionCharge, p)

	if c.Core.Status.Duration("childemelee") == 0 {
		return f, a
	}

	hitWeakPoint, ok := p["hitWeakPoint"]
	if !ok {
		hitWeakPoint = 0
	}

	for i, mult := range eCharge {
		d := c.Snapshot(
			fmt.Sprintf("Charge %v", i),
			core.AttackTagExtra,
			core.ICDTagExtraAttack,
			core.ICDGroupDefault,
			core.StrikeTypeSlash,
			core.Hydro,
			25,
			mult[c.TalentLvlSkill()],
		)
		if hitWeakPoint != 0 {
			d.HitWeakPoint = true
		}

		c.AddTask(func() {
			c.Core.Combat.ApplyDamage(&d)
		}, "childe-charge-attack", f-5)
	}
	return f, a
}

//Unleashes a set of weaponry made of pure water, dealing Hydro DMG to surrounding opponents and entering Melee Stance.
//Melee Stance: Converts Tartaglia’s Normal and Charged Attacks into Hydro DMG.Cannot be overridden by any other elemental infusion.
func (c *char) Skill(p map[string]int) (int, int) {
	if c.Core.Status.Duration("childemelee") > 0 {
		f, a := c.ActionFrames(core.ActionSkill, p)
		c.Core.Status.DeleteStatus("childemelee")
		newCD := float64(c.Core.F - c.eCast + 6*60)
		//Foul Legacy: Tide Withholder. Decreases the CD of Foul Legacy: Raging Tide by 20%
		if c.Base.Cons >= 1 {
			newCD *= 0.8
		}
		if c.Base.Cons >= 6 && c.c6 {
			newCD = 0
			c.c6 = false
		}
		c.Core.Log.Debugw("Childe leaving melee stance", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
			c.Core.Status.Duration("childemelee"))
		c.SetCD(core.ActionSkill, int(newCD))

		c.ResetNormalCounter()
		return f, a
	}

	f, a := c.ActionFrames(core.ActionSkill, p)
	c.eCast = c.Core.F
	c.Core.Status.AddStatus("childemelee", 30*60)
	c.Core.Log.Debugw("Foul Legacy acivated", "frame", c.Core.F, "event", core.LogCharacterEvent, "expiry", c.Core.F+30*60)

	d := c.Snapshot(
		"Foul Legacy: Raging Tide",
		core.AttackTagElementalArt,
		core.ICDTagElementalArt,
		core.ICDGroupDefault,
		core.StrikeTypeDefault,
		core.Hydro,
		50,
		skill[c.TalentLvlSkill()],
	)
	d.Targets = core.TargetAll
	c.QueueDmg(&d, f)

	//If the skill ended automatically after 30s, the CD is even longer. (45s cd)
	c.AddTask(func() {
		//if he's already out melee, do nothing
		if c.Core.Status.Duration("childemelee") == 0 {
			return
		}

		c.Core.Status.DeleteStatus("childemelee")
		newCD := float64(45 * 60) // cd 45s
		//Foul Legacy: Tide Withholder. Decreases the CD of Foul Legacy: Raging Tide by 20%
		if c.Base.Cons >= 1 {
			newCD *= 0.8
		}
		if c.Base.Cons >= 6 && c.c6 {
			newCD = 0
			c.c6 = false
		}
		c.SetCD(core.ActionSkill, int(newCD))
	}, "childe-exit-melee", 30*60)

	c.SetCD(core.ActionSkill, 60)
	return f, a
}

//Performs a different attack depending on the stance in which it is cast.
//Ranged Stance: Flash of Havoc (Fire a Hydro-imbued magic arrow, dealing AoE Hydro DMG. Apply Riptide status to enemies hit. Returns 20 Energy after use.)
//Melee Stance: Light of Obliteration (Performs a slash with a large AoE, dealing massive Hydro DMG. Triggers Riptide Blast)
func (c *char) Burst(p map[string]int) (int, int) {
	mult := burst[c.TalentLvlBurst()]

	f, a := c.ActionFrames(core.ActionBurst, p)

	skillName := "Ranged Stance: Flash of Havoc"
	if c.Core.Status.Duration("childemelee") > 0 {
		skillName = "Melee Stance: Light of Obliteration"
		mult = meleeBurst[c.TalentLvlBurst()]
	}

	c.AddTask(func() {
		d := c.Snapshot(
			skillName,
			core.AttackTagElementalBurst,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Hydro,
			50,
			mult,
		)
		d.Targets = core.TargetAll
		c.Core.Combat.ApplyDamage(&d)
	}, "childe-burst-clear", f-5) //random 5 frame

	c.Energy = 0
	if c.Core.Status.Duration("childemelee") == 0 {
		c.AddEnergy(20)
		c.Core.Log.Debugw("Childe ranged burst restoring 20 energy", "frame", c.Core.F, "event", core.LogEnergyEvent, "new energy", c.Energy)
	} else {
		// C6 AnnihilationWhen Havoc: Obliteration is cast in Melee Stance, the CD of Foul Legacy: Raging Tide is reset.
		// This effect will only take place once Tartaglia returns to his Ranged Stance.
		if c.Base.Cons >= 6 {
			c.c6 = true
		}
	}
	c.SetCD(core.ActionBurst, 900)
	return f, a
}

func (c *char) rtHook() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		t := args[0].(core.Target)
		crit := args[3].(bool)
		//childe is active
		if ds.ActorIndex != c.CharIndex() {
			return false
		}
		//source is hydro
		if ds.Element != core.Hydro {
			return false
		}
		// dont proc if src from riptides
		if ds.Abil == "Riptide Flash" || ds.Abil == "Riptide Slash" ||
			ds.Abil == "Riptide Blast" || ds.Abil == "Riptide Burst" ||
			ds.Abil == "C4 Riptide Flash" || ds.Abil == "C4 Riptide Slash" {
			return false
		}

		switch ds.AttackTag {
		case core.AttackTagNormal:
			if c.Core.Status.Duration("childemelee") > 0 {
				// melee normal
				if c.rtExpiry[t.Index()] > c.Core.F {
					// c.Core.Log.Debugw("Riptide Slash checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
					// 	c.Core.Status.Duration("childemelee"), "target", t, "fl", c.rtflashICD[t], "sl", c.rtslashICD[t], "rtExpiry", c.rtExpiry[t])
					if c.rtslashICD[t.Index()] < c.Core.F {
						c.AddTask(func() {
							d := c.Snapshot(
								"Riptide Slash",
								core.AttackTagElementalArt,
								core.ICDTagNone,
								core.ICDGroupDefault,
								core.StrikeTypeDefault,
								core.Hydro,
								25,
								rtSlash[c.TalentLvlSkill()],
							)
							d.Targets = core.TargetAll

							c.Core.Combat.ApplyDamage(&d)
							c.Core.Log.Debugw("Riptide Slash ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
								c.Core.Status.Duration("childemelee"), "target", t.Index(), "fl", c.rtflashICD[t.Index()], "sl", c.rtslashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])
						}, "Riptide Slash", 5)
						c.rtslashICD[t.Index()] = c.Core.F + 90 //1.5s icd
					}
				}

				// A4:Sword of TorrentsWhen Tartaglia is in Foul Legacy: Raging Tide’s Melee Stance,
				// on dealing a CRIT hit, Normal and Charged Attacks apply the Riptide status effect to opponents.
				if crit {
					c.applyRT(t.Index())
				}
			}
		case core.AttackTagExtra:
			if c.Core.Status.Duration("childemelee") > 0 {
				// melee charge
				if c.rtExpiry[t.Index()] > c.Core.F {
					// c.Core.Log.Debugw("Riptide Slash checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
					// 	c.Core.Status.Duration("childemelee"), "target", t, "fl", c.rtflashICD[t], "sl", c.rtslashICD[t], "rtExpiry", c.rtExpiry[t])
					if c.rtslashICD[t.Index()] < c.Core.F {
						c.AddTask(func() {
							d := c.Snapshot(
								"Riptide Slash",
								core.AttackTagElementalArt,
								core.ICDTagNone,
								core.ICDGroupDefault,
								core.StrikeTypeDefault,
								core.Hydro,
								25,
								rtSlash[c.TalentLvlSkill()],
							)
							d.Targets = core.TargetAll

							c.Core.Combat.ApplyDamage(&d)
							c.Core.Log.Debugw("Riptide Slash ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
								c.Core.Status.Duration("childemelee"), "target", t.Index(), "fl", c.rtflashICD[t.Index()], "sl", c.rtslashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])
						}, "Riptide Slash", 5)
						c.rtslashICD[t.Index()] = c.Core.F + 90 //1.5s icd
					}
				}

				// A4:Sword of TorrentsWhen Tartaglia is in Foul Legacy: Raging Tide’s Melee Stance,
				// on dealing a CRIT hit, Normal and Charged Attacks apply the Riptide status effect to opponents.
				if crit {
					c.applyRT(t.Index())
				}
			} else {
				// aim mode
				if c.rtExpiry[t.Index()] > c.Core.F {
					// c.Core.Log.Debugw("Riptide Flash checking for tick", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
					// 	c.Core.Status.Duration("childemelee"), "target", t, "fl", c.rtflashICD[t], "sl", c.rtslashICD[t], "rtExpiry", c.rtExpiry[t])
					if c.rtflashICD[t.Index()] < c.Core.F {
						c.AddTask(func() {
							d := c.Snapshot(
								"Riptide Flash",
								core.AttackTagNormal,
								core.ICDTagTartagliaRiptideFlash,
								core.ICDGroupDefault,
								core.StrikeTypeDefault,
								core.Hydro,
								25,
								rtFlash[0][c.TalentLvlAttack()],
							)
							d.Targets = core.TargetAll

							//proc 3 hits
							for i := 1; i <= 3; i++ {
								x := d.Clone()
								c.QueueDmg(&x, i)
							}
							c.Core.Log.Debugw("Riptide Flash ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
								c.Core.Status.Duration("childemelee"), "target", t.Index(), "fl", c.rtflashICD[t.Index()], "sl", c.rtslashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])
						}, "Riptide Flash", 5)
						c.rtflashICD[t.Index()] = c.Core.F + 42 //0.7s icd
					}
				}

				c.applyRT(t.Index())
			}
		case core.AttackTagElementalBurst:
			if c.Core.Status.Duration("childemelee") > 0 {
				//Riptide Blast: Clears Riptide status. DMG Dealt is considered Elemental Burst Damage.
				//clear riptide status
				if c.rtExpiry[t.Index()] > c.Core.F {
					// c.Core.Log.Debugw(fmt.Sprintf("%v checking for tick", sname), "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
					// 	c.Core.Status.Duration("childemelee"), "target", t, "fl", c.rtflashICD[t], "sl", c.rtslashICD[t], "rtExpiry", c.rtExpiry[t])
					c.AddTask(func() {
						d := c.Snapshot(
							"Riptide Blast",
							core.AttackTagElementalBurst,
							core.ICDTagNone,
							core.ICDGroupDefault,
							core.StrikeTypeDefault,
							core.Hydro,
							50,
							rtBlast[c.TalentLvlBurst()],
						)
						d.Targets = core.TargetAll

						c.Core.Combat.ApplyDamage(&d)
						c.Core.Log.Debugw("Riptide Blast ticked", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
							c.Core.Status.Duration("childemelee"), "target", t.Index(), "fl", c.rtflashICD[t.Index()], "sl", c.rtslashICD[t.Index()], "rtExpiry", c.rtExpiry[t.Index()])
					}, "Riptide Blast", 5)

					c.rtExpiry[t.Index()] = 0

					c.Core.Log.Debugw("Childe cleared riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
						c.Core.Status.Duration("childemelee"), "target", t.Index(), "Expiry", c.rtExpiry[t.Index()])

					if c.Base.Cons >= 4 {
						c.funcC4[t.Index()] = false
						c.mlBurstUsed = true
					}
				}
			} else {
				c.applyRT(t.Index())
			}
		default:
		}
		return false
	}, "childe-riptide")
}

//Q: does all type of childe's riptide share the same icd of particle?
//A: Hard 3s cooldown on particle generation. Particles can only be gained through riptide procs.
//Q: do all riptide flash, burst, slash and blast share same icd of particle generation?
//A: flash and slash share it iirc, burst doesn't have one, and I don't think blast gens energy
func (c *char) rtParticleGen() {
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != c.CharIndex() {
			return false
		}
		if ds.Abil != "Riptide Flash" && ds.Abil != "Riptide Slash" &&
			ds.Abil != "C4 Riptide Flash" && ds.Abil != "C4 Riptide Slash" {
			return false
		}
		if c.rtParticleICD > c.Core.F {
			// c.Core.Log.Debugw("childe particle gen on icd", "frame", c.Core.F, "event", core.LogCharacterEvent, "icd", c.rtParticleICD)
			return false
		}
		if c.rtParticleICD < c.Core.F {
			c.Core.Log.Debugw("childe gen a hydro particle", "frame", c.Core.F, "event", core.LogCharacterEvent, "icd", c.rtParticleICD)
			c.rtParticleICD = c.Core.F + 180 //3 sec
			c.QueueParticle("tartaglia", 1, core.Hydro, 100)
		}
		return false
	}, "childe-particle-gen")
}

func (c *char) c4TickFunc(t int) func() {
	return func() {
		if !c.funcC4[t] {
			return
		}
		if c.rtExpiry[t] > c.Core.F {
			c.AddTask(c.c4TickFunc(t), "childe-c4-ticker", 240) //check every 4 sec
		} else {
			//riptide expired
			c.funcC4[t] = false
			return
		}
		if c.Base.Cons >= 4 {
			if c.mlBurstUsed {
				c.mlBurstUsed = false
				return
			}
		}

		//All of Riptide effects triggered by C4 are considered Normal Attack DMG.
		if c.Core.Status.Duration("childemelee") > 0 {
			c.Core.Log.Debugw(fmt.Sprintf("C4 %v ticking", "Riptide Slash"), "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("childemelee"), "rt", c.rtExpiry[t], "target", t, "c4", c.funcC4[t])
			//riptide slash
			d := c.Snapshot(
				"C4 Riptide Slash",
				core.AttackTagNormal,
				core.ICDTagNone,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				25,
				rtSlash[c.TalentLvlSkill()],
			)
			d.Targets = core.TargetAll
			c.QueueDmg(&d, 1)
		} else {
			//riptide flash
			c.Core.Log.Debugw(fmt.Sprintf("C4 %v ticking", "Riptide Flash"), "frame", c.Core.F, "event", core.LogCharacterEvent, "dur",
				c.Core.Status.Duration("childemelee"), "rt", c.rtExpiry[t], "target", t, "c4", c.funcC4[t])
			d := c.Snapshot(
				"C4 Riptide Flash",
				core.AttackTagNormal,
				core.ICDTagTartagliaRiptideFlash,
				core.ICDGroupDefault,
				core.StrikeTypeDefault,
				core.Hydro,
				25,
				rtFlash[0][c.TalentLvlAttack()],
			)
			d.Targets = core.TargetAll
			for i := 1; i <= 3; i++ {
				x := d.Clone()
				c.QueueDmg(&x, i)
			}
		}
	}
}

func (c *char) applyRT(t int) {
	//apply riptide status to enemies hit
	if c.rtExpiry[t] < c.Core.F {
		c.Core.Log.Debugw("Childe applied riptide", "frame", c.Core.F, "event", core.LogCharacterEvent, "target", t, "Expiry", c.Core.F+c.rtA1)
	}
	c.rtExpiry[t] = c.Core.F + c.rtA1

	if c.Base.Cons >= 4 && !c.funcC4[t] {
		c.funcC4[t] = true
		c.AddTask(c.c4TickFunc(t), "childe-c4-tick", 6) //tick procs every 4 sec
	}
}
