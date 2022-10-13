package cyno

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const a1Key = "cyno-a1"

// When Cyno is in the Pactsworn Pathclearer state activated by Sacred Rite: Wolf's Swiftness,
// Cyno will enter the Endseer stance at intervals. If he activates Secret Rite: Chasmic Soulfarer whle affected by this stance,
// he will activate the Judication effect, increasing the DMG of this Secret Rite: Chasmic Soulfarer by 35%,
// and firing off 3 Duststalker Bolts that deal 50% of Cyno's ATK as Electro DMG.
// Duststalker Bolt DMG is considered Elemental Skill DMG.
func (c *char) a1() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.a1Extended = false
	c.AddStatus(a1Key, 84, true)
	c.QueueCharTask(c.a1, 234)
}

// If cyno dashes with the a1 modifier, he will increase the modifier's
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
// Pactsworn Pathclearer's Normal Attack DMG is increased by 100% of his Elemental Mastery.
// Duststalker Bolt DMG from his Ascension Talent Featherfall Judgment is increased by 250% of his Elemental Mastery.
func (c *char) a4() {
	// I just added this at flat damage on the attack frames lol
}
