package kuki

import "github.com/genshinsim/gcsim/pkg/core"

//Gyoei Narukami Kariyama Rite's AoE is increased by 50%.
func (c *char) c1() {
	c.c1AoeMod = 3.5

}

//Grass Ring of Sanctification's duration is increased by 3s.
func (c *char) c2() {
	c.skilldur = 900 //12+3s

}

//When the Normal, Charged, or Plunging Attacks of the character affected by Shinobu's Grass Ring of Sanctification hit opponents,
// a Thundergrass Mark will land on the opponent's position and deal AoE Electro DMG based on 9.7% of Shinobu's Max HP.
//This effect can occur once every 5s.
func (c *char) c4() {
	//TODO: idk if the damage is instant or not
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)
		//ignore if C4 on icd
		if c.c4ICD > c.Core.F {
			return false
		}
		//On normal,charge and plunge attack
		if ae.Info.AttackTag != core.AttackTagNormal && ae.Info.AttackTag != core.AttackTagExtra && ae.Info.AttackTag != core.AttackTagPlunge {
			return false
		}
		//make sure the person triggering the attack is on field still
		if ae.Info.ActorIndex != c.Core.ActiveChar {
			return false
		}
		if c.Core.Status.Duration("kukibell") == 0 {
			return false
		}

		//TODO:frames for this and ICD tag
		ai := core.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "C4 proc",
			AttackTag:  core.AttackTagElementalArt,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Electro,
			Durability: 25,
			Mult:       0,
			FlatDmg:    c.MaxHP() * 0.097,
		}

		//Particle check is 45% for particle
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), 5, 5)
		if c.Core.Rand.Float64() < .45 {
			c.QueueParticle("Kuki", 1, core.Electro, 100) // TODO: idk the particle timing yet fml (or probability)
		}
		c.c4ICD = c.Core.F + 300 //5 sec icd
		return false
	}, "kuki-c4")

}

//When Kuki Shinobu takes lethal DMG, this instance of DMG will not take her down.
//This effect will automatically trigger when her HP reaches 1 and will trigger once every 60s.
//When Shinobu's HP drops below 25%, she will gain 150 Elemental Mastery for 15s. This effect will trigger once every 60s.

func (c *char) c6() {
	c.Core.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.F < c.c6icd && c.c6icd != 0 {
			return false
		}
		//check if hp less than 25%
		if c.HP()/c.MaxHP() > .25 {
			return false
		}
		//if dead, revive back to 1 hp
		if c.HP() <= -1 {
			c.HPCurrent = 1
		}

		//increase EM by 150 for 15s
		val := make([]float64, core.EndStatType)
		val[core.EM] = 150
		c.AddMod(core.CharStatMod{
			Key:    "kuki-c6",
			Amount: func() ([]float64, bool) { return val, true },
			Expiry: c.Core.F + 900,
		})

		c.c6icd = c.Core.F + 3600

		return false
	}, "kuki-c6")

}
