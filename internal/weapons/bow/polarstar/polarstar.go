package polarstar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("polar star", weapon)
	core.RegisterWeaponFunc("polarstar", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {

	dmg := .09 + float64(r)*.03
	stack := .075 + float64(r)*.025
	max := .06 + float64(r)*.02

	normal := 0
	charged := 0
	skill := 0
	burst := 0

	mATK := make([]float64, attributes.EndStatType)
	char.AddMod(core.CharStatMod{
		Key:	"polar-star",
		Expiry:	-1,
		Amount: func() ([]float64, bool) {
			count := 0
			if normal > c.F {
				count++
			}
			if charged > c.F {
				count++
			}
			if skill > c.F {
				count++
			}
			if burst > c.F {
				count++
			}

			atkbonus := stack * float64(count)
			if count >= 4 {
				atkbonus += max
			}
			mATK[attributes.ATKP] = atkbonus

			return mATK, true
		},
	})

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}

		cd := c.F + 60*12
		switch atk.Info.AttackTag {
		case combat.AttackTagNormal:
			normal = cd
		case combat.AttackTagExtra:
			charged = cd
		case combat.AttackTagElementalArt, combat.AttackTagElementalArtHold:
			skill = cd
		case combat.AttackTagElementalBurst:
			burst = cd
		}

		return false
	}, fmt.Sprintf("polar-star-%v", char.Name()))

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = dmg
	char.AddAttackMod("polar-star",
		-1,
		func(atk *combat.AttackEvent, t combat.Target,) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case combat.AttackTagElementalArt, combat.AttackTagElementalArtHold, combat.AttackTagElementalBurst:
				return mDmg, true
			}
			return nil, false
		})

	return "polarstar"
}
