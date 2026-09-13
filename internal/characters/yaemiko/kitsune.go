package yaemiko

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const kitsuneDur = 900

type kitsune struct {
	src         int
	deleted     bool
	kitsuneArea info.AttackPattern
}

func (c *char) makeKitsune() {
	k := &kitsune{}
	k.src = c.Core.F
	k.deleted = false

	// spawn kitsune detection area on player pos
	k.kitsuneArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, c.kitsuneDetectionRadius)

	// start ticking
	c.Core.Tasks.Add(c.kitsuneTick(k), 120-skillStart)
	// add task to delete this one if times out (and not deleted by anything else)
	c.Core.Tasks.Add(func() {
		// i think we can just check for .deleted here
		if k.deleted {
			return
		}
		// ok now we can delete this
		c.popOldestKitsune()
	}, kitsuneDur+c.revelationBonusSkillDur()-skillStart) // e ani + duration

	if c.kitsuneCount() == 0 {
		c.Core.Status.Add(yaeTotemStatus, kitsuneDur+c.revelationBonusSkillDur()-skillStart)
	}
	// pop oldest first
	if c.kitsuneCount() == 3 {
		c.popOldestKitsune()
		c.a1OnSkillPopKitsune()
	}
	c.kitsunes = append(c.kitsunes, k)
	c.SetTag(yaeTotemCount, c.kitsuneCount())
}

func (c *char) popAllKitsune() {
	for i := range c.kitsunes {
		c.kitsunes[i].deleted = true
	}
	c.kitsunes = c.kitsunes[:0]
	c.Core.Status.Delete(yaeTotemStatus)
	c.SetTag(yaeTotemCount, 0)
}

func (c *char) popOldestKitsune() {
	if len(c.kitsunes) == 0 {
		// nothing to pop??
		return
	}

	c.kitsunes[0].deleted = true
	c.kitsunes = c.kitsunes[1:]

	// here check for status
	if len(c.kitsunes) > 0 {
		dur := c.Core.F - c.kitsunes[0].src + (kitsuneDur + c.revelationBonusSkillDur() - skillStart)
		if dur < 0 {
			log.Panicf("oldest totem should have expired already? dur: %v totem: %v", dur, *c.kitsunes[0])
		}
		c.Core.Status.Add(yaeTotemStatus, dur)
	} else {
		c.Core.Status.Delete(yaeTotemStatus)
	}

	c.SetTag(yaeTotemCount, len(c.kitsunes))
}

func (c *char) kitsuneBurst(ai info.AttackInfo, pattern info.AttackPattern) {
	for i := 0; i < c.kitsuneCount(); i++ {
		c.Core.QueueAttack(ai, pattern, burstThunderbolt1Hitmark+i*24, burstThunderbolt1Hitmark+i*24)
		if c.Base.Cons >= 1 {
			c.Core.Tasks.Add(func() {
				c.AddEnergy("yae-c1", 8)
			}, burstThunderbolt1Hitmark+i*24)
		}
		c.a1()
		c.Core.Log.NewEvent("sky kitsune thunderbolt", glog.LogCharacterEvent, c.Index()).
			Write("src", c.kitsunes[i].src).
			Write("delay", burstThunderbolt1Hitmark+i*24)
	}

	if !c.revelation {
		c.popAllKitsune()
	}
}

func (c *char) kitsuneTick(totem *kitsune) func() {
	return func() {
		// if deleted do nothing
		if totem.deleted {
			return
		}

		// spawn 1 attack
		// priority: enemy > gadget
		tick := func(pos info.Point) {
			// c6
			// Sesshou Sakura start at Level 2 when created. Max level increased to 4, and their attacks will ignore 45% of the opponents' DEF.

			lvl := c.sakuraLevel()
			// safety check
			if lvl < 1 {
				panic("sakura level should not be < 1 during tick")
			}

			c.Core.Log.NewEvent("sky kitsune tick at level", glog.LogCharacterEvent, c.Index()).
				Write("sakura level", lvl)

			flatDmg, revCB := c.revelationEnhanceDMG()

			ai := info.AttackInfo{
				Abil:             "Sesshou Sakura Tick",
				ActorIndex:       c.Index(),
				AttackTag:        attacks.AttackTagElementalArt,
				Mult:             skill[lvl-1][c.TalentLvlSkill()],
				ICDTag:           attacks.ICDTagElementalArt,
				ICDGroup:         attacks.ICDGroupDefault,
				StrikeType:       attacks.StrikeTypeDefault,
				Element:          attributes.Electro,
				Durability:       25,
				FlatDmg:          flatDmg,
				IgnoreDefPercent: c.c6DefIgnore(),
			}

			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTarget(pos, nil, 0.5),
				1,
				1,
				c.particleCB,
				revCB,
				c.c4MakeCB(),
			)
		}

		// try to target an enemy first
		enemy := c.Core.Combat.RandomEnemyWithinArea(totem.kitsuneArea, nil)
		if enemy != nil {
			tick(enemy.Pos())
		} else {
			// target gadget if no enemy was targeted
			gadget := c.Core.Combat.RandomGadgetWithinArea(totem.kitsuneArea, nil)
			if gadget != nil {
				tick(gadget.Pos())
			}
		}

		// tick per ~2.9s seconds
		c.Core.Tasks.Add(c.kitsuneTick(totem), 176)
	}
}

func (c *char) kitsuneCount() int {
	return len(c.kitsunes)
}

func (c *char) sakuraLevel() int {
	count := c.kitsuneCount()
	if count <= 0 {
		// this is for the base case when there are no totems (other wise we'll end up with 1 if C2)
		return 0
	}
	if count > 3 {
		panic("wtf more than 3 totems")
	}
	return count + c.c2SakuraLevelBonus()
}
