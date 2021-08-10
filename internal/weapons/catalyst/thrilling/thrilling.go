package thrilling

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("thrilling tales of dragon slayers", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {
	last := 0
	isActive := false

	s.AddInitHook(func() {
		isActive = s.ActiveCharIndex() == c.CharIndex()
	})

	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .16 + float64(r)*0.06

	s.AddEventHook(func(s core.Sim) bool {
		if !isActive && s.ActiveCharIndex() == c.CharIndex() {
			//swapped to current char
			isActive = true
			return false
		}

		//swap from current char to new char
		if isActive && s.ActiveCharIndex() != c.CharIndex() {
			isActive = false

			//do nothing if off cd
			if last != 0 && s.Frame()-last < 1200 {
				return false
			}
			//trigger buff if not on cd

			last = s.Frame()
			expiry := s.Frame() + 600

			active, _ := s.CharByPos(s.ActiveCharIndex())
			active.AddMod(core.CharStatMod{
				Key: "thrilling tales",
				Amount: func(a core.AttackTag) ([]float64, bool) {
					return m, expiry > s.Frame()
				},
				Expiry: -1,
			})

			log.Debugw("ttds activated", "frame", s.Frame(), "event", core.LogWeaponEvent, "char", active.CharIndex(), "expiry", expiry)
		}

		return false
	}, fmt.Sprintf("thrilling-%v", c.Name()), core.PostSwapHook)
}
