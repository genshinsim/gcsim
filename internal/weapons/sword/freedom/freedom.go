package freedom

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("freedom-sworn", weapon)
	core.RegisterWeaponFunc("freedomsworn", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.075 + float64(r)*0.025
	char.AddMod(coretype.CharStatMod{
		Key: "freedom-dmg",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	atkBuff := make([]float64, core.EndStatType)
	atkBuff[core.ATKP] = .15 + float64(r)*0.05

	buffNACAPlunge := make([]float64, core.EndStatType)
	buffNACAPlunge[core.DmgP] = .12 + 0.04*float64(r)

	icd := 0
	stacks := 0
	cooldown := 0

	stackFunc := func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if cooldown > c.Frame {
			return false
		}
		if icd > c.Frame {
			return false
		}

		icd = c.Frame + 30
		stacks++
		c.Log.NewEvent("freedomsworn gained sigil", coretype.LogWeaponEvent, char.Index(), "sigil", stacks)

		if stacks == 2 {
			stacks = 0
			c.AddStatus("freedom", 12*60)
			cooldown = c.Frame + 20*60
			for _, char := range c.Chars {
				// Attack buff snapshots so it needs to be in a separate mod
				char.AddMod(coretype.CharStatMod{
					Key: "freedom-proc",
					Amount: func() ([]float64, bool) {
						return atkBuff, true
					},
					Expiry: c.Frame + 12*60,
				})

				char.AddPreDamageMod(coretype.PreDamageMod{
					Key: "freedom-proc",
					Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
						switch atk.Info.AttackTag {
						case coretype.AttackTagNormal, coretype.AttackTagExtra, core.AttackTagPlunge:
							return buffNACAPlunge, true
						}
						return nil, false
					},
					Expiry: c.Frame + 12*60,
				})
			}
		}
		return false
	}

	for i := core.ReactionEventStartDelim + 1; i < core.ReactionEventEndDelim; i++ {
		c.Subscribe(i, stackFunc, "freedom-"+char.Name())
	}

	return "freedomsworn"
}
