package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

//After defeating an enemy, ATK is increased by 12/15/18/21/24% for 30s.
//This effect has a maximum of 3 stacks, and the duration of each stack is independent of the others.
func Blackcliff(char core.Character, c *core.Core, r int, param map[string]int) {

	atk := 0.09 + float64(r)*0.03
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key: "blackcliff",
		Amount: func() ([]float64, bool) {
			count := 0
			for _, v := range stacks {
				if v > c.F {
					count++
				}
			}
			m[core.ATKP] = atk * float64(count)
			return m, true
		},
		Expiry: -1,
	})

	c.Events.Subscribe(core.OnTargetDied, func(args ...interface{}) bool {
		stacks[index] = c.F + 1800
		index++
		if index == 3 {
			index = 0
		}
		return false
	}, fmt.Sprintf("blackcliff-%v", char.Name()))

}

func Favonius(char core.Character, c *core.Core, r int, param map[string]int) {

	p := 0.50 + float64(r)*0.1
	cd := 810 - r*90
	icd := 0
	//add on crit effect
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}

		if c.Rand.Float64() > p {
			return false
		}
		c.Log.NewEvent("favonius proc'd", core.LogWeaponEvent, char.CharIndex())

		char.QueueParticle("favonius-"+char.Name(), 3, core.NoElement, 80)

		icd = c.F + cd

		return false
	}, fmt.Sprintf("favo-%v", char.Name()))

}

func Lithic(char core.Character, c *core.Core, r int, param map[string]int) {

	stacks := 0
	val := make([]float64, core.EndStatType)

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		for _, char := range c.Chars {
			if char.Zone() == core.ZoneLiyue {
				stacks++
			}
		}
		val[core.CR] = (0.02 + float64(r)*0.01) * float64(stacks)
		val[core.ATKP] = (0.06 + float64(r)*0.01) * float64(stacks)
		return true
	}, fmt.Sprintf("lithic-%v", char.Name()))

	char.AddMod(core.CharStatMod{
		Key:    "lithic",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
}

func Royal(char core.Character, c *core.Core, r int, param map[string]int) {
	stacks := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		crit := args[3].(bool)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if crit {
			stacks = 0
		} else {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
		}
		return false
	}, fmt.Sprintf("royal-%v", char.Name()))

	rate := 0.06 + float64(r)*0.02
	m := make([]float64, core.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:    "royal",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			m[core.CR] = float64(stacks) * rate
			return m, true
		},
	})

}

//After damaging an opponent with an Elemental Skill, the skill has a 40/50/60/70/80%
//chance to end its own CD. Can only occur once every 30/26/22/19/16s.
func Sacrificial(char core.Character, c *core.Core, r int, param map[string]int) {

	last := 0
	prob := 0.3 + float64(r)*0.1
	cd := (34 - r*4) * 60

	if r >= 4 {
		cd = (19 - (r-4)*3) * 60
	}
	//add on crit effect
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
			return false
		}
		if last != 0 && c.F-last < cd {
			return false
		}
		if char.Cooldown(core.ActionSkill) == 0 {
			return false
		}
		if c.Rand.Float64() < prob {
			char.ResetActionCooldown(core.ActionSkill)
			last = c.F
			c.Log.NewEvent("sacrificial proc'd", core.LogWeaponEvent, char.CharIndex())
		}
		return false
	}, fmt.Sprintf("sac-%v", char.Name()))

}

//For every point of the entire party's combined maximum Energy capacity,
//the Elemental Burst DMG of the character equipping this weapon is increased by 0.12%.
//A maximum of 40% increased Elemental Burst DMG can be achieved this way.
//r1 0.12 40%
//r2 0.15 50%
//r3 0.18 60%
//r4 0.21 70%
//r5 0.24 80%
func Wavebreaker(char core.Character, c *core.Core, r int, param map[string]int) {

	per := 0.09 + 0.03*float64(r)
	max := 0.3 + 0.1*float64(r)

	var amt float64

	c.Events.Subscribe(core.OnInitialize, func(args ...interface{}) bool {
		var energy float64
		//calculate total team energy
		for _, x := range c.Chars {
			energy += x.MaxEnergy()
		}

		amt = energy * per / 100
		if amt > max {
			amt = max
		}
		c.Log.NewEvent("wavebreaker dmg calc", core.LogWeaponEvent, char.CharIndex(), "total", energy, "per", per, "max", max, "amt", amt)
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = amt
		char.AddPreDamageMod(core.PreDamageMod{
			Expiry: -1,
			Key:    "wavebreaker",
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				if atk.Info.AttackTag == core.AttackTagElementalBurst {
					return m, true
				}
				return nil, false
			},
		})
		return true
	}, fmt.Sprintf("wavebreaker-%v", char.Name()))

}

// Golden Majesty:
// Increases Shield Strength by 20/25/30/35/40%.
// Scoring hits on opponents increases ATK by 4/5/6/7/8% for 8s.
// Max 5 stacks. Can only occur once every 0.3s.
// While protected by a shield, this ATK increase effect is increased by 100%.
func GoldenMajesty(char core.Character, c *core.Core, r int, param map[string]int) {

	shd := .15 + float64(r)*.05
	atkbuff := 0.03 + 0.01*float64(r)

	c.Shields.AddBonus(func() float64 { return shd })

	icd := -1
	stacks := 0
	expiry := 0
	m := make([]float64, core.EndStatType)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)

		if ae.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		if c.F < icd {
			return false
		}
		icd = c.F + 18

		// reset stacks if expired
		if c.F > expiry {
			stacks = 0
		}

		stacks++
		if stacks > 5 {
			stacks = 5
		}

		expiry = c.F + 60*8
		char.AddMod(core.CharStatMod{
			Key:    "golden-majesty",
			Expiry: expiry,
			Amount: func() ([]float64, bool) {
				m[core.ATKP] = atkbuff * float64(stacks)
				if c.Shields.IsShielded(char.CharIndex()) {
					m[core.ATKP] *= 2
				}
				return m, true
			},
		})

		return false
	}, fmt.Sprintf("golden-majesty-%v", char.Name()))
}

func NoEffectWeapon(key string) core.NewWeaponFunc {
	return func(char core.Character, c *core.Core, r int, param map[string]int) string {
		return key
	}
}
