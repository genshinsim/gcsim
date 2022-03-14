package polarstar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("polar star", weapon)
	core.RegisterWeaponFunc("polarstar", weapon)
}

// Elemental Skill and Elemental Burst DMG increased by 12/15/18/21/24%.
// After a Normal Attack, Charged Attack, Elemental Skill or Elemental Burst hits an opponent, 1 stack of Ashen Nightstar will be gained for 12s.
// When 1/2/3/4 stacks of Ashen Nightstar are present, ATK is increased by (10/20/30/48)/(12.5/25/37.5/60)/(15/30/45/72)/(17.5/35/52.5/84)/(20/40/60/96)%.
// The stack of Ashen Nightstar created by the Normal Attack, Charged Attack, Elemental Skill or Elemental Burst will be counted independently of the others.
func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {

	dmg := .09 + float64(r)*.03
	stack := .075 + float64(r)*.025
	max := .06 + float64(r)*.02

	normal := 0
	charged := 0
	skill := 0
	burst := 0

	mATK := make([]float64, core.EndStatType)
	char.AddMod(coretype.CharStatMod{
		Key:    "polar-star",
		Expiry: -1,
		Amount: func() ([]float64, bool) {
			count := 0
			if normal > c.Frame {
				count++
			}
			if charged > c.Frame {
				count++
			}
			if skill > c.Frame {
				count++
			}
			if burst > c.Frame {
				count++
			}

			atkbonus := stack * float64(count)
			if count >= 4 {
				atkbonus += max
			}
			mATK[core.ATKP] = atkbonus

			return mATK, true
		},
	})

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.ActiveChar != char.Index() {
			return false
		}

		cd := c.Frame + 60*12
		switch atk.Info.AttackTag {
		case coretype.AttackTagNormal:
			normal = cd
		case coretype.AttackTagExtra:
			charged = cd
		case core.AttackTagElementalArt, core.AttackTagElementalArtHold:
			skill = cd
		case core.AttackTagElementalBurst:
			burst = cd
		}

		return false
	}, fmt.Sprintf("polar-star-%v", char.Name()))

	mDmg := make([]float64, core.EndStatType)
	mDmg[core.DmgP] = dmg
	char.AddPreDamageMod(coretype.PreDamageMod{
		Key:    "polar-star",
		Expiry: -1,
		Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case core.AttackTagElementalArt, core.AttackTagElementalArtHold, core.AttackTagElementalBurst:
				return mDmg, true
			}
			return nil, false
		},
	})

	return "polarstar"
}
