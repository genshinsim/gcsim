package reminiscence

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("reminiscence of shime", New)
	core.RegisterSetFunc("shimenawasreminiscence", New)
}

func New(c core.Character, s *core.Core, count int, params map[string]int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "rem-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	//11:51 AM] Episodde｜ShimenawaChildePeddler: Basically I found out that the fox set energy tax have around a 10 frame delay.
	//so I was testing if you can evade the fox set 15 energy tax by casting burst within those 10 frame after using an elemental
	//skill (not on hit). Turn out it work with childe :Childejoy:
	//The finding is now in #energy-drain-effects-have-a-delay if you want to take a closer look
	if count >= 4 {
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.50
		s.Events.Subscribe(core.PreSkill, func(args ...interface{}) bool {
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			if c.CurrentEnergy() > 15 {
				//consume 15 energy, increased normal/charge/plunge dmg by 50%
				s.Tasks.Add(func() {
					c.AddEnergy(-15)
				}, 10)
				c.AddMod(core.CharStatMod{
					Key: "rem-4pc",
					Amount: func(ds core.AttackTag) ([]float64, bool) {
						if ds != core.AttackTagNormal && ds != core.AttackTagExtra && ds != core.AttackTagPlunge {
							return nil, false
						}
						return m, true
					},
					Expiry: s.F + 600,
				})
			}
			return false
		}, fmt.Sprintf("rem-4pc-%v", c.Name()))

	}
	//add flat stat to char
}
