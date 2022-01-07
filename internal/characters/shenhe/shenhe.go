package shenhe

import (
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Shenhe, NewChar)
}

type char struct {
	*character.Tmpl
	quillcount []int
	c4count    int
	c4expiry   int

	// eChargeMax       int
	// eCharges         int
	// skillRecoverySrc int
	// recoverQueue     []int //queue of recovery

	cdQueueWorkerStartedAt []int
	cdQueue                [][]int
	availableCDCharge      []int
	additionalCDCharge     []int
}

const (
	quillKey = "shenhequill"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 80
	c.EnergyMax = 80
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = core.ZoneLiyue
	c.Base.Element = core.Cryo

	c.c4count = 0
	c.c4expiry = 0

	c.cdQueueWorkerStartedAt = make([]int, core.EndActionType)
	c.cdQueue = make([][]int, core.EndActionType)
	c.additionalCDCharge = make([]int, core.EndActionType)
	c.availableCDCharge = make([]int, core.EndActionType)

	for i := 0; i < len(c.cdQueue); i++ {
		c.cdQueue[i] = make([]int, 0, 4)
		c.availableCDCharge[i] = 1
	}

	if c.Base.Cons >= 1 {
		c.additionalCDCharge[core.ActionSkill] = 1
		c.availableCDCharge[core.ActionSkill]++
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}

	c.quillDamageMod()

	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	// if c.Base.Cons >= 6 {
	// 	c.c6()
	// }
	c.a2()
	c.quillcount = make([]int, len(c.Core.Chars))
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 25
	default:
		c.Core.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Key.String(), a.String())
		return 0
	}

}

// inspired from barbara c2
// TODO: technically always assumes you are inside shenhe burst
func (c *char) a2() {
	val := make([]float64, core.EndStatType)
	val[core.CryoP] = 0.15
	for _, char := range c.Core.Chars {
		// if i == c.Index {
		// 	continue
		// }
		char.AddMod(core.CharStatMod{
			Key:    "shenhe-a2",
			Expiry: -1,
			Amount: func() ([]float64, bool) {
				if c.Core.Status.Duration("shenheburst") > 0 {
					return val, true
				} else {
					return nil, false
				}
			},
		})
	}
}

func (c *char) c6() {
	val := make([]float64, core.EndStatType)
	val[core.CD] = 0.15
	for _, char := range c.Core.Chars {
		char.AddPreDamageMod(core.PreDamageMod{
			Key:    "shenhe-c2",
			Expiry: -1,
			Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
				//check if tags active
				if c.Core.Status.Duration("shenheburst") <= 0 {
					return nil, false
				}
				if atk.Info.Element != core.Cryo {
					return nil, false
				}
				return val, true
			},
		})
	}
}

func (c *char) c4() {
	c.AddPreDamageMod(core.PreDamageMod{
		Expiry: -1,
		Key:    "shenhe-c4",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			val := make([]float64, core.EndStatType)

			if atk.Info.AttackTag != core.AttackTagElementalArt && atk.Info.AttackTag != core.AttackTagElementalArtHold {
				return nil, false
			}
			if c.Core.F >= c.c4expiry {
				return nil, false
			}
			val[core.DmgP] += 0.05 * float64(c.c4count)
			c.c4count = 0
			c.c4expiry = 0
			return val, true
		},
	})
}
