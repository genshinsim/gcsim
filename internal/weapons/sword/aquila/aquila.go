package aquila

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("aquila favonia", weapon)
	core.RegisterWeaponFunc("aquilafavonia", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) string {
	m := make([]float64, core.EndStatType)
	m[core.ATKP] = .15 + .05*float64(r)
	char.AddMod(core.CharStatMod{
		Key: "acquila favonia",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	dmg := 1.7 + .3*float64(r)
	heal := .85 + .15*float64(r)

	last := -1

	c.Events.Subscribe(core.OnCharacterHurt, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if c.F-last < 900 && last != -1 {
			return false
		}
		last = c.F
		ai := core.AttackInfo{
			ActorIndex: char.CharIndex(),
			Abil:       "Aquila Favonia",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       dmg,
		}
		snap := char.Snapshot(&ai)
		c.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), 1)

		atk := snap.BaseAtk*(1+snap.Stats[core.ATKP]) + snap.Stats[core.ATK]

		// c.Log.NewEvent("acquila heal triggered", core.LogWeaponEvent, char.CharIndex(), "atk", atk, "heal amount", atk*heal)
		c.Health.Heal(core.HealInfo{
			Caller:  char.CharIndex(),
			Target:  c.ActiveChar,
			Message: "Aquila Favonia",
			Src:     atk * heal,
			Bonus:   char.Stat(core.Heal),
		})
		return false
	}, fmt.Sprintf("aquila-%v", char.Name()))
	return "aquilafavonia"
}
