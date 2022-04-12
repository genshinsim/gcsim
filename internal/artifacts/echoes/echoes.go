package echoes

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("echoes of an offering", New)
	core.RegisterSetFunc("echoesofanoffering", New)
	core.RegisterSetFunc("echoes", New)
}

// 2pc - ATK +18%.
// 4pc - When Normal Attacks hit opponents, there is a 36% chance that it will trigger Valley Rite, which will increase Normal Attack DMG by 70% of ATK.
//  This effect will be dispelled 0.05s after a Normal Attack deals DMG.
//  If a Normal Attack fails to trigger Valley Rite, the odds of it triggering the next time will increase by 20%.
//  This trigger can occur once every 0.2s.
func New(c core.Character, s *core.Core, count int, params map[string]int) {
	prob := 0.36
	icd := 0
	procDuration := 3 // 0.05s
	procExpireF := 0

	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "echoes-2pc",
			Amount: func() ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}

	var dmgAdded float64

	if count >= 4 {
		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			// if the active char is not the equipped char then ignore
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			atk := args[1].(*core.AttackEvent)

			// If attack does not belong to the equipped character then ignore
			if atk.Info.ActorIndex != c.CharIndex() {
				return false
			}

			// If this is not a normal attack then ignore
			if atk.Info.AttackTag != core.AttackTagNormal {
				return false
			}

			// if buff is already active then buff attack
			if s.F < procExpireF {
				dmgAdded = (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]) * 0.7
				atk.Info.FlatDmg += dmgAdded
				s.Log.NewEvent("echoes 4pc adding dmg", core.LogArtifactEvent, c.CharIndex(),
					"buff_expiry", procExpireF,
					"dmg_added", dmgAdded,
				)
				return false
			}

			// If Artifact set effect is still on CD then ignore
			if s.F < icd {
				return false
			}

			if s.Rand.Float64() > prob {
				prob += 0.2
				return false
			}

			dmgAdded = (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]) * 0.7
			atk.Info.FlatDmg += dmgAdded

			procExpireF = s.F + procDuration
			icd = s.F + 12 //0.2s

			s.Log.NewEvent("echoes 4pc proc'd", core.LogArtifactEvent, c.CharIndex(),
				"probability", prob,
				"icd", icd,
				"buff_expiry", procExpireF,
				"dmg_added", atk.Info.FlatDmg,
			)

			prob = 0.36

			return false
		}, fmt.Sprintf("echoes-4pc-%v", c.Name()))
	}
}
