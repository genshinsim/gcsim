package kokomi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 49
	burstKey     = "kokomiburst"
)

func init() {
	burstFrames = frames.InitAbilSlice(78) // Q -> D/J
	burstFrames[action.ActionAttack] = 77  // Q -> N1
	burstFrames[action.ActionCharge] = 77  // Q -> CA
	burstFrames[action.ActionSkill] = 77   // Q -> E
	burstFrames[action.ActionWalk] = 77    // Q -> W
	burstFrames[action.ActionSwap] = 76    // Q -> Swap
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
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 50,
		Mult:       0,
	}
	ai.FlatDmg = burstDmg[c.TalentLvlBurst()] * c.MaxHP()

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), burstHitmark, burstHitmark)

	c.Core.Status.Add(burstKey, 10*60)

	// update jellyfish flat damage
	c.skillFlatDmg = c.burstDmgBonus(attacks.AttackTagElementalArt)

	if c.Base.Ascension >= 1 {
		c.Core.Tasks.Add(c.a1, 46)
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
		CanQueueAfter:   burstFrames[action.ActionSwap],
		State:           action.BurstState,
	}
}

// Helper function for determining whether burst damage bonus should apply
func (c *char) burstDmgBonus(a attacks.AttackTag) float64 {
	if c.Core.Status.Duration("kokomiburst") == 0 {
		return 0
	}
	switch a {
	case attacks.AttackTagNormal:
		return burstBonusNormal[c.TalentLvlBurst()] * c.MaxHP()
	case attacks.AttackTagExtra:
		return burstBonusCharge[c.TalentLvlBurst()] * c.MaxHP()
	case attacks.AttackTagElementalArt:
		return burstBonusSkill[c.TalentLvlBurst()] * c.MaxHP()
	default:
		return 0
	}
}

// - implements burst healing, C2 and C6 handling
//
// When her Normal and Charged Attacks hit opponents,
// Kokomi will restore HP for all nearby party members,
// and the amount restored is based on her Max HP.
func (c *char) makeBurstHealCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != combat.TargettableEnemy {
			return
		}
		if c.Core.Status.Duration("kokomiburst") == 0 {
			return
		}
		if done {
			return
		}
		done = true

		heal := burstHealPct[c.TalentLvlBurst()]*c.MaxHP() + burstHealFlat[c.TalentLvlBurst()]
		for _, char := range c.Core.Player.Chars() {
			src := heal

			// C2 handling
			// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
			// Nereid's Ascension Normal and Charged Attacks: 0.6% of Kokomi's Max HP.
			if c.Base.Cons >= 2 && char.HPCurrent/char.MaxHP() <= .5 {
				bonus := 0.006 * c.MaxHP()
				src += bonus
				c.Core.Log.NewEvent("kokomi c2 proc'd", glog.LogCharacterEvent, char.Index).
					Write("bonus", bonus)
			}
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  char.Index,
				Message: "Ceremonial Garment",
				Src:     src,
				Bonus:   c.Stat(attributes.Heal),
			})
		}

		if c.Base.Cons >= 6 {
			c.c6()
		}
	}
}

// Clears Kokomi burst when she leaves the field
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		prev := args[0].(int)
		// update jellyfish flat damage. regardless if burst is active or not
		if prev == c.Index {
			c.swapEarlyF = c.Core.F
			c.skillFlatDmg = c.burstDmgBonus(attacks.AttackTagElementalArt)
		}
		c.Core.Status.Delete("kokomiburst")
		return false
	}, "kokomi-exit")
}
