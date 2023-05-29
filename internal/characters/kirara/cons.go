package kirara

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4IcdStatus = "kirara-c4-icd"
	c6Status    = "kirara-c6"
)

// When Kirara is in the Urgent Neko Parcel state of Meow-teor Kick, she will grant other party members she crashes into Critical Transport Shields.
// The DMG absorption of Critical Transport Shield is 40% of the maximum absorption Meow-teor Kick's normal Shields of Safe Transport
// are capable of, and will absorb Dendro DMG with 250% effectiveness.
// Critical Transport Shields last 12s and can be triggered once on each character every 10s.
// co-op only
func (c *char) c2() {}

// After active character(s) protected by Shields of Safe Transport or Critical Transport Shields hit opponents with Normal, Charged, or Plunging Attacks,
// Kirara will perform a coordinated attack with them using Small Cat Grass Cardamoms, dealing 200% of her ATK as Dendro DMG. DMG dealt this way is
// considered Elemental Burst DMG. This effect can be triggered once every 3.8s. This CD is shared between all party members.
func (c *char) c4() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if c.StatusIsActive(c4IcdStatus) {
			return false
		}
		existingShield := c.Core.Player.Shields.Get(shield.ShieldKiraraSkill)
		if existingShield == nil {
			return false
		}

		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal,
			attacks.AttackTagExtra,
			attacks.AttackTagPlunge:
		default:
			return false
		}

		// TODO: snapshot? damage delay?
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Steed of Skanda",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       2,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 2), 0, 0)
		c.AddStatus(c4IcdStatus, 3.8*60, true)
		return false
	}, "kirara-c4")
}

// All nearby party members will gain 12% All Elemental DMG Bonus within 15s after Kirara uses her Elemental Skill or Burst.
func (c *char) c6() {
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(c6Status, 15*60),
			Amount: func() ([]float64, bool) {
				return c.c6Buff, true
			},
		})
	}
}
