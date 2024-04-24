package kuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// C4:
// When the Normal, Charged, or Plunging Attacks of the character affected by Shinobu's Grass Ring of Sanctification hit opponents,
//
//	a Thundergrass Mark will land on the opponent's position and deal AoE Electro DMG based on 9.7% of Shinobu's Max HP.
//
// This effect can occur once every 5s.
func (c *char) c4() {
	//TODO: idk if the damage is instant or not
	const c4IcdKey = "kuki-c4-icd"
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		trg := args[0].(combat.Target)
		// ignore if C4 on icd
		if c.StatusIsActive(c4IcdKey) {
			return false
		}
		// On normal,charge and plunge attack
		if ae.Info.AttackTag != attacks.AttackTagNormal && ae.Info.AttackTag != attacks.AttackTagExtra && ae.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		// make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if c.Core.Status.Duration(ringKey) == 0 {
			return false
		}
		c.AddStatus(c4IcdKey, 300, true) // 5s * 60

		//TODO:frames for this and ICD tag
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Thundergrass Mark",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0,
			FlatDmg:    c.MaxHP() * 0.097,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 2), 5, 5, c.particleCB)

		return false
	}, "kuki-c4")
}

// C6:
// When Kuki Shinobu takes lethal DMG, this instance of DMG will not take her down.
// This effect will automatically trigger when her HP reaches 1 and will trigger once every 60s.
// When Shinobu's HP drops below 25%, she will gain 150 Elemental Mastery for 15s. This effect will trigger once every 60s.
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 150
	const c6IcdKey = "kuki-c6-icd"
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(info.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		if c.StatusIsActive(c6IcdKey) {
			return false
		}
		// check if hp less than 25%
		if c.CurrentHPRatio() > 0.25 {
			return false
		}
		// if dead, revive back to 1 hp
		if c.CurrentHPRatio() <= 0 {
			c.SetHPByAmount(1)
		}
		c.AddStatus(c6IcdKey, 3600, false) // 60s * 60

		// increase EM by 150 for 15s
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("kuki-c6", 900),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		return false
	}, "kuki-c6")
}
