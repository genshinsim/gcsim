package chiori

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1TailorMadeWindowKey    = "chiori-a2-tailor-made"
	a1TailorMadeWindowLength = 120 //TODO: i made this up; should be from button press
	a1GeoInfusionKey         = "chiori-tailoring"
	a1SeizeTheMoment         = "chiori-tapestry"
	a1SeizeTheMomentICD      = "chiori-seize-the-moment-icd"

	a4BuffKey = "chiori-a4"
)

// Gain different effects depending on the next action you take within a short
// duration after using Fluttering Hasode's upward sweep. If you Press the
// Elemental Skill, you will trigger the Tapestry effect. If you your Normal
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
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if c.a1AttackCount >= 2 {
			return false
		}
		if c.Core.Status.Duration(a1SeizeTheMoment) == 0 {
			return false
		}
		if c.Core.Status.Duration(a1SeizeTheMomentICD) > 0 {
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

		c.a1AttackCount++
		c.Core.Status.Add(a1SeizeTheMomentICD, 120)

		ai := combat.AttackInfo{
			Abil:       "Fluttering Hasode (Seize the Moment)",
			ActorIndex: c.Index,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagChioriSkill,
			ICDGroup:   attacks.ICDGroupChioriSkill,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       thrustAtkScaling[c.TalentLvlSkill()],
		}

		snap := c.Snapshot(&ai)
		ai.FlatDmg = snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		ai.FlatDmg *= thrustDefScaling[c.TalentLvlSkill()]
		//TODO: hit box size
		c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), nil, 1.2), 0)

		return false
	}, a1SeizeTheMoment)
}

func (c *char) a1Tapestry() {
	c.a1Triggered = true
	//TODO: should this be per char status for hitlag?
	c.Core.Status.Add(a1SeizeTheMoment, 8*60)
	c.DeleteStatus(a1TailorMadeWindowKey)
	c.Core.Log.NewEvent("a1 seize the moment triggered", glog.LogCharacterEvent, c.Index)
	c.c4()
	c.c6()
}

func (c *char) tryTriggerA1Tailoring() {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.StatusIsActive(a1TailorMadeWindowKey) {
		return
	}
	c.a1Tailoring()
}

func (c *char) a1Tailoring() {
	c.a1Triggered = true
	// she can't use skill again if she triggers this
	c.DeleteStatus(a1TailorMadeWindowKey)
	c.c4()
	c.c6()
	c.Core.Log.NewEvent("a1 geo infusion triggered", glog.LogCharacterEvent, c.Index)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		a1GeoInfusionKey,
		attributes.Geo,
		300, // 5 s
		true,
		attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
	)
}

func (c *char) a1Window() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Triggered = false
	c.a1AttackCount = 0
	c.AddStatus(a1TailorMadeWindowKey, a1TailorMadeWindowLength, false) //TODO: hitlag on this?
	c.QueueCharTask(func() {
		if c.a1Triggered {
			return
		}
		c.a1Tailoring()
	}, a1TailorMadeWindowLength)
}

// When a nearby party member creates a Geo Construct, Chiori will gain 20% Geo DMG Bonus for 20s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4buff = make([]float64, attributes.EndStatType)
	c.a4buff[attributes.GeoP] = 0.20
	c.Core.Events.Subscribe(event.OnConstructSpawned, func(args ...interface{}) bool {
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a4BuffKey, 20*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return c.a4buff, true
			},
		})
		return false
	}, "chiori-a4")
}
