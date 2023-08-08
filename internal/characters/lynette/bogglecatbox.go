package lynette

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

type BogglecatBox struct {
	*gadget.Gadget
	char        *char
	pos         geometry.Point
	vividTravel int
}

func (c *char) newBogglecatBox(vividTravel int) *BogglecatBox {
	b := &BogglecatBox{}

	player := c.Core.Combat.Player()
	b.pos = geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{Y: 1.8},
		player.Direction(),
	)

	// TODO: double check estimation of hitbox
	b.Gadget = gadget.New(c.Core, b.pos, 1, combat.GadgetTypBogglecatBox)
	b.char = c
	b.vividTravel = vividTravel

	b.Duration = burstDuration // TODO: proper frames
	b.char.AddStatus(burstKey, b.Duration, false)

	c.Core.Tasks.Add(func() {
		// initial hit
		initialAI := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Magic Trick: Astonishing Shift",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(initialAI, combat.NewCircleHitOnTarget(player, geometry.Point{Y: 1.5}, 4.5), 0, 0)

		// bogglecat ticks
		bogglecatAI := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Bogglecat Box",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       bogglecat[c.TalentLvlBurst()],
		}
		// TODO: double check tick count and interval
		// queue up ticks
		for t := 1 * 60; t <= b.Duration; t += 1 * 60 {
			c.Core.QueueAttack(bogglecatAI, combat.NewCircleHitOnTarget(b.pos, nil, 6), t, t)
		}
	}, burstHitmark-burstSpawn)

	// check for absorb every 0.3s besides absorbing on being hit
	b.OnThinkInterval = b.absorbCheck
	b.ThinkInterval = 0.3 * 60

	b.Core.Log.NewEvent("Lynette Bogglecat Box added", glog.LogCharacterEvent, c.Index).Write("src", b.Src())

	return b
}

func (b *BogglecatBox) HandleAttack(atk *combat.AttackEvent) float64 {
	b.Core.Events.Emit(event.OnGadgetHit, b, atk)

	b.Core.Log.NewEvent(fmt.Sprintf("lynette bogglecat box hit by %s", atk.Info.Abil), glog.LogCharacterEvent, b.char.Index)

	if atk.Info.Durability <= 0 {
		return 0
	}
	atk.Info.Durability *= reactions.Durability(b.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
	if atk.Info.Durability <= 0 {
		return 0
	}

	// only allow contact with cryo/pyro/hydro/electro
	switch atk.Info.Element {
	case attributes.Cryo:
	case attributes.Pyro:
	case attributes.Hydro:
	case attributes.Electro:
	default:
		return 0
	}

	b.absorbRoutine(atk.Info.Element)

	return 0
}

func (b *BogglecatBox) absorbRoutine(absorbedElement attributes.Element) {
	b.Core.Log.NewEvent(fmt.Sprintf("lynette bogglecat box came into contact with %s", absorbedElement.String()), glog.LogCharacterEvent, b.char.Index)

	// vivid shots
	vividShotAI := combat.AttackInfo{
		ActorIndex: b.char.Index,
		Abil:       "Vivid Shot",
		AttackTag:  attacks.AttackTagElementalBurst,
		// should be ElementalBurstMix, but it just needs to be different from all the other icd tags used by the char so no need to add extra icd tag
		ICDTag:     attacks.ICDTagElementalBurstAnemo,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    absorbedElement, // take element from absorbed element
		Durability: 25,
		Mult:       vivid[b.char.TalentLvlBurst()],
	}
	// queue up vivid shots
	// TODO: proper interval frames
	interval := int(2.3 * 60)
	for t := interval; t <= b.Duration; t += interval {
		b.Core.Tasks.Add(func() {
			// target random enemy within 15m of burst pos
			enemy := b.Core.Combat.RandomEnemyWithinArea(combat.NewCircleHitOnTarget(b.pos, nil, 15), nil)
			// queue up 1 or 2 (C2) vivid shots
			for i := 0; i < b.char.vividCount; i++ {
				// TODO: snapshot correct here?
				b.Core.QueueAttack(vividShotAI, combat.NewCircleHitOnTarget(enemy, nil, 1), 0, b.vividTravel)
			}
		}, t)
	}

	// apply A4
	b.char.a4(b.Duration)

	// remove the gadget, because it should not be hitable after contact
	b.Kill()
}

func (b *BogglecatBox) absorbCheck() {
	absorbedElement := b.Core.Combat.AbsorbCheck(combat.NewCircleHitOnTarget(b.pos, nil, 0.48), attributes.Cryo, attributes.Pyro, attributes.Hydro, attributes.Electro)
	if absorbedElement == attributes.NoElement {
		return
	}
	b.absorbRoutine(absorbedElement)
}

func (b *BogglecatBox) SetDirection(trg geometry.Point) {}
func (b *BogglecatBox) SetDirectionToClosestEnemy()     {}
func (b *BogglecatBox) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
