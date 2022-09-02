package cyno

//When Cyno is in the Pactsworn Pathclearer state activated by Sacred Rite: Wolf's Swiftness,
//Cyno will enter the Endseer stance at intervals. If he activates Secret Rite: Chasmic Soulfarer whle affected by this stance,
//he will activate the Judication effect, increasing the DMG of this Secret Rite: Chasmic Soulfarer by 35%,
//and firing off 3 Duststalker Bolts that deal 50% of Cyno's ATK as Electro DMG.
//Duststalker Bolt DMG is considered Elemental Skill DMG.

// mucho texto
func (c *char) a1() {
	if !c.StatusIsActive(burstKey) { //if not on burst gtfo
		return
	}

	c.AddStatus(a4key, 78, true) //TODO: idk if this gets affected by hitlag, koli do ur job, also idk duration either
	c.Core.Tasks.Add(func() {
		c.a1()
	}, 284) //TODO: 284 interval between subsequent endseers (source: i made it the fuck up)
}

// Cyno's DMG values will be increased based on his Elemental Mastery as follows:
// Pactsworn Pathclearer's Normal Attack DMG is increased by 100% of his Elemental Mastery.
// Duststalker Bolt DMG from his Ascension Talent Featherfall Judgment is increased by 250% of his Elemental Mastery.
func (c *char) a4() {
	//I just added this at flat damage on the attack frames lol
}
