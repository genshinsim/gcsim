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

var a1DMGBuff = []float64{0.0, 0.15, 0.35, 0.65}
var a1ConversionChance = []float64{0.0, 0.33, 0.66, 1.0}

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
			m[attributes.DmgP] = a1DMGBuff[len(c.partyPHECTypes)]
			return m, true
		},
	})
}

func (c *char) a1Conversion() attributes.Element {
	if c.Base.Ascension < 1 {
		return attributes.Anemo
	}
	chance := a1ConversionChance[len(c.partyPHECTypes)]
	if c.Core.Rand.Float64() > chance {
		return attributes.Anemo
	}
	return c.partyPHECTypes[c.Core.Rand.Intn(len(c.partyPHECTypes))]
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Burning Shadowhunt Shell",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Anemo,
		Durability:     25,
		Mult:           1.5 * skillShadowhunt[c.TalentLvlSkill()],
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(_ ...interface{}) bool {
		bulletElem := attributes.Anemo
		if c.partyPHECTypes != nil {
			bulletElem = c.partyPHECTypes[c.Core.Rand.Intn(len(c.partyPHECTypes))]
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
