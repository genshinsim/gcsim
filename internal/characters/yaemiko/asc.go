package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// When casting Great Secret Art: Tenko Kenshin, each Sesshou Sakura destroyed
// resets the cooldown for 1 charge of Yakan Evocation: Sesshou Sakura.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// When there are at least 3 Sesshou Sakura nearby at the same time, unleashing Yakan Evocation:
// Sesshou Sakura will cause Yae Miko to unleash an additional lightning strike which deals Electro
// DMG at 40% of her ATK.
func (c *char) a1OnSkillPopKitsune() {
	if c.Base.Ascension < 1 {
		return
	}

	if !c.revelation {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Yae A1",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       0.4,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)

	if c.isRadianceSSC() {
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.Abil += stellarConductText
		ai.IgnoreDefPercent = 1
		ai.Mult = 0.5
	}

	c.Core.QueueAttack(ai, ap, 64, 64)
}

// Every point of Elemental Mastery Yae Miko possesses will increase Sesshou Sakura DMG by 0.15%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("yaemiko-a4", -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil
			}
			m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
			return m
		},
	})
}
