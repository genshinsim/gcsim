package chasca

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const c6key = "chasca-c6"
const c6IcdKey = "chasca-c6-icd"

func (c *char) c1() float64 {
	if c.Base.Cons < 1 {
		return 0.0
	}
	return 0.333
}

func (c *char) c1Conversion() {
	if c.Base.Cons < 1 {
		return
	}
	if c.bulletsNext[2] == attributes.Anemo {
		return
	}
	c.bulletsNext[1] = c.partyPHECTypesUnique[c.Core.Rand.Intn(len(c.partyPHECTypesUnique))]
}

func (c *char) c2A1Stack() int {
	if c.Base.Cons < 2 {
		return 0
	}
	return 1
}

func (c *char) c2cb(src int) combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	return func(ac combat.AttackCB) {
		if c.c2Src == src {
			return
		}
		c.c2Src = src

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Shining Shadowhunt Shell (C2)",
			AttackTag:      attacks.AttackTagExtra,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        ac.AttackEvent.Info.Element,
			Durability:     25,
			Mult:           4,
		}
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)
		c.Core.QueueAttack(ai, ap, 0, 1)
	}
}

func (c *char) c4cb(src int) combat.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}
	return func(ac combat.AttackCB) {
		c.AddEnergy("chasca-c4", 1.5)
		if c.c4Src == src {
			return
		}
		c.c4Src = src

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Radiant Shadowhunt Shell (C4)",
			AttackTag:      attacks.AttackTagExtra,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        ac.AttackEvent.Info.Element,
			Durability:     25,
			Mult:           4,
		}
		// TODO: get the actual target range
		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)
		c.Core.QueueAttack(ai, ap, 0, 1)
	}
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.StatusIsActive(c6IcdKey) {
		return
	}
	c.AddStatus(c6IcdKey, 3*60, true)
	c.AddStatus(c6key, 3*60, true)
}

func (c *char) c6buff(snap *combat.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}

	if !c.StatusIsActive(c6key) {
		return
	}
	old := snap.Stats[attributes.CD]
	snap.Stats[attributes.CD] += 1.20
	c.Core.Log.NewEvent("c6 adding crit dmg", glog.LogCharacterEvent, c.Index).
		Write("old", old).
		Write("new", snap.Stats[attributes.CD])
}

func (c *char) removeC6() {
	if c.Base.Cons < 6 {
		return
	}
	c.DeleteStatus(c6key)
}

func (c *char) c6ChargeTime(count int) int {
	if c.Base.Cons < 6 {
		return cumuSkillAimLoadFrames[count-1]
	}
	if c.StatusIsActive(c6key) {
		return cumuSkillAimLoadFramesC6Instant[count-1]
	}
	return cumuSkillAimLoadFramesC6[count-1]
}
