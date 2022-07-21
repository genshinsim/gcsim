package kokomi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const burstHitmark = 49

func init() {
	burstFrames = frames.InitAbilSlice(77)
}

// Burst - This function only handles initial damage and status setting
// Damage bonus modification is handled in a separate function based on status
// The might of Watatsumi descends, dealing Hydro DMG to surrounding opponents, before robing Kokomi in a Ceremonial Garment made from the flowing waters of Sangonomiya.
// Ceremonial Garment:
// - Sangonomiya Kokomi's Normal Attack, Charged Attack and Bake-Kurage DMG are increased based on her Max HP.
// - When her Normal and Charged Attacks hit opponents, Kokomi will restore HP for all nearby party members, and the amount restored is based on her Max HP.
// - Increases Sangonomiya Kokomi's resistance to interruption and allows her to move on the water's surface.
// These effects will be cleared once Sangonomiya Kokomi leaves the field.
func (c *char) Burst(p map[string]int) action.ActionInfo {
	// TODO: Snapshot timing is not yet known. Assume it's dynamic for now
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Nereid's Ascension",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       0,
	}
	ai.FlatDmg = burstDmg[c.TalentLvlBurst()] * c.MaxHP()

	c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 5, false, combat.TargettableEnemy), burstHitmark, burstHitmark)

	c.Core.Status.Add("kokomiburst", 10*60)

	// update jellyfish flat damage
	c.skillFlatDmg = c.burstDmgBonus(combat.AttackTagElementalArt)

	// Ascension 1 - reset duration of E Skill and also resnapshots it
	// Should not activate HoD consistent with in game since it is not a skill usage
	if c.Core.Status.Duration("kokomiskill") > 0 {
		// +1 to avoid same frame expiry issues with skill tick
		c.Core.Status.Add("kokomiskill", 12*60+1)
	}

	// C4 attack speed buff
	if c.Base.Cons >= 4 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.1
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("kokomi-c4", 10*60),
			AffectedStat: attributes.AtkSpd,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	// Cannot be prefed particles
	c.ConsumeEnergy(57)
	c.SetCDWithDelay(action.ActionBurst, 18*60, 46)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

// Helper function for determining whether burst damage bonus should apply
func (c *char) burstDmgBonus(a combat.AttackTag) float64 {
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return 0
	}
	switch a {
	case combat.AttackTagNormal:
		return burstBonusNormal[c.TalentLvlBurst()] * c.MaxHP()
	case combat.AttackTagExtra:
		return burstBonusCharge[c.TalentLvlBurst()] * c.MaxHP()
	case combat.AttackTagElementalArt:
		return burstBonusSkill[c.TalentLvlBurst()] * c.MaxHP()
	default:
		return 0
	}
}

// Implements event handler for healing during burst
// Also checks constellations
func (c *char) burstActiveHook() {
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}

		if c.Core.Status.Duration("kokomiburst") == 0 {
			return false
		}

		switch atk.Info.AttackTag {
		case combat.AttackTagNormal, combat.AttackTagExtra:
		default:
			return false
		}

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Ceremonial Garment",
			Src:     burstHealPct[c.TalentLvlBurst()]*c.MaxHP() + burstHealFlat[c.TalentLvlBurst()],
			Bonus:   c.Stat(attributes.Heal),
		})

		if c.Base.Cons >= 2 {
			c.c2()
		}
		if c.Base.Cons >= 4 {
			c.c4()
		}
		if c.Base.Cons >= 6 {
			c.c6()
		}

		return false
	}, "kokomi-q-healing")
}

// Clears Kokomi burst when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		// update jellyfish flat damage. regardless if burst is active or not
		if prev == c.Index {
			c.swapEarlyF = c.Core.F
			c.skillFlatDmg = c.burstDmgBonus(combat.AttackTagElementalArt)
		}
		c.Core.Status.Delete("kokomiburst")
		return false
	}, "kokomi-exit")
}
