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

	if count >= 2 {
		mATK := make([]float64, core.EndStatType)
		mATK[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "echoes-2pc",
			Amount: func() ([]float64, bool) {
				return mATK, true
			},
			Expiry: -1,
		})
	}

	if count >= 4 {
		mDMG := make([]float64, core.EndStatType)
		var triggerAtk *core.AttackEvent

		s.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
			// if the active char is not the equipped char then ignore
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			atk := args[1].(*core.AttackEvent)

			// If this is not a normal attack then ignore
			if atk.Info.AttackTag != core.AttackTagNormal {
				return false
			}

			// If Artifact set effect is still on CD then ignore
			if s.F < icd {
				return false
			}

			icd = s.F + 0.2*60

			if s.Rand.Float64() < prob {
				triggerAtk = atk
				dmgAdded := (atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[core.ATKP]) + atk.Snapshot.Stats[core.ATK]) * 0.7
				mDMG[core.Dmg] = dmgAdded

				c.AddMod(core.CharStatMod{
					Key: "echoes-4pc",
					Amount: func() ([]float64, bool) {
						return mDMG, true
					},
					Expiry: -1,
				})

				s.Log.NewEvent("echoes 4pc proc", core.LogArtifactEvent, c.CharIndex(),
					"probability", prob,
					"dmg_added", dmgAdded,
				)

				prob = 0.36
			} else {
				prob += 0.2
			}

			return false
		}, fmt.Sprintf("echoes-4pc-%v", c.Name()))

		s.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
			// if the active char is not the equipped char then ignore
			if s.ActiveChar != c.CharIndex() {
				return false
			}

			atkOnDmg := args[1].(*core.AttackEvent)

			// If this is not a normal attack then ignore
			if atkOnDmg.Info.AttackTag != core.AttackTagNormal {
				return false
			}

			// If this is the same attack that triggered echoes-4pc then ignore
			if atkOnDmg == triggerAtk {
				return false
			}

			// Only set a CD on the mod if it's currently active
			if c.ModIsActive("echoes-4pc") {
				c.AddMod(core.CharStatMod{
					Key: "echoes-4pc",
					Amount: func() ([]float64, bool) {
						return mDMG, true
					},
					Expiry: 0.05 * 60, // 3 frames
				})
			}

			return false
		}, fmt.Sprintf("echoes-4pc-ondamage-%v", c.Name()))

	}
}
