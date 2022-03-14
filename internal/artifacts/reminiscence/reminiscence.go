package reminiscence

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterSetFunc("reminiscence of shime", New)
	core.RegisterSetFunc("shimenawasreminiscence", New)
	core.RegisterSetFunc("shim", New)
}

func New(c coretype.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(coretype.CharStatMod{
			Key:    "shim-2pc",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	//11:51 AM] Episodde｜ShimenawaChildePeddler: Basically I found out that the fox set energy tax have around a 10 frame delay.
	//so I was testing if you can evade the fox set 15 energy tax by casting burst within those 10 frame after using an elemental
	//skill (not on hit). Turn out it work with childe :Childejoy:
	//The finding is now in #energy-drain-effects-have-a-delay if you want to take a closer look
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.50
		cd := -1
		s.Subscribe(core.PreSkill, func(args ...interface{}) bool {
			if s.Player.ActiveChar != c.Index()() {
				return false
			}
			if c.CurrentEnergy() < 15 {
				return false
			}
			if s.Frame < cd {
				return false
			}
			cd = s.Frame + 60*10

			//consume 15 energy, increased normal/charge/plunge dmg by 50%
			s.Tasks.Add(func() {
				c.AddEnergy("shim-4pc", -15)
			}, 10)
			c.AddPreDamageMod(coretype.PreDamageMod{
				Key:    "shim-4pc",
				Expiry: s.Frame + 60*10,
				Amount: func(atk *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
					return m, atk.Info.AttackTag == coretype.AttackTagNormal || atk.Info.AttackTag == coretype.AttackTagExtra || atk.Info.AttackTag == core.AttackTagPlunge
				},
			})

			return false
		}, fmt.Sprintf("shim-4pc-%v", c.Name()))

	}
}
