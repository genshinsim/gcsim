package tartaglia

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

// While aiming, the power of Hydro will accumulate on the arrowhead.
// A arrow fully charged with the torrent will deal Hydro DMG and apply the Riptide status.
func (c *char) aimedApplyRiptide(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// apply if attack has hydro or e state is up and it's physical
	// covers:
	// - fully-charged aimed shot regardless of E state
	// - phys aimed shot in E state
	if a.AttackEvent.Info.Element == attributes.Hydro || (c.StatusIsActive(meleeKey) && a.AttackEvent.Info.Element == attributes.Physical) {
		c.applyRiptide("aimed shot", t)
	}
}

// Swiftly fires a Hydro-imbued magic arrow, dealing AoE Hydro DMG and applying the Riptide status.
func (c *char) rangedBurstApplyRiptide(a combat.AttackCB) {
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	c.applyRiptide("ranged burst", t)
}

func (c *char) applyRiptide(src string, t *enemy.Enemy) {
	if c.Base.Cons >= 4 && !t.StatusIsActive(riptideKey) {
		c.c4Src = c.Core.F
		t.QueueEnemyTask(c.rtC4Tick(c.Core.F, t), 60*3.9)
	}

	t.AddStatus(riptideKey, c.riptideDuration, true)
	c.Core.Log.NewEvent(
		fmt.Sprintf("riptide applied (%v)", src),
		glog.LogCharacterEvent,
		c.Index,
	).
		Write("target", t.Key()).
		Write("expiry", t.StatusExpiry(riptideKey))
}

// if tartaglia is in melee stance, triggers Riptide Slash against opponents on the field affected by Riptide every 4s, otherwise, triggers Riptide Flash.
// this constellation effect is not subject to ICD.
func (c *char) rtC4Tick(src int, t *enemy.Enemy) func() {
	return func() {
		if c.c4Src != src {
			c.Core.Log.NewEvent("tartaglia c4 src check ignored, src diff", glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.c4Src)
			return
		}
		if !t.StatusIsActive(riptideKey) {
			return
		}

		if c.StatusIsActive(meleeKey) {
			c.rtSlashTick(t)
		} else {
			c.rtFlashTick(t)
		}

		t.QueueEnemyTask(c.rtC4Tick(src, t), 60*3.9)
		c.Core.Log.NewEvent("tartaglia c4 applied", glog.LogCharacterEvent, c.Index).
			Write("src", src).
			Write("target", t.Key())
	}
}

// Riptide Flash: A fully-charged Aimed Shot that hits an opponent affected
// by Riptide deals consecutive bouts of AoE DMG. Can occur once every 0.7s.
func (c *char) rtFlashCallback(a combat.AttackCB) {
	// make sure it's actually an enemey
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// do nothing if attack does not have hydro
	// - prevent phys ca from triggering this
	if a.AttackEvent.Info.Element != attributes.Hydro {
		return
	}
	// do nothing if no riptide on target
	if !t.StatusIsActive(riptideKey) {
		return
	}
	// do nothing if flash still on icd
	if t.StatusIsActive(riptideFlashICDKey) {
		return
	}
	// add 0.7s icd
	t.AddStatus(riptideFlashICDKey, 42, true)

	c.rtFlashTick(t)
}

func (c *char) rtFlashTick(t *enemy.Enemy) {
	// queue damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Flash",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagTartagliaRiptideFlash,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       rtFlash[c.TalentLvlAttack()],
	}

	// proc 3 hits
	for i := 1; i <= 3; i++ {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3), 1, 1, c.particleCB)
	}

	c.Core.Log.NewEvent(
		"riptide flash triggered",
		glog.LogCharacterEvent,
		c.Index,
	).
		Write("dur", c.StatusExpiry(meleeKey)-c.Core.F).
		Write("target", t.Key()).
		Write("riptide_flash_icd", t.StatusExpiry(riptideFlashICDKey)).
		Write("riptide_expiry", t.StatusExpiry(riptideKey))
}

// Hitting an opponent affected by Riptide with a melee attack unleashes a Riptide Slash that deals AoE Hydro DMG.
// DMG dealt in this way is considered Elemental Skill DMG, and can only occur once every 1.5s.
func (c *char) rtSlashCallback(a combat.AttackCB) {
	// make sure it's actually an enemey
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// do nothing if E is not activated
	// - this can also be triggered by an aimed shot if E state is up when it hits the enemy
	if !c.StatusIsActive(meleeKey) {
		return
	}
	// do nothing if no riptide on target
	if !t.StatusIsActive(riptideKey) {
		return
	}
	// do nothing if slash still on icd
	if t.StatusIsActive(riptideSlashICDKey) {
		return
	}
	// add 1.5s icd
	t.AddStatus(riptideSlashICDKey, 90, true)

	c.rtSlashTick(t)
}

func (c *char) rtSlashTick(t *enemy.Enemy) {
	// trigger attack
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Slash",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       rtSlash[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 3), 1, 1, c.particleCB)

	c.Core.Log.NewEvent(
		"riptide slash ticked",
		glog.LogCharacterEvent,
		c.Index,
	).
		Write("dur", c.StatusExpiry(meleeKey)-c.Core.F).
		Write("target", t.Key()).
		Write("riptide_slash_icd", t.StatusExpiry(riptideSlashICDKey)).
		Write("riptide_expiry", t.StatusExpiry(riptideKey))
}

// When the obliterating waters hit an opponent affected by Riptide, it clears their Riptide status
// and triggers a Hydro Explosion that deals AoE Hydro DMG. DMG dealt in this way is considered Elemental Burst DMG.
func (c *char) rtBlastCallback(a combat.AttackCB) {
	// make sure it's actually an enemey
	t, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}
	// only triggers if target affected by riptide
	if !t.StatusIsActive(riptideKey) {
		return
	}
	// TODO: this shares icd with slash???
	if t.StatusIsActive(riptideSlashICDKey) {
		return
	}
	// queue damage
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Riptide Blast",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       rtBlast[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 5), 1, 1)

	c.Core.Log.NewEvent(
		"riptide blast triggered",
		glog.LogCharacterEvent,
		c.Index,
	).
		Write("dur", c.StatusExpiry(meleeKey)-c.Core.F).
		Write("target", t.Key()).
		Write("rtExpiry", t.StatusExpiry(riptideKey))

	// clear riptide status
	t.DeleteStatus(riptideKey)
}

// Riptide Burst: Defeating an opponent affected by Riptide creates a Hydro burst
// that inflicts the Riptide status on nearby opponents hit.
// Handles Childe riptide burst and C2 on death effects
func (c *char) onDefeatTargets() {
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		// do nothing if not an enemy
		if !ok {
			return false
		}
		// do nothing if no riptide on target
		if !t.StatusIsActive(riptideKey) {
			return false
		}
		c.Core.Tasks.Add(func() {
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Riptide Burst",
				AttackTag:  attacks.AttackTagNormal,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeSlash,
				Element:    attributes.Hydro,
				Durability: 50,
				Mult:       rtBurst[c.TalentLvlAttack()],
			}
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(t, nil, 5), 0, 0)
		}, 5)
		// TODO: re-index riptide expiry frame array if needed
		if c.Base.Cons >= 2 {
			c.AddEnergy("tartaglia-c2", 4)
		}
		return false
	}, "tartaglia-on-enemy-death")
}
