package chasca

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var a1DMGBuff = []float64{0.0, 0.15, 0.35, 0.65, 0.65} // has an extra 0.65 for c2 stack
var a1ConversionChance = []float64{0.0, 0.333, 0.667, 1.0}

func (c *char) a1DMGBuff() {
	if c.Base.Ascension < 1 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	// assuming we don't need to add and remove this buff constantly
	// since it would be active for all E-CAs anyways
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("chasca-a1", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.ICDTag != attacks.ICDTagChascaShining {
				return nil, false
			}
			m[attributes.DmgP] = a1DMGBuff[len(c.partyPHECTypesUnique)+c.c2A1Stack()]
			return m, true
		},
	})
}

func (c *char) a1Conversion() attributes.Element {
	if c.Base.Ascension < 1 {
		return attributes.Anemo
	}
	if len(c.partyPHECTypesUnique) == 0 {
		return attributes.Anemo
	}
	chance := a1ConversionChance[len(c.partyPHECTypesUnique)]
	chance += c.c1()
	if c.Core.Rand.Float64() > chance {
		return attributes.Anemo
	}
	c.c6()
	return c.partyPHECTypesUnique[c.Core.Rand.Intn(len(c.partyPHECTypesUnique))]
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Burning Shadowhunt Shot",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           1.5 * skillShadowhunt[c.TalentLvlSkill()],
		IsDeployable:   true,
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(_ ...interface{}) bool {
		bulletElem := attributes.Anemo
		if len(c.partyPHECTypesUnique) > 0 {
			bulletElem = c.partyPHECTypesUnique[c.Core.Rand.Intn(len(c.partyPHECTypesUnique))]
		}
		switch bulletElem {
		case attributes.Anemo:
			ai.Abil = "Burning Shadowhunt Shell"
			ai.Element = attributes.Anemo
			ai.Mult = 1.5 * skillShadowhunt[c.TalentLvlSkill()]
		default:
			ai.Abil = fmt.Sprintf("Burning Shadowhunt Shell (%s)", bulletElem)
			ai.Element = bulletElem
			ai.Mult = 1.5 * skillShining[c.TalentLvlSkill()]
		}
		ap := combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key())
		c.Core.QueueAttack(ai, ap, 0, 60)
		return false
	}, "chasca-a4")
}
