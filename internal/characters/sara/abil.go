package sara

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other bow characters
// Has "travel" parameter, used to set the number of frames that the arrow is in the air (default = 10)
func (c *char) Attack(p map[string]int) (int, int) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}

	f, a := c.ActionFrames(core.ActionAttack, p)
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypePierce,
		Element:    core.Physical,
		Durability: 25,
		Mult:       attack[c.NormalCounter][c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	c.AdvanceNormalIndex()

	return f, a
}

// Aimed charge attack damage queue generator
// Additionally handles crowfeather state, E skill damage, and A4
// A4 effect is: When Tengu Juurai: Ambush hits opponents, Kujou Sara will restore 1.2 Energy to all party members for every 100% Energy Recharge she has. This effect can be triggered once every 3s.
// Has two parameters, "travel", used to set the number of frames that the arrow is in the air (default = 10)
// weak_point, used to determine if an arrow is hitting a weak point (default = 1 for true)
func (c *char) Aimed(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionAim, p)

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot, ok := p["weakspot"]

	ai := core.AttackInfo{
		ActorIndex:   c.Index,
		Abil:         "Aim Charge Attack",
		AttackTag:    core.AttackTagExtra,
		ICDTag:       core.ICDTagNone,
		ICDGroup:     core.ICDGroupDefault,
		StrikeType:   core.StrikeTypePierce,
		Element:      core.Electro,
		Durability:   25,
		Mult:         aimChargeFull[c.TalentLvlAttack()],
		HitWeakPoint: weakspot == 1,
	}
	// d.AnimationFrames = f
	c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), f, f+travel)

	// Cover state handling - drops crowfeather, which explodes after 1.5 seconds
	if c.Core.Status.Duration("saracover") > 0 {
		// Not sure what kind of strike type this is
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Electro,
			Durability: 25,
			Mult:       skill[c.TalentLvlSkill()],
		}

		//TODO: snapshot?
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, f+travel+90, c.a4)
		c.attackBuff(f + travel + 90)

		// Particles are emitted after the ambush thing hits
		c.QueueParticle("sara", 3, core.Electro, f+travel+90+100)

		c.Core.Status.DeleteStatus("saracover")
	}

	return f, a
}

// Implements skill handling. Fairly barebones since most of the actual stuff happens elsewhere
// Retreats rapidly with the speed of a tengu, summoning the protection of the Crowfeather. Gains Crowfeather Cover for 18s, and when Kujou Sara fires a fully-charged Aimed Shot, Crowfeather Cover will be consumed, and will leave a Crowfeather at the target location.
// Crowfeathers will trigger Tengu Juurai: Ambush after a short time, dealing Electro DMG and granting the active character within its AoE an ATK Bonus based on Kujou Sara's Base ATK. The ATK Bonuses from different Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
// Also implements C2
// Unleashing Tengu Stormcall will leave a Weaker Crowfeather at Kujou Sara's original position that will deal 30% of its original DMG.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)

	// Snapshot for all of the crowfeathers are taken upon cast
	c.Core.Status.AddStatus("saracover", 18*60)

	// C2 handling
	if c.Base.Cons >= 2 {
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Tengu Juurai: Ambush C2",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			StrikeType: core.StrikeTypePierce,
			Element:    core.Electro,
			Durability: 25,
			Mult:       0.3 * skill[c.TalentLvlSkill()],
		}
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f, 90, c.a4)
		c.attackBuff(90)
	}

	c.SetCD(core.ActionSkill, 600)

	return f, a
}

// Implements A4 energy regen.
// According to library finding, text description is inaccurate
// it's more like for every 1% of ER, she grants 0.012 flat energy
func (c *char) a4(a core.AttackCB) {
	if c.Core.F < c.a4LastProc {
		return
	}
	c.a4LastProc = c.Core.F + 180
	energyAddAmt := 1.2 * (1 + c.Stat(core.ER))
	for _, char := range c.Core.Chars {
		char.AddEnergy("sara-a4", energyAddAmt)
	}
}

// Implements C1 CD reduction. Waits until delay (when it hits the enemy), then procs the effect
// Triggers on her E and Q
func (c *char) c1() {
	if c.Core.F < c.c1LastProc {
		return
	}
	c.c1LastProc = c.Core.F + 180
	c.ReduceActionCooldown(core.ActionSkill, 60)
	c.Core.Log.NewEvent("c1 reducing skill cooldown", core.LogCharacterEvent, c.Index, "new_cooldown", c.Cooldown(core.ActionSkill))
}

// Implements burst handling.
// Casts down Tengu Juurai: Titanbreaker, dealing AoE Electro DMG. Afterwards, Tengu Juurai: Titanbreaker spreads out into 4 consecutive bouts of Tengu Juurai: Stormcluster, dealing AoE Electro DMG.
// Tengu Juurai: Titanbreaker and Tengu Juurai: Stormcluster can provide the active character within their AoE with the same ATK Bonus as given by the Elemental Skill, Tengu Stormcall. The ATK Bonus provided by various kinds of Tengu Juurai will not stack, and their effects and duration will be determined by the last Tengu Juurai to take effect.
// Has parameters: "wave_cluster_hits", which controls how many of the mini-clusters in each wave hit an opponent.
// Also has "waveAttackProcs", used to determine which waves proc the attack buff.
// Format for both is a digit of length 5 - rightmost value is the starting proc (titanbreaker hit), and it moves from right to left
// For example, if you want waves 3 and 4 only to proc the attack buff, set attack_procs=11000
// For "wave_cluster_hits", use numbers in each slot to control the # of hits. So for center hit, then 3 hits from each wave, set wave_cluster_hits=33331
// Default for both is for the main titanbreaker and 1 wave to hit and also proc the buff
// Also implements C4
// The number of Tengu Juurai: Stormcluster released by Subjugation: Koukou Sendou is increased to 6.
func (c *char) Burst(p map[string]int) (int, int) {

	waveClusterHits, ok := p["wave_cluster_hits"]
	if !ok {
		waveClusterHits = 41
		if c.Base.Cons >= 2 {
			waveClusterHits = 61
		}
	}
	waveAttackProcs, ok := p["waveAttackProcs"]
	if !ok {
		waveAttackProcs = 11
	}

	f, a := c.ActionFrames(core.ActionBurst, p)

	// Entire burst snapshots sometime after activation but before 1st hit.
	// For now, assume that it snapshots on cd delay
	// Flagged as no ICD since the stormclusters do not share ICD with the main hit
	// No ICD should not functionally matter as this only hits once

	//titan breaker
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Tengu Juurai: Titanbreaker",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Electro,
		Durability: 25,
		Mult:       burstMain[c.TalentLvlBurst()],
	}
	// dTitanbreaker.Targets = core.TargetAll

	// dStormcluster.Targets = core.TargetAll

	var c1cb core.AttackCBFunc
	if c.Base.Cons >= 1 {
		c1cb = func(a core.AttackCB) { c.c1() }
	}

	if waveClusterHits%10 == 1 {
		// Actual hit procs after the full cast duration, or 50 frames
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 47, 50, c1cb)
	}
	if waveAttackProcs%10 == 1 {
		c.attackBuff(50)
	}

	//stormcluster
	// Each cluster wave hits ~50 frames after titanbreaker and each preceding wave
	// TODO: Replace with frame counts from KQM when those are available
	ai.Abil = "Tengu Juurai: Stormcluster"
	ai.ICDTag = core.ICDTagElementalBurst
	ai.Mult = burstCluster[c.TalentLvlBurst()]
	for waveN := 0; waveN < 4; waveN++ {
		// Handles the potential manual user override through the input tags
		// For each wave, get the corresponding digit from the numeric sequence (e.g. for 4441, wave 2 = 4)
		waveHits := int((waveClusterHits % PowInt(10, waveN+2)) / PowInt(10, waveN+2-1))
		waveAttackProc := int((waveAttackProcs % PowInt(10, waveN+2)) / PowInt(10, waveN+2-1))
		if waveHits > 0 {
			for j := 0; j < waveHits; j++ {
				c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(5, false, core.TargettableEnemy), 47, 100+18*waveN, c1cb)
			}
		}
		if waveAttackProc == 1 {
			c.attackBuff(100 + 18*waveN)
		}
	}

	c.SetCDWithDelay(core.ActionBurst, 20*60, 47)
	c.ConsumeEnergy(54)

	return f, a
}

// Handles attack boost from Sara's skills
// Checks for the onfield character at the delay frame, then applies buff to that character
// Also handles Sara C6
// The Electro DMG of characters who have had their ATK increased by Tengu Juurai has its Crit DMG increased by 60%.
// Uses event subscription as it can't get snapshotted
func (c *char) attackBuff(delay int) {
	c.AddTask(func() {
		buff := atkBuff[c.TalentLvlSkill()] * float64(c.Base.Atk+c.Weapon.Atk)

		active := c.Core.Chars[c.Core.ActiveChar]

		active.AddTag("sarabuff", c.Core.F+360)
		// c.Core.Status.AddStatus(fmt.Sprintf("sarabuff%v", active.Name()), 360)
		c.Core.Log.NewEvent("sara attack buff applied", core.LogCharacterEvent, c.Index, "char", active.CharIndex(), "buff", buff, "expiry", c.Core.F+360)

		m := make([]float64, core.EndStatType)
		m[core.ATK] = buff
		// AddMod function already only takes the most recent version of this buff
		active.AddMod(core.CharStatMod{
			Key: "sara-attack-buff",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: c.Core.F + 360,
		})

		if c.Base.Cons >= 1 {
			c.c1()
		}
		if c.Base.Cons >= 6 {
			c.c6(active)
		}
	}, "sara-attack-buff", delay)
}

func (c *char) c6(char core.Character) {
	val := make([]float64, core.EndStatType)
	val[core.CD] = 0.6
	char.AddPreDamageMod(core.PreDamageMod{
		Key:    "sara-c6",
		Expiry: c.Core.F + 360,
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.Element != core.Electro {
				return nil, false
			}
			return val, true
		},
	})
}

// Get integer power - required for burst
func PowInt(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}
