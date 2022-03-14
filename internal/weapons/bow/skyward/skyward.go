package skyward

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterWeaponFunc("skyward harp", weapon)
	core.RegisterWeaponFunc("skywardharp", weapon)
}

func weapon(char coretype.Character, c *core.Core, r int, param map[string]int) string {
	//add passive crit, atk speed not sure how to do right now??
	//looks like jsut reduce the frames of normal attacks by 1 + 12%
	m := make([]float64, core.EndStatType)
	m[core.CD] = 0.15 + float64(r)*0.05
	cd := 270 - 30*r
	p := 0.5 + 0.1*float64(r)
	char.AddMod(coretype.CharStatMod{
		Key: "skyward harp",
		Amount: func() ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	icd := 0

	c.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		//check if char is correct?
		if atk.Info.ActorIndex != char.Index() {
			return false
		}
		if c.ActiveChar != char.Index() {
			return false
		}
		//check if cd is up
		if icd > c.Frame {
			return false
		}
		if c.Rand.Float64() > p {
			return false
		}

		//add a new action that deals % dmg immediately
		//superconduct attack
		ai := core.AttackInfo{
			ActorIndex: char.Index(),
			Abil:       "Skyward Harp Proc",
			AttackTag:  core.AttackTagWeaponSkill,
			ICDTag:     core.ICDTagNone,
			ICDGroup:   core.ICDGroupDefault,
			Element:    core.Physical,
			Durability: 100,
			Mult:       1.25,
		}
		c.Combat.QueueAttack(ai, core.NewDefCircHit(2, true, coretype.TargettableEnemy), 0, 1)

		//trigger cd
		icd = c.Frame + cd

		return false
	}, fmt.Sprintf("skyward-harp-%v", char.Name()))

	return "skywardharp"
}
