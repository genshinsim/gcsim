package cyno

// When Cyno is in the Pactsworn Pathclearer state activated by Sacred Rite: Wolf's Swiftness,
// Cyno will enter the Endseer stance at intervals. If he activates Secret Rite: Chasmic Soulfarer whle affected by this stance,
// he will activate the Judication effect, increasing the DMG of this Secret Rite: Chasmic Soulfarer by 35%,
// and firing off 3 Duststalker Bolts that deal 50% of Cyno's ATK as Electro DMG.
// Duststalker Bolt DMG is considered Elemental Skill DMG.
func (c *char) a1() {
	if !c.StatusIsActive(burstKey) {
		return
	}
	c.AddStatus(a1Key, 84, true)
	c.QueueCharTask(c.a1, 234)
}

// Cyno's DMG values will be increased based on his Elemental Mastery as follows:
// Pactsworn Pathclearer's Normal Attack DMG is increased by 100% of his Elemental Mastery.
// Duststalker Bolt DMG from his Ascension Talent Featherfall Judgment is increased by 250% of his Elemental Mastery.
func (c *char) a4() {
	// I just added this at flat damage on the attack frames lol
}
