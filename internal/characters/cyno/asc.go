package cyno

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "cyno-a1"

// When Cyno is in the Pactsworn Pathclearer state activated by Sacred Rite: Wolf's Swiftness,
// Cyno will enter the Endseer stance at intervals. If he activates Secret Rite: Chasmic Soulfarer whle affected by this stance,
// he will activate the Judication effect, increasing the DMG of this Secret Rite: Chasmic Soulfarer by 35%,
// and firing off 3 Duststalker Bolts that deal 50% of Cyno's ATK as Electro DMG.
// Duststalker Bolt DMG is considered Elemental Skill DMG.
//
// - checks for ascension level in burst.go to avoid queuing this up only to fail the ascension level check
//
// - other parts of a1 depend on the status that this task applies so they don't need an extra ascension check
func (c *char) a1() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.a1Extended = false
	c.AddStatus(a1Key, 84, true)
	c.QueueCharTask(c.a1, 234)
}

func (c *char) a1Buff() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.35
	// game also implements dmg buff with 1s modifier
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("cyno-a1-dmg", 60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			// actual game uses AttackTagElementalArtExtra for a1, this is a decent
			// workaround
			if atk.Info.Abil != skillBName {
				return nil, false
			}
			return m, true
		},
	})
}

// If Cyno dashes with the a1 modifier, he will increase the modifier's
// durability by 20. This translates to a 0.28s extension.
func (c *char) a1Extension() {
	c.Core.Events.Subscribe(event.OnDash, func(_ ...interface{}) bool {
		if c.a1Extended {
			return false
		}
		active := c.Core.Player.ActiveChar()
		if !(active.Index == c.Index && active.StatusIsActive(a1Key)) {
			return false
		}
		c.ExtendStatus(a1Key, 17)
		c.a1Extended = true
		c.Core.Log.NewEvent("a1 dash pp slide", glog.LogCharacterEvent, c.Index).
			Write("expiry", c.StatusExpiry(a1Key))
		return false
	}, "cyno-a1-dash")
}

// Cyno's DMG values will be increased based on his Elemental Mastery as follows:
//
// - Pactsworn Pathclearer's Normal Attack DMG is increased by 150% of his Elemental Mastery.
func (c *char) a4NormalAttack() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 1.5
}

// Cyno's DMG values will be increased based on his Elemental Mastery as follows:
//
// - Duststalker Bolt DMG from his Passive Talent Featherfall Judgment is increased by 250% of his Elemental Mastery.
func (c *char) a4Bolt() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 2.5
}
