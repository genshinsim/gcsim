package baizhu

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2ICDKey = "baizhu-c2-icd"

// Universal Diagnosis gains 1 additional charge.
func (c *char) c1() {
	c.SetNumCharges(action.ActionSkill, 2)
}

// When your own active character hits a nearby opponent with their attacks, Baizhu will unleash a Gossamer Sprite: Splice.
// Gossamer Sprite: Splice will initiate 1 attack before returning, dealing 250% of Baizhu's ATK as Dendro DMG and healing for 20% of Universal Diagnosis's Gossamer Sprite's normal healing.
// DMG dealt this way is considered Elemental Skill DMG.
// This effect can be triggered once every 5s.
func (c *char) c2() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		t := args[0].(combat.Target)
		// only trigger with the active character
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if c.StatusIsActive(c2ICDKey) {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Gossamer Sprite: Splice. (Baizhu's C2)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupBaizhuC2,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Dendro,
			Durability: 25,
			Mult:       2.5,
		}
		c.c6done = false
		var c6cb combat.AttackCBFunc
		if c.Base.Cons >= 6 {
			c6cb = c.makeC6CB()
		}
		// TODO: accurate C2 hitmark and return travel values
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(t, nil, 0.6),
			0,
			skillFirstHitmark, // reuse skill for now
			c6cb,
		)

		// C2 healing
		c.Core.Tasks.Add(func() {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "Baizhu's C2: Healing",
				Src:     (skillHealPP[c.TalentLvlSkill()]*c.MaxHP() + skillHealFlat[c.TalentLvlSkill()]) * 0.2,
				Bonus:   c.Stat(attributes.Heal),
			})
		}, skillReturnTravel) // reuse skill for now

		c.AddStatus(c2ICDKey, 60*5, false) // 5s
		return false
	}, "baizhu-c2")
}

// For 15s after Healing Holism is used, Baizhu will increase all nearby party members' Elemental Mastery by 80.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 80
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("baizhu-c4", 900),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
}

// Increases the DMG dealt by Holistic Revivification's Spiritveins by 8% of Baizhu's Max HP.
// Additionally, when Gossamer Sprite or Gossamer Sprite: Splice hit opponents, there is a 100% chance of generating one of Healing Holism's
// Seamless Shields. This effect can only be triggered once by a Gossamer Sprite or Gossamer Sprite: Splice.
func (c *char) makeC6CB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if done {
			return
		}
		done = true
		c.summonSeamlessShield()
	}
}
