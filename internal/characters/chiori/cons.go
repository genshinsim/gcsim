package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
)

const (
	c4ICD     = "chiori-c4-icd"
	c4Key     = "chiori-c4"
	c4Lockout = "chiori-c4-lockout" // tracking 15s window
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
func (c *char) c1Active() bool {
	if c.Base.Cons < 1 {
		return false
	}
	if c.geoCount < 2 {
		return false
	}
	return true
}

// For 10s after using Hiyoku: Twin Blades, a simplified automaton doll, "Kinu,"
// will be summoned next to your active character every 3s. Kinu will attack
// nearby opponents, dealing AoE Geo DMG equivalent to 170% of Tamoto's DMG. DMG
// dealt this way is considered Elemental Skill DMG.
//
// Kinu will leave the field after 1 attack or after lasting 3s.
func (c *char) c2init() {
	if c.Base.Cons < 2 {
		return
	}
	for _, v := range c.Core.Player.Chars() {
		if v.Base.Element == attributes.Geo {
			c.geoCount++
		}
	}
}
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.kill(c.c2Ticker)
	t := newTicker(c.Core, 600)
	t.cb = c.createKinu
	t.interval = 180 // every 3 sec
	//TODO: no idea when first one comes out. so let's just do 1 frame for now
	c.Core.Tasks.Add(t.tick, 1)
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
func (c *char) c4init() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if !c.StatusIsActive(c4Key) {
			return false
		}
		if c.Core.Status.Duration(c4ICD) > 0 {
			return false
		}
		if c.c4count >= 3 {
			return false
		}
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}

		c.c4count++
		c.Core.Status.Add(c4ICD, 60) //TODO: icd/
		c.createKinu()

		return false
	}, "chiori-c4")
}

func (c *char) c4() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 4 {
		return
	}
	if c.StatusIsActive(c4Lockout) {
		return
	}
	c.c4count = 0
	c.DeleteStatus(c4ICD)

	c.AddStatus(c4Key, 8*60, true) //TODO: hitlag?
	// lock out c4 from being triggered again for 15s
	c.AddStatus(c4Lockout, 15*60, true) //TODO: hitlag
}

// After triggering a follow-up effect of the Passive Talent "Tailor-Made,"
// Chiori's own Fluttering Hasode's CD is decreased by 12s. Must unlock the
// Passive "Tailor-Made" first.
//
// In addition, the DMG dealt by Chiori's own Normal Attacks is increased by an
// amount equal to 235% of her own DEF.
func (c *char) c6() {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 6 {
		return
	}
	c.ReduceActionCooldown(action.ActionSkill, 12*60)
}

func (c *char) c6AttackMod(ai *combat.AttackInfo, snap *combat.Snapshot) {
	if c.Base.Ascension < 1 {
		return
	}
	if c.Base.Cons < 6 {
		return
	}
	ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
	ai.FlatDmg *= ai.Mult
}
