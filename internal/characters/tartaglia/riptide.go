package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

//While aiming, the power of Hydro will accumulate on the arrowhead.
//A arrow fully charged with the torrent will deal Hydro DMG and apply the Riptide status.
func (c *char) aimedApplyRiptide(a combat.AttackCB) {
	c.applyRiptide("aimed shot", a.Target)
}

//Swiftly fires a Hydro-imbued magic arrow, dealing AoE Hydro DMG and applying the Riptide status.
func (c *char) rangedBurstApplyRiptide(a combat.AttackCB) {
	c.applyRiptide("ranged burst", a.Target)
}

//When Tartaglia is in Foul Legacy: Raging Tide's Melee Stance, on dealing a CRIT hit,
//Normal and Charged Attacks apply the Riptide status effect to opponents.
func (c *char) meleeApplyRiptide(a combat.AttackCB) {
	//only applies if is crit
	if a.IsCrit {
		c.applyRiptide("melee", a.Target)
	}
}

func (c *char) applyRiptide(src string, t combat.Target) {
	if c.Base.Cons >= 4 && t.GetTag(riptideKey) < c.Core.F {
		c.Core.Tasks.Add(func() { c.rtC4Tick(t) }, 60*4)
	}

	t.SetTag(riptideKey, c.Core.F+riptideDuration)
	c.Core.Log.NewEvent(
		fmt.Sprintf("riptide applied (%v)", src),
		glog.LogCharacterEvent,
		c.Index,
		"target", t.Index(),
		"expiry", c.Core.F+riptideDuration,
	)
}

// if tartaglia is in melee stance, triggers Riptide Slash against opponents on the field affected by Riptide every 4s, otherwise, triggers Riptide Flash.
// this constellation effect is not subject to ICD.
func (c *char) rtC4Tick(t combat.Target) {
	if t.GetTag(riptideKey) < c.Core.F {
		return
	}

	if c.Core.Status.Duration("tartagliamelee") > 0 {
		c.rtSlashTick(t)
	} else {
		c.rtFlashTick(t)
	}

	c.Core.Tasks.Add(func() { c.rtC4Tick(t) }, 60*4)
}

// Riptide Flash: A fully-charged Aimed Shot that hits an opponent affected
// by Riptide deals consecutive bouts of AoE DMG. Can occur once every 0.7s.
func (c *char) rtFlashCallback(a combat.AttackCB) {
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

	c.rtFlashTick(a.Target)
}

func (c *char) rtFlashTick(t combat.Target) {
	//queue damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Flash",
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagTartagliaRiptideFlash,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       rtFlash[c.TalentLvlAttack()],
	}

	//proc 3 hits
	x, y := t.Shape().Pos()
	for i := 1; i <= 3; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHit(x, y, 0.5, false, combat.TargettableEnemy), 1, 1)
	}

	c.Core.Log.NewEvent(
		"riptide flash triggered",
		glog.LogCharacterEvent,
		c.Index,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", t.Index(),
		"riptide_flash_icd", t.GetTag(riptideFlashICDKey),
		"riptide_expiry", t.GetTag(riptideKey),
	)

	//queue particles
	if c.rtParticleICD < c.Core.F {
		c.rtParticleICD = c.Core.F + 180 //3 sec
		c.Core.QueueParticle("tartaglia", 1, attributes.Hydro, 100)
	}
}

//Hitting an opponent affected by Riptide with a melee attack unleashes a Riptide Slash that deals AoE Hydro DMG.
//DMG dealt in this way is considered Elemental Skill DMG, and can only occur once every 1.5s.
func (c *char) rtSlashCallback(a combat.AttackCB) {
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

	c.rtSlashTick(a.Target)
}

func (c *char) rtSlashTick(t combat.Target) {
	//trigger attack
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Slash",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       rtSlash[c.TalentLvlSkill()],
	}

	x, y := t.Shape().Pos()
	c.Core.QueueAttack(ai, combat.NewCircleHit(x, y, 2, false, combat.TargettableEnemy), 1, 1)

	c.Core.Log.NewEvent(
		"riptide slash ticked",
		glog.LogCharacterEvent,
		c.Index,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", t.Index(),
		"riptide_slash_icd", t.GetTag(riptideSlashICDKey),
		"riptide_expiry", t.GetTag(riptideKey),
	)

	//queue particle if not on icd
	if c.rtParticleICD < c.Core.F {
		c.rtParticleICD = c.Core.F + 180 //3 sec
		c.Core.QueueParticle("tartaglia", 1, attributes.Hydro, 100)
	}
}

//When the obliterating waters hit an opponent affected by Riptide, it clears their Riptide status
//and triggers a Hydro Explosion that deals AoE Hydro DMG. DMG dealt in this way is considered Elemental Burst DMG.
func (c *char) rtBlastCallback(a combat.AttackCB) {
	//only triggers if target affected by riptide
	if a.Target.GetTag(riptideKey) < c.Core.F {
		return
	}
	//TODO: this shares icd with slash???
	if a.Target.GetTag(riptideSlashICDKey) > c.Core.F {
		return
	}
	//queue damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Blast",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       rtBlast[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewDefCircHit(3, false, combat.TargettableEnemy), 1, 1)

	c.Core.Log.NewEvent(
		"riptide blast triggered",
		glog.LogCharacterEvent,
		c.Index,
		"dur", c.Core.Status.Duration("tartagliamelee"),
		"target", a.Target.Index(),
		"rtExpiry", a.Target.GetTag(riptideKey),
	)

	//clear riptide status
	a.Target.RemoveTag(riptideKey)
}

//Riptide Burst: Defeating an opponent affected by Riptide creates a Hydro burst
//that inflicts the Riptide status on nearby opponents hit.
// Handles Childe riptide burst and C2 on death effects
func (c *char) onDefeatTargets() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		t := args[0].(combat.Target)
		//do nothing if no riptide on target
		if t.GetTag(riptideKey) < c.Core.F {
			return false
		}
		c.Core.Tasks.Add(func() {
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Riptide Burst",
				AttackTag:  combat.AttackTagNormal,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				StrikeType: combat.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 50,
				Mult:       rtBurst[c.TalentLvlAttack()],
			}
			c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 0)
		}, 5)
		//TODO: re-index riptide expiry frame array if needed
		if c.Base.Cons >= 2 {
			c.AddEnergy("tartaglia-c2", 4)
		}
		return false
	}, "tartaglia-on-enemy-death")
}
