package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	// A1
	// window in which a1 can be used after skill
	a1WindowKey = "chiori-a1-window"
	// Tapestry - coordinated attacks
	a1SeizeTheMomentKey         = "chiori-seize-the-moment"
	a1SeizeTheMomentDuration    = 8 * 60
	a1SeizeTheMomentICDKey      = "chiori-seize-the-moment-icd"
	a1SeizeTheMomentICD         = 2 * 60
	a1SeizeTheMomentAttackLimit = 2
	// Tailoring - geo infusion
	a1GeoInfusionKey      = "chiori-tailoring"
	a1GeoInfusionDuration = 5 * 60
	// A4
	a4BuffKey  = "chiori-a4"
	a4Duration = 20 * 60
)

// Gain different effects depending on the next action you take within a short
// duration after using Fluttering Hasode's upward sweep. If you Press the
// Elemental Skill, you will trigger the Tapestry effect. If you (Press/Tap) your Normal
// Attack, the Tailoring effect will be triggered instead.
//
// Tapestry
// - Switches to the next character in your roster.
// - Grants all your party members "Seize the Moment": When your active party
// member's Normal Attacks, Charged Attacks, and Plunging Attacks hit a nearby
// opponent, "Tamoto" will execute a coordinated attack, dealing 100% of
// Fluttering Hasode's upward sweep DMG as AoE Geo DMG at the opponent's
// location. DMG dealt this way is considered Elemental Skill DMG.
// - "Seize the Moment" lasts 8s, and 1 of "Tamoto"'s coordinated attack can be
// unleashed every 2s. 2 such coordinated attacks can occur per "Seize the
// Moment" effect duration.
//
// Tailoring
// - Chiori gains Geo infusion for 5s.
//
// When on the field, if Chiori does not either Press her Elemental Skill or use
// a Normal Attack within a short time after using Fluttering Hasode's upward
// sweep, the Tailoring effect will be triggered by default.
func (c *char) a1TapestrySetup() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		// seize the moment not active
		if !c.StatusIsActive(a1SeizeTheMomentKey) {
			return false
		}
		// seize the moment on icd
		if c.StatusIsActive(a1SeizeTheMomentICDKey) {
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
		// atk not within 30m of player
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.IsWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 30)) {
			return false
		}

		// apply icd
		c.AddStatus(a1SeizeTheMomentICDKey, a1SeizeTheMomentICD, true)

		// deal dmg
		ai := combat.AttackInfo{
			Abil:       "Fluttering Hasode (Seize the Moment)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagChioriSkill,
			ICDGroup:   attacks.ICDGroupChioriSkill,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       thrustAtkScaling[c.TalentLvlSkill()],
		}
		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= thrustDefScaling[c.TalentLvlSkill()]
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(t, nil, 2.5), 0)

		// increment attack count and delete seize the moment if reached limit
		c.a1AttackCount++
		if c.a1AttackCount == a1SeizeTheMomentAttackLimit {
			c.DeleteStatus(a1SeizeTheMomentKey)
		}

		return false
	}, a1SeizeTheMomentKey)
}

// tap and hold skill have different start and duration for a1 window
// so a1 activation is based on the supplied values
func (c *char) activateA1Window(start, duration int) {
	if c.Base.Ascension < 1 {
		return
	}
	c.QueueCharTask(func() {
		c.AddStatus(a1WindowKey, duration, true)
		// When on the field, if Chiori does not either Press her Elemental Skill or use
		// a Normal Attack within a short time after using Fluttering Hasode's upward
		// sweep, the Tailoring effect will be triggered by default.
		c.a1Triggered = false
		c.QueueCharTask(func() {
			if c.a1Triggered {
				return
			}
			c.a1Tailoring()
		}, duration)
	}, start)
}

func (c *char) commonA1Trigger() {
	c.a1Triggered = true
	// she can't use skill again if she triggers this
	c.DeleteStatus(a1WindowKey)
	c.c4Activation()
	c.c6CooldownReduction()
}

func (c *char) a1Tapestry() {
	c.commonA1Trigger()

	c.Core.Log.NewEvent("a1 tapestry triggered", glog.LogCharacterEvent, c.Index)
	c.AddStatus(a1SeizeTheMomentKey, a1SeizeTheMomentDuration, true)
	c.a1AttackCount = 0
}

func (c *char) a1Tailoring() {
	c.commonA1Trigger()

	c.Core.Log.NewEvent("a1 tailoring triggered", glog.LogCharacterEvent, c.Index)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		a1GeoInfusionKey,
		attributes.Geo,
		a1GeoInfusionDuration,
		true,
		attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
	)
}

// tailoring proc via na can fail if it was already triggered via a1 window expiring
func (c *char) tryTriggerA1TailoringNA() {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.StatusIsActive(a1WindowKey) {
		return
	}
	c.a1Tailoring()
}

// When a nearby party member creates a Geo Construct, Chiori will gain 20% Geo DMG Bonus for 20s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4Buff[attributes.GeoP] = 0.20
	c.Core.Events.Subscribe(event.OnConstructSpawned, func(args ...interface{}) bool {
		c.applyA4Buff()
		return false
	}, a4BuffKey)
}

// needs to be callable separately because of c1 rock doll activating a4
func (c *char) applyA4Buff() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a4BuffKey, a4Duration),
		AffectedStat: attributes.GeoP,
		Amount: func() ([]float64, bool) {
			return c.a4Buff, true
		},
	})
}
