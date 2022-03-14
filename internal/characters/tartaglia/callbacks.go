package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

//While aiming, the power of Hydro will accumulate on the arrowhead.
//A arrow fully charged with the torrent will deal Hydro DMG and apply the Riptide status.
func (c *char) aimedApplyRiptide(a core.AttackCB) {
	c.applyRiptide("aimed shot", a.Target)
}

//Swiftly fires a Hydro-imbued magic arrow, dealing AoE Hydro DMG and applying the Riptide status.
func (c *char) rangedBurstApplyRiptide(a core.AttackCB) {
	c.applyRiptide("ranged burst", a.Target)
}

//When Tartaglia is in Foul Legacy: Raging Tide's Melee Stance, on dealing a CRIT hit,
//Normal and Charged Attacks apply the Riptide status effect to opponents.
func (c *char) meleeApplyRiptide(a core.AttackCB) {
	//only applies if is crit
	if a.IsCrit {
		c.applyRiptide("melee", a.Target)
	}
}

func (c *char) applyRiptide(src string, t coretype.Target) {
	if c.Base.Cons >= 4 && t.GetTag(riptideKey) < c.Core.Frame {
		c.AddTask(func() { c.c4(t) }, "tartaglia-c4", 60*4)
	}

	t.SetTag(riptideKey, c.Core.Frame+riptideDuration)
	c.coretype.Log.NewEvent(
		fmt.Sprintf("riptide applied (%v)", src),
		coretype.LogCharacterEvent,
		c.Index,
		"target", t.Index(),
		"expiry", c.Core.Frame+riptideDuration,
	)
}

// Riptide Flash: A fully-charged Aimed Shot that hits an opponent affected
// by Riptide deals consecutive bouts of AoE DMG. Can occur once every 0.7s.
func (c *char) rtFlashCallback(a core.AttackCB) {
	//do nothing if no riptide on target
	if a.Target.GetTag(riptideKey) < c.Core.Frame {
		return
	}
	//do nothing if flash still on icd
	if a.Target.GetTag(riptideFlashICDKey) > c.Core.Frame {
		return
	}
	//add 0.7s icd
	a.Target.SetTag(riptideFlashICDKey, c.Core.Frame+42)

	c.rtFlashTick(a.Target)
}

func (c *char) rtFlashTick(t coretype.Target) {
	//queue damage
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Flash",
		AttackTag:  coretype.AttackTagNormal,
		ICDTag:     core.ICDTagTartagliaRiptideFlash,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.Hydro,
		Durability: 25,
		Mult:       rtFlash[c.TalentLvlAttack()],
	}

	//proc 3 hits
	x, y := t.Shape().Pos()
	for i := 1; i <= 3; i++ {
		c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 0.5, false, core.TargettableEnemy), 1, 1)
	}

	c.coretype.Log.NewEvent(
		"riptide flash triggered",
		coretype.LogCharacterEvent,
		c.Index,
		"dur", c.Core.StatusDuration("tartagliamelee"),
		"target", t.Index(),
		"riptide_flash_icd", t.GetTag(riptideFlashICDKey),
		"riptide_expiry", t.GetTag(riptideKey),
	)

	//queue particles
	if c.rtParticleICD < c.Core.Frame {
		c.rtParticleICD = c.Core.Frame + 180 //3 sec
		c.QueueParticle("tartaglia", 1, core.Hydro, 100)
	}
}

//Hitting an opponent affected by Riptide with a melee attack unleashes a Riptide Slash that deals AoE Hydro DMG.
//DMG dealt in this way is considered Elemental Skill DMG, and can only occur once every 1.5s.
func (c *char) rtSlashCallback(a core.AttackCB) {
	//do nothing if no riptide on target
	if a.Target.GetTag(riptideKey) < c.Core.Frame {
		return
	}
	//do nothing if slash still on icd
	if a.Target.GetTag(riptideSlashICDKey) > c.Core.Frame {
		return
	}
	//add 1.5s icd
	a.Target.SetTag(riptideSlashICDKey, c.Core.Frame+90)

	c.rtSlashTick(a.Target)
}

func (c *char) rtSlashTick(t coretype.Target) {
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

	x, y := t.Shape().Pos()
	c.Core.Combat.QueueAttack(ai, core.NewCircleHit(x, y, 2, false, coretype.TargettableEnemy), 1, 1)

	c.coretype.Log.NewEvent(
		"riptide slash ticked",
		coretype.LogCharacterEvent,
		c.Index,
		"dur", c.Core.StatusDuration("tartagliamelee"),
		"target", t.Index(),
		"riptide_slash_icd", t.GetTag(riptideSlashICDKey),
		"riptide_expiry", t.GetTag(riptideKey),
	)

	//queue particle if not on icd
	if c.rtParticleICD < c.Core.Frame {
		c.rtParticleICD = c.Core.Frame + 180 //3 sec
		c.QueueParticle("tartaglia", 1, core.Hydro, 100)
	}
}

//When the obliterating waters hit an opponent affected by Riptide, it clears their Riptide status
//and triggers a Hydro Explosion that deals AoE Hydro DMG. DMG dealt in this way is considered Elemental Burst DMG.
func (c *char) rtBlastCallback(a core.AttackCB) {
	//only triggers if target affected by riptide
	if a.Target.GetTag(riptideKey) < c.Core.Frame {
		return
	}
	//TODO: this shares icd with slash???
	if a.Target.GetTag(riptideSlashICDKey) > c.Core.Frame {
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

	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(3, false, coretype.TargettableEnemy), 0, 1)

	c.coretype.Log.NewEvent(
		"riptide blast triggered",
		coretype.LogCharacterEvent,
		c.Index,
		"dur", c.Core.StatusDuration("tartagliamelee"),
		"target", a.Target.Index(),
		"rtExpiry", a.Target.GetTag(riptideKey),
	)

	//clear riptide status
	a.Target.RemoveTag(riptideKey)
}
