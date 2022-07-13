package heizou

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

//When Vacuum Slugger hits opponents affected by Hydro/Pyro/Cryo/Electro,
//these opponents will be afflicted with Windmuster Iris.
//This Windmuster Iris will explode after a moment and dissipate,
//dealing AoE DMG of the corresponding aforementioned elemental type.
func (c *char) irisDmg(src string, t core.Target) {

	//TODO: does burst iris snapshot
	aiAbs := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Windmuster Iris",
		AttackTag:  core.AttackTagElementalBurst,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeDefault,
		Element:    core.NoElement,
		Durability: 25,
		Mult:       burstIris[c.TalentLvlBurst()],
	}
	//snapAbs := c.Snapshot(&aiAbs)
	//t.SetTag("iris", c.Core.F+20)
	c.Core.Log.NewEvent(
		"Iris Applied",
		core.LogCharacterEvent,
		c.Index,
		"target", t.Index(),
		"expiry", c.Core.F+20,
	)
	//TODO: Iris tiiming
	x, y := t.Shape().Pos()

	switch t.AuraType() {
	case core.Pyro, core.Hydro, core.Electro, core.Cryo:
		aiAbs.Element = t.AuraType()
		c.Core.Combat.QueueAttack(aiAbs, core.NewCircleHit(x, y, 0.5, false, core.TargettableEnemy), 1, 1)
	case core.EC:
		aiAbs.Element = core.Hydro
		c.Core.Combat.QueueAttack(aiAbs, core.NewCircleHit(x, y, 0.5, false, core.TargettableEnemy), 1, 1)
	case core.Frozen:
		aiAbs.Element = core.Cryo
		c.Core.Combat.QueueAttack(aiAbs, core.NewCircleHit(x, y, 0.5, false, core.TargettableEnemy), 1, 1)
	default:
		c.Core.Log.NewEvent(
			"No valid aura detected, omiting iris",
			core.LogCharacterEvent,
			c.Index,
			"aura type", t.AuraType(),
		)
	}

}
