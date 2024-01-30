package chevreuse

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey    = "chev-c1-icd"
	c2ICDKey    = "chev-c2-icd"
	c4StatusKey = "chev-c4"
)

// When the active character with the "Coordinated Tactics" status (not including Chevreuse herself)
// triggers the Overloaded reaction, they will recover 6 Energy.
// This effect can be triggered once every 10s.
// You must first unlock the Passive Talent "Vanguard's Coordinated Tactics."
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	if !c.onlyPyroElectro {
		return
	}

	c.Core.Events.Subscribe(event.OnOverload, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		// does not include chevreuse
		if atk.Info.ActorIndex == c.Index {
			return false
		}
		// does not trigger off-field
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		// does not trigger if on cd
		if c.StatusIsActive(c1ICDKey) {
			return false
		}
		c.AddStatus(c1ICDKey, 10*60, true)

		active := c.Core.Player.ByIndex(atk.Info.ActorIndex)
		active.AddEnergy("chev-c1", 6)

		return false
	}, "chev-c1")
}

// After Holding Short-Range Rapid Interdiction Fire and hitting a target,
// 2 chain explosions will be triggered near the location where said target is hit.
// Each explosion deals Pyro DMG equal to 120% of Chevreuse's ATK.
// This effect can be triggered up to once every 10s,
// and DMG dealt this way is considered Elemental Skill DMG.
func (c *char) c2() combat.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}
	// triggers on hitting anything, not just enemy
	return func(a combat.AttackCB) {
		if c.StatusIsActive(c2ICDKey) {
			return
		}
		c.AddStatus(c2ICDKey, 10*60, true)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Sniper Induced Explosion (C2)",
			AttackTag:  attacks.AttackTagElementalArt,
			// should be ElementalArtExtra but no other chevreuse attack shares this tag so this is ok
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   25,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       1.2,
		}

		// calc pos
		player := c.Core.Combat.Player()
		targetPos := a.Target.Pos()
		// TODO: using player direction here is inaccurate since it should maybe be the direction from which it was hit?
		bomb1Pos := geometry.CalcOffsetPoint(targetPos, geometry.Point{X: -1.5}, player.Direction())
		bomb2Pos := geometry.CalcOffsetPoint(targetPos, geometry.Point{X: 1.5}, player.Direction())

		// calc delay
		// random between 0.6s and 1s from hit
		// not shared between bombs
		bomb1Delay := int(60 * (0.6 + c.Core.Rand.Float64()*(1-0.6)))
		bomb2Delay := int(60 * (0.6 + c.Core.Rand.Float64()*(1-0.6)))

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(bomb1Pos, nil, 3),
			bomb1Delay,
			bomb1Delay,
		)

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(bomb2Pos, nil, 3),
			bomb2Delay,
			bomb2Delay,
		)
	}
}

// After using Ring of Bursting Grenades, the Hold mode of Short-Range Rapid Interdiction Fire
// will not go on cooldown when Chevreuse uses it.
// This effect is removed after Short-Range Rapid Interdiction Fire
// has been fired twice using Hold or after 6s.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	c.AddStatus(c4StatusKey, 6*60, true)
	c.c4ShotsLeft = 2
}

// After 12s of the healing effect from Short-Range Rapid Interdiction Fire,
// all nearby party members recover HP equivalent to 10% of Chevreuse's Max HP once.
func (c *char) c6TeamHeal() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6HealQueued = false

	for _, char := range c.Core.Player.Chars() {
		c.c6(char)
	}

	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "In Pursuit of Ending Evil (C6)",
		Src:     0.1 * c.MaxHP(),
		Bonus:   c.Stat(attributes.Heal),
	})
}

func (c *char) c6(char *character.CharWrapper) {
	if c.Base.Cons < 6 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.PyroP] = 0.20
	m[attributes.ElectroP] = 0.20

	char.AddStatMod(character.StatMod{
		Base: modifier.NewBaseWithHitlag(fmt.Sprintf("chev-c6-%v-stack", c.c6StackCounts[char.Index]+1), 8*60),
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	c.c6StackCounts[char.Index] = (c.c6StackCounts[char.Index] + 1) % 3
}
