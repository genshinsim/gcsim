package amber

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const manualExplosionAbil = "Baron Bunny (Manual Explosion)"

type Bunny struct {
	*gadget.Gadget
	*reactable.Reactable
	ae   *combat.AttackEvent
	char *char
	src  int
}

func (b *Bunny) AuraContains(e ...attributes.Element) bool {
	for ele := range e {
		if b.Reactable.Durability[ele] <= reactable.ZeroDur {
			return false
		}
	}
	return true
}

func (b *Bunny) HandleAttack(atk *combat.AttackEvent) float64 {
	b.Core.Events.Emit(event.OnGadgetHit, b, atk)

	b.Core.Log.NewEvent(fmt.Sprintf("baron bunny hit by %s", atk.Info.Abil), glog.LogCharacterEvent, b.char.Index)

	b.ShatterCheck(atk)

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
					Write("target", b.Key()).
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

func (s *Bunny) attachEle(atk *combat.AttackEvent) {
	// check for ICD first
	existing := s.Reactable.ActiveAuraString()
	applied := atk.Info.Durability
	s.AttachOrRefill(atk)
	if s.Core.Flags.LogDebug {
		s.Core.Log.NewEvent(
			"application",
			glog.LogElementEvent,
			atk.Info.ActorIndex,
		).
			Write("attack_tag", atk.Info.AttackTag).
			Write("applied_ele", atk.Info.Element.String()).
			Write("dur", applied).
			Write("abil", atk.Info.Abil).
			Write("target", s.Key()).
			Write("existing", existing).
			Write("after", s.Reactable.ActiveAuraString())

	}
}

func (r *Bunny) React(a *combat.AttackEvent) {
	//only check the ones possible
	switch a.Info.Element {
	case attributes.Electro:
		r.TryFrozenSuperconduct(a)
		r.TrySuperconduct(a)
	case attributes.Pyro:
		r.TryMelt(a)
	case attributes.Cryo:
	case attributes.Hydro:
		r.TryFreeze(a)
	case attributes.Anemo:
		r.TrySwirlCryo(a)
		r.TrySwirlFrozen(a)
	case attributes.Geo:
		r.TryCrystallizeCryo(a)
		r.TryCrystallizeFrozen(a)
	case attributes.Dendro:
	}
}

func (s *Bunny) Tick() {
	//this is needed since gadget tick
	s.Reactable.Tick()
	s.Gadget.Tick()
}

func (c *char) makeBunny() *Bunny {

	b := &Bunny{}

	// TODO: I think it's supposed by default be an offset from the player, and not on the player?
	// I haven't played amber in years
	// player := c.Core.Combat.Player()
	// bunnyPos := geometry.CalcOffsetPoint(
	// 	player.Pos(),
	// 	geometry.Point{Y: 1},
	// 	player.Direction(),
	// )
	bunnyPos := c.Core.Combat.Player().Pos()
	b.Gadget = gadget.New(c.Core, bunnyPos, 1, combat.GadgetTypBaronBunny)
	b.Reactable = &reactable.Reactable{}
	b.Reactable.Init(b, c.Core)

	// duration is 8.2s
	b.Duration = 492

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
	b.Gadget.OnKill = b.OnKill
	return b
}

func (b *Bunny) OnKill() {
	// Explode
	b.explode()

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

func (b *Bunny) explode() {
	b.char.Core.Log.NewEvent("amber exploding bunny", glog.LogCharacterEvent, b.char.Index).
		Write("src", b.src)
	b.char.Core.QueueAttackEvent(b.ae, 1)
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
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {

		atk := args[1].(*combat.AttackEvent)
		if len(c.bunnies) == 0 {
			return false
		}
		//TODO: only amber trigger?
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if atk.Info.AttackTag != attacks.AttackTagOverloadDamage {
			return false
		}
		c.bunnies[0].ae.Info.Abil = manualExplosionAbil

		for _, v := range c.bunnies {
			v.Kill()
		}

		return false
	}, "amber-bunny-explode-overload")
}
