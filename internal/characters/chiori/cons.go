package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	c2Duration        = 10 * 60
	c2SpawnInterval   = 3 * 60
	c2MinRandom       = 0.8
	c2MaxRandom       = 1.8
	c2CenterOffset    = 0.2
	c4Key             = "chiori-c4"
	c4Duration        = 8 * 60
	c4Lockout         = "chiori-c4-lockout"
	c4LockoutDuration = 15 * 60
	c4ICDKey          = "chiori-c4-icd"
	c4ICD             = 1 * 60
	c4AttackLimit     = 3
	c4MinRandom       = 1.8
	c4MaxRandom       = 2.8
	c4CenterOffset    = 0
)

// The AoE of the automaton doll "Tamoto" summoned by Fluttering Hasode is
// increased by 50%. Additionally, if there is a Geo party member other than
// Chiori, Fluttering Hasode will trigger the following after the dash is
// completed:
// - Summon an additional Tamoto. Only one additional Tamoto can exist at the
// same time, whether summoned by Chiori this way or through the presence of a
// Geo Construct.
// - Triggers the Passive Talent "The Finishing Touch." This effect requires you
// to first unlock the Passive Talent "The Finishing Touch."
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	geoCount := 0
	for _, v := range c.Core.Player.Chars() {
		if v.Base.Element == attributes.Geo {
			geoCount++
		}
		if geoCount >= 2 {
			c.c1Active = true
			break
		}
	}
	// 50% from description most likely refers to the volume of the AoE
	// volume of a cylinder is pi*r^2*h, so radius needs to be multiplied by 1.5^2=2.25
	c.skillSearchAoE *= 2.25
}

// For 10s after using Hiyoku: Twin Blades, a simplified automaton doll, "Kinu,"
// will be summoned next to your active character every 3s. Kinu will attack
// nearby opponents, dealing AoE Geo DMG equivalent to 170% of Tamoto's DMG. DMG
// dealt this way is considered Elemental Skill DMG.
//
// Kinu will leave the field after 1 attack or after lasting 3s.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Log.NewEvent("c2 activated", glog.LogCharacterEvent, c.Index)

	// kill existing c2 ticker
	c.kill(c.c2Ticker)

	// spawn new c2 ticker
	// yes, the c2 duration and spawn interval is hitlag affected
	t := newTicker(c.Core, c2Duration, c.QueueCharTask)
	t.cb = c.createKinu(c.Core.F, c2CenterOffset, c2MinRandom, c2MaxRandom)
	t.interval = c2SpawnInterval
	c.QueueCharTask(t.tick, c2SpawnInterval)
	c.c2Ticker = t
}

// For 8s after triggering either follow-up effect of the Passive Talent
// "Tailor-Made," when your current active character's Normal, Charged, or
// Plunging Attacks hit a nearby opponent, a simplified automaton doll, "Kinu,"
// will be summoned near this opponent. You can summon 1 Kinu every 1s in this
// way, and up to 3 Kinu may be summoned this way during each instance of
// "Tailor-Made"'s Seize the Moment or Tailoring effect. The above effect can be
// triggered up to once every 15s.
//
// Must unlock the Passive Talent "Tailor-Made" first.
func (c *char) c4Activation() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4Lockout) {
		return
	}

	c.Core.Log.NewEvent("c4 activated", glog.LogCharacterEvent, c.Index)

	c.AddStatus(c4Lockout, c4LockoutDuration, true) // applied to chiori

	c.c4AttackCount = 0
	c.DeleteStatus(c4ICDKey)
	c.AddStatus(c4Key, c4Duration, false) // applied to team, so not hitlag affected
}

func (c *char) c4() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		// c4 status not active
		if !c.StatusIsActive(c4Key) {
			return false
		}
		// c4 attack on icd
		if c.StatusIsActive(c4ICDKey) {
			return false
		}
		// attack not na/ca/plunge
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}
		// atk not by active char
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}

		// apply icd
		c.AddStatus(c4ICDKey, c4ICD, true) // applied to chiori

		c.Core.Log.NewEvent("c4 spawning kinu", glog.LogCharacterEvent, c.Index)

		// spawn kinu
		c.createKinu(c.Core.F, c4CenterOffset, c4MinRandom, c4MaxRandom)()

		// increment attack count and delete c4 if reached limit
		c.c4AttackCount++
		if c.c4AttackCount == c4AttackLimit {
			c.DeleteStatus(c4Key)
		}

		return false
	}, "chiori-c4")
}

// After triggering a follow-up effect of the Passive Talent "Tailor-Made,"
// Chiori's own Fluttering Hasode's CD is decreased by 12s. Must unlock the
// Passive "Tailor-Made" first.
func (c *char) c6CooldownReduction() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 6 {
		return
	}
	c.ReduceActionCooldown(action.ActionSkill, 12*60)
}

// In addition, the DMG dealt by Chiori's own Normal Attacks is increased by an
// amount equal to 235% of her own DEF.
func (c *char) c6NAIncrease(ai *combat.AttackInfo, snap *combat.Snapshot) {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 6 {
		return
	}
	ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
	ai.FlatDmg *= 2.35
}
