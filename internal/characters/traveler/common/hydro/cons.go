package hydro

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const c4ICDKey = "travelerhydro-c4-icd"

// When using Aquacrest Saber, an Aquacrest Aegis that can absorb 10% of the Traveler's Max HP in DMG will be
// created and will absorb Hydro DMG with 250% effectiveness. It will persist until the Traveler finishes using the skill.
// Once every 2s, after a Dewdrop hits an opponent, if the Traveler is being protected by Aquacrest Aegis,
// the DMG Absorption of the Aegis will be restored to 10% of the Traveler's Max HP. If the Traveler is not presently
// being protected by an Aegis, one will be redeployed.
func (c *char) c4() {
	existingShield := c.Core.Player.Shields.Get(shield.TravelerHydroC4)
	if existingShield != nil {
		// update hp
		shd, _ := existingShield.(*shield.Tmpl)
		shd.HP = 0.1 * c.MaxHP()
		c.Core.Log.NewEvent("update shield hp", glog.LogCharacterEvent, c.Index).
			Write("hp", shd.HP)
		return
	}

	// add shield
	c.Core.Player.Shields.Add(&shield.Tmpl{
		Src:        c.Core.F,
		ShieldType: shield.TravelerHydroC4,
		Name:       "Traveler (Hydro) C4",
		HP:         0.1 * c.MaxHP(),
		Ele:        attributes.Hydro,
		Expires:    c.Core.F + 15*60,
	})
}

func (c *char) makeC4CB() combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if c.StatusIsActive(c4ICDKey) {
			return
		}

		c.c4()
		c.AddStatus(c4ICDKey, 2*60, true)
	}
}

func (c *char) c4Remove() {
	if c.Base.Cons < 4 {
		return
	}

	existingShield := c.Core.Player.Shields.Get(shield.TravelerHydroC4)
	if existingShield == nil {
		return
	}
	shd, _ := existingShield.(*shield.Tmpl)
	shd.Expires = c.Core.F + 1
}

// When the Traveler picks up a Sourcewater Droplet, they will restore HP to a nearby party member with the lowest
// remaining HP percentage based on 6% of said member's Max HP.
func (c *char) c6() {
	lowest := c.Index
	chars := c.Core.Player.Chars()
	for i, char := range chars {
		if char.CurrentHPRatio() < chars[lowest].CurrentHPRatio() {
			lowest = i
		}
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  lowest,
		Type:    player.HealTypePercent,
		Message: "Tides of Justice",
		Src:     0.06,
		Bonus:   c.Stat(attributes.Heal),
	})
}
