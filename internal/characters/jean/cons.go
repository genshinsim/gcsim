package jean

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C1:
// Increases the pulling speed of Gale Blade after holding for more than 1s, and increases the DMG dealt by 40%.
func (c *char) c1(snap *combat.Snapshot) {
	// add 40% dmg
	snap.Stats[attributes.DmgP] += .4
	c.Core.Log.NewEvent("jean c1 adding 40% dmg", glog.LogCharacterEvent, c.Index).
		Write("final dmg%", snap.Stats[attributes.DmgP])
}

// C2:
// When Jean picks up an Elemental Orb/Particle, all party members have their Movement SPD and ATK SPD increased by 15% for 15s.
func (c *char) c2() {
	c.c2buff = make([]float64, attributes.EndStatType)
	c.c2buff[attributes.AtkSpd] = 0.15
	c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
		// only trigger if Jean catches the particle
		if c.Core.Player.Active() != c.Index {
			return false
		}
		// apply C2 to all characters
		for _, this := range c.Core.Player.Chars() {
			this.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("jean-c2", 900),
				AffectedStat: attributes.AtkSpd,
				Amount: func() ([]float64, bool) {
					return c.c2buff, true
				},
			})
		}
		return false
	}, "jean-c2")
}

// C4:
// Within the Field created by Dandelion Breeze, all opponents have their Anemo RES decreased by 40%.
func (c *char) c4() {
	// gets called once right before burst start and then at the same time as heal ticks (every 1s)
	// add debuff to all targets for 1.2 s
	enemies := c.Core.Combat.EnemiesWithinArea(c.burstArea, nil)
	for _, e := range enemies {
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("jean-c4", 72), // 1.2s
			Ele:   attributes.Anemo,
			Value: -0.4,
		})
	}
}

// C6:
// Incoming DMG is decreased by 35% within the Field created by Dandelion Breeze.
// Upon leaving the Dandelion Field, this effect lasts for 3 attacks or 10s.
func (c *char) c6() {
	c.Core.Log.NewEvent("jean-c6 not implemented", glog.LogCharacterEvent, c.Index)
}
