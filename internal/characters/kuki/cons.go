package kuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

//When the Normal, Charged, or Plunging Attacks of the character affected by Shinobu's Grass Ring of Sanctification hit opponents,
// a Thundergrass Mark will land on the opponent's position and deal AoE Electro DMG based on 9.7% of Shinobu's Max HP.
//This effect can occur once every 5s.
func (c *char) c4() {
	//TODO: idk if the damage is instant or not
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		//ignore if C4 on icd
		if c.c4ICD > c.Core.F {
			return false
		}
		//On normal,charge and plunge attack
		if ae.Info.AttackTag != combat.AttackTagNormal && ae.Info.AttackTag != combat.AttackTagExtra && ae.Info.AttackTag != combat.AttackTagPlunge {
			return false
		}
		//make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}
		if c.Core.Status.Duration("kukibell") == 0 {
			return false
		}

		//TODO:frames for this and ICD tag
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "C4 proc",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0,
			FlatDmg:    c.MaxHP() * 0.097,
		}

		c.Core.QueueAttack(ai, combat.NewCircleHit(c.Core.Combat.Player(), 2, false, combat.TargettableEnemy), 5, 5)
		c.c4ICD = c.Core.F + 300 //5 sec icd
		return false
	}, "kuki-c4")
}

//When Kuki Shinobu takes lethal DMG, this instance of DMG will not take her down.
//This effect will automatically trigger when her HP reaches 1 and will trigger once every 60s.
//When Shinobu's HP drops below 25%, she will gain 150 Elemental Mastery for 15s. This effect will trigger once every 60s.
func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 150
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(_ ...interface{}) bool {
		if c.Core.F < c.c6ICD {
			return false
		}
		//check if hp less than 25%
		if c.HPCurrent/c.MaxHP() > .25 {
			return false
		}
		//if dead, revive back to 1 hp
		if c.HPCurrent <= -1 {
			c.HPCurrent = 1
		}

		//increase EM by 150 for 15s
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("kuki-c6", 900),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		c.c6ICD = c.Core.F + 3600

		return false
	}, "kuki-c6")
}
