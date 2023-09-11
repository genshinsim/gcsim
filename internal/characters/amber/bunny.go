package amber

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const manualExplosionAbil = "Baron Bunny (Manual Explosion)"

type Bunny struct {
	*gadget.Gadget
	*reactable.Reactable
	ae   *combat.AttackEvent
	char *char
}

func (b *Bunny) HandleAttack(atk *combat.AttackEvent) float64 {
	b.Core.Events.Emit(event.OnGadgetHit, b, atk)

	b.Core.Log.NewEvent(fmt.Sprintf("baron bunny hit by %s", atk.Info.Abil), glog.LogCharacterEvent, b.char.Index)

	b.ShatterCheck(atk)

	//TODO: Check if sucrose E or faruzan E on Bunny is 25 dur or 50 dur

	if atk.Info.Durability > 0 {
		atk.Info.Durability *= reactions.Durability(b.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
		if atk.Info.Durability > 0 && atk.Info.Element != attributes.Physical {
			existing := b.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			b.React(atk)
			if b.Core.Flags.LogDebug && atk.Reacted {
				b.Core.Log.NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
				).
					Write("attack_tag", atk.Info.AttackTag).
					Write("applied_ele", atk.Info.Element.String()).
					Write("dur", applied).
					Write("abil", atk.Info.Abil).
					Write("target", "Bunny").
					Write("existing", existing).
					Write("after", b.Reactable.ActiveAuraString())

			}
		}
	}

	//apply damage delay is only there to make sure aura gets applied at the end of current frame
	//however because we can only hold cryo, we'll only call this if atk is cryo and there
	//is durability left
	if atk.Info.Element != attributes.Cryo {
		return 0
	}
	if atk.Info.Durability < reactable.ZeroDur {
		return 0
	}
	if atk.Reacted {
		return 0
	}
	b.Core.Combat.Tasks.Add(func() {
		b.attachEle(atk)
	}, 0)
	return 0
}

func (b *Bunny) attachEle(atk *combat.AttackEvent) {
	// check for ICD first
	existing := b.Reactable.ActiveAuraString()
	applied := atk.Info.Durability
	b.AttachOrRefill(atk)
	if b.Core.Flags.LogDebug {
		b.Core.Log.NewEvent(
			"application",
			glog.LogElementEvent,
			atk.Info.ActorIndex,
		).
			Write("attack_tag", atk.Info.AttackTag).
			Write("applied_ele", atk.Info.Element.String()).
			Write("dur", applied).
			Write("abil", atk.Info.Abil).
			Write("target", "Bunny").
			Write("existing", existing).
			Write("after", b.Reactable.ActiveAuraString())

	}
}

func (b *Bunny) React(a *combat.AttackEvent) {
	//only check the ones possible
	switch a.Info.Element {
	case attributes.Electro:
		b.TryFrozenSuperconduct(a)
		b.TrySuperconduct(a)
	case attributes.Pyro:
		b.TryMelt(a)
	// Cryo cannot react because the only allowed aura is Cryo.
	// case attributes.Cryo:
	case attributes.Hydro:
		b.TryFreeze(a)
	case attributes.Anemo:
		b.TrySwirlCryo(a)
		b.TrySwirlFrozen(a)
	case attributes.Geo:
		b.TryCrystallizeCryo(a)
		b.TryCrystallizeFrozen(a)
	case attributes.Dendro:
	}
}

func (b *Bunny) Tick() {
	//this is needed since gadget tick
	b.Gadget.Tick()
	b.Reactable.Tick()
}

func (c *char) makeBunny() *Bunny {

	b := &Bunny{}

	// Bunny is offset 1.3-1.5m in the Y direction for Tap E.
	// TODO: Implement hold E for different distances
	// TODO: Implement collision check for moving Baron Bunny off enemies and players
	player := c.Core.Combat.Player()
	bunnyPos := geometry.CalcOffsetPoint(
		player.Pos(),
		geometry.Point{Y: 1.4},
		player.Direction(),
	)
	b.Gadget = gadget.New(c.Core, bunnyPos, 1, combat.GadgetTypBaronBunny)
	b.Reactable = &reactable.Reactable{}
	b.Reactable.Init(b, c.Core)

	b.Gadget.Duration = 484

	b.char = c

	ai := combat.AttackInfo{
		Abil:       "Baron Bunny",
		ActorIndex: c.Index,
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       bunnyExplode[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	b.ae = &combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	}
	b.ae.Callbacks = append(b.ae.Callbacks, c.makeParticleCB())
	c.bunnies = append(c.bunnies, b)
	b.Gadget.OnKill = b.explode
	b.Gadget.OnExpiry = b.explode
	c.Core.Combat.AddGadget(b)
	return b
}

func (b *Bunny) explode() {
	// Explode
	b.char.Core.Log.NewEvent("amber exploding bunny", glog.LogCharacterEvent, b.char.Index).
		Write("src", b.Gadget.Src())
	b.char.Core.QueueAttackEvent(b.ae, 1)

	// remove self from list of bunnies
	for i := 0; i < len(b.char.bunnies); i++ {
		if b.char.bunnies[i] == b {
			b.char.bunnies = append(b.char.bunnies[:i], b.char.bunnies[i+1:]...)
		}
	}
}

func (c *char) makeParticleCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Pyro, c.ParticleDelay)
	}
}

func (c *char) manualExplode() {
	//do nothing if there are no bunnies
	if len(c.bunnies) == 0 {
		c.Core.Log.NewEvent("Did not find any Bunnies", glog.LogCharacterEvent, c.Index)
		return
	}
	//only explode the first bunny
	c.bunnies[0].ae.Info.Abil = manualExplosionAbil
	c.bunnies[0].Kill()

}

// explode all bunnies on overload
func (c *char) overloadExplode() {
	c.Core.Events.Subscribe(event.OnOverload, func(args ...interface{}) bool {
		target := args[0].(*enemy.Enemy)
		atk := args[1].(*combat.AttackEvent)
		if len(c.bunnies) == 0 {
			return false
		}
		//TODO: only amber trigger?
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}

		for _, v := range c.bunnies {
			// OL is a radius 3 hit on the target

			if v.IsWithinArea(combat.NewCircleHitOnTarget(target.Pos(), nil, 3)) {
				v.ae.Info.Abil = manualExplosionAbil
				v.Kill()
			}
		}

		return false
	}, "amber-bunny-explode-overload")
}
