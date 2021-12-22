package tartaglia

import "github.com/genshinsim/gcsim/pkg/core"

//While aiming, the power of Hydro will accumulate on the arrowhead.
//A arrow fully charged with the torrent will deal Hydro DMG and apply the Riptide status.
func (c *char) aimedApplyRiptide(a core.AttackCB) {
	a.Target.SetTag(riptideKey, c.Core.F+riptideDuration)
	c.Core.Log.Debugw(
		"riptide applied (CA)",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"char", c.Index,
		"target", a.Target.Index(),
		"expiry", c.Core.F+riptideDuration,
	)
}

//When Tartaglia is in Foul Legacy: Raging Tide's Melee Stance, on dealing a CRIT hit,
//Normal and Charged Attacks apply the Riptide status effect to opponents.
func (c *char) meleeApplyRiptide(a core.AttackCB) {
	//only applies if is crit
	if a.IsCrit {
		a.Target.SetTag(riptideKey, c.Core.F+riptideDuration)
		c.Core.Log.Debugw(
			"riptide applied (melee)",
			"frame", c.Core.F,
			"event", core.LogCharacterEvent,
			"char", c.Index,
			"target", a.Target.Index(),
			"expiry", c.Core.F+riptideDuration,
		)
	}
}

// Riptide Flash: A fully-charged Aimed Shot that hits an opponent affected
// by Riptide deals consecutive bouts of AoE DMG. Can occur once every 0.7s.
func (c *char) rtFlashCallback(a core.AttackCB) {
	//do nothing if no riptide on target
	if a.Target.GetTag(riptideKey) < c.Core.F {
		return
	}
	//do nothing if flash still on icd
	if a.Target.GetTag(riptideFlashICDKey) > c.Core.F {
		return
	}
	//add 0.7s icd
	a.Target.SetTag(riptideFlashICDKey, c.Core.F+42)
	//queue damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Flash",
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagTartagliaRiptideFlash,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       rtFlash[0][c.TalentLvlAttack()],
	}

	//proc 3 hits
	for i := 1; i <= 3; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewDefSingleTarget(1, core.TargettableEnemy), 0, 1)
	}
	c.Core.Log.Debugw(
		"riptide flash triggered",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", a.Target.Index(),
		"riptide_flash_icd", a.Target.GetTag(riptideFlashICDKey),
		"riptide_expiry", a.Target.GetTag(riptideKey),
	)

	//queue particles
	if c.rtParticleICD < c.Core.F {
		c.rtParticleICD = c.Core.F + 180 //3 sec
		c.QueueParticle("tartaglia", 1, core.Hydro, 100)
	}
}

//Hitting an opponent affected by Riptide with a melee attack unleashes a Riptide Slash that deals AoE Hydro DMG.
//DMG dealt in this way is considered Elemental Skill DMG, and can only occur once every 1.5s.
func (c *char) rtSlashCallback(a core.AttackCB) {
	//do nothing if no riptide on target
	if a.Target.GetTag(riptideKey) < c.Core.F {
		return
	}
	//do nothing if slash still on icd
	if a.Target.GetTag(riptideSlashICDKey) > c.Core.F {
		return
	}
	//add 1.5s icd
	a.Target.SetTag(riptideSlashICDKey, c.Core.F+90)

	//trigger attack
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Slash",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       rtSlash[c.TalentLvlSkill()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 0, 0)

	c.Core.Log.Debugw(
		"Riptide Slash ticked",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", a.Target.Index(),
		"riptide_slash_icd", a.Target.GetTag(riptideSlashICDKey),
		"riptide_expiry", a.Target.GetTag(riptideKey),
	)

	//queue particle if not on icd
	if c.rtParticleICD < c.Core.F {
		c.rtParticleICD = c.Core.F + 180 //3 sec
		c.QueueParticle("tartaglia", 1, core.Hydro, 100)
	}
}

//When the obliterating waters hit an opponent affected by Riptide, it clears their Riptide status
//and triggers a Hydro Explosion that deals AoE Hydro DMG. DMG dealt in this way is considered Elemental Burst DMG.
func (c *char) rtBlastCallback(a core.AttackCB) {
	//only triggers if target affected by riptide
	if a.Target.GetTag(riptideKey) < c.Core.F {
		return
	}
	//TODO: this shares icd with slash???
	if a.Target.GetTag(riptideSlashICDKey) > c.Core.F {
		return
	}
	//queue damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Blast",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 50,
		Mult:       rtBlast[c.TalentLvlBurst()],
	}

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, core.TargettableEnemy), 0, 1)

	c.Core.Log.Debugw(
		"riptide blast triggered",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", a.Target.Index(),
		"rtExpiry", a.Target.GetTag(riptideKey),
	)

	//clear riptide status
	a.Target.RemoveTag(riptideKey)
}
