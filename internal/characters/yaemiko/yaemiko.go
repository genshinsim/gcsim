package yaemiko

import (
	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
)

func init() {
	core.RegisterCharFunc(core.YaeMiko, NewChar)
}

type char struct {
	*character.Tmpl
	kitsunes               []*kitsune
	cdQueueWorkerStartedAt []int
	cdQueue                [][]int
	availableCDCharge      []int
	additionalCDCharge     []int
	a2skillTimer           int
	a2burstTimer           int
	totemParticleICD       int
}

const (
	yaeTotemStatus = "oldestTotemTime"
	yaeTotemCount  = "totems"
)

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Electro
	c.Energy = 90
	c.EnergyMax = 90
	c.Weapon.Class = core.WeaponClassCatalyst
	c.NormalHitNum = 3

	c.BurstCon = 5
	c.SkillCon = 3

	c.cdQueueWorkerStartedAt = make([]int, core.EndActionType)
	c.cdQueue = make([][]int, core.EndActionType)
	c.additionalCDCharge = make([]int, core.EndActionType)
	c.availableCDCharge = make([]int, core.EndActionType)
	c.kitsunes = make([]*kitsune, 0, 5)
	c.totemParticleICD = 0

	for i := 0; i < len(c.cdQueue); i++ {
		c.cdQueue[i] = make([]int, 0, 4)
		c.availableCDCharge[i] = 1
	}

	c.additionalCDCharge[core.ActionSkill] = 2
	c.availableCDCharge[core.ActionSkill] += 2
	c.a2burstTimer = 0
	c.a2skillTimer = 0
	c.Tags["eCharge"] = c.availableCDCharge[core.ActionSkill]
	return &c, nil
}

func (c *char) Init(index int) {
	c.Tmpl.Init(index)
	// c.a2()
	c.a4()
	if c.Base.Cons >= 4 {
		c.c4()
	}
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	case core.ActionCharge:
		return 50
	default:
		return 0
	}
}

// func (c *char) a2() {
// 	//Other nearby party members can decrease the CD of Yae Miko's Yakan Evocation: Sesshou Sakura:
// 	// • Hitting opponents with Elemental Skill DMG decreases it by 1s and can occur once every 1.8s.
// 	// • Hitting opponents with Elemental Burst DMG decreases it by 1s and can occur once every 1.8s.

// 	c.Core.Events.Subscribe(core.OnAttackWillLand, func(args ...interface{}) bool {
// 		atk := args[1].(*core.AttackEvent)
// 		if c.Index == atk.Info.ActorIndex {
// 			// do not trigger for yae attacks
// 			return false
// 		}
// 		switch atk.Info.AttackTag {
// 		case core.AttackTagElementalBurst:
// 			if c.Core.F < c.a2burstTimer+1.8*60 {
// 				return false
// 			} else {
// 				c.ReduceActionCooldown(core.ActionSkill, 60)
// 				c.a2burstTimer = c.Core.F
// 			}
// 		case core.AttackTagElementalArt:
// 			if c.Core.F < c.a2skillTimer+1.8*60 {
// 				return false
// 			} else {
// 				c.ReduceActionCooldown(core.ActionSkill, 60)
// 				c.a2skillTimer = c.Core.F
// 			}
// 		case core.AttackTagElementalArtHold:
// 			if c.Core.F < c.a2skillTimer+1.8*60 {
// 				return false
// 			} else {
// 				c.ReduceActionCooldown(core.ActionSkill, 60)
// 				c.a2skillTimer = c.Core.F
// 			}
// 		default:
// 			return false
// 		}
// 		return false
// 	}, "yaemiko-a2")

// }

func (c *char) a4() {
	c.AddPreDamageMod(core.PreDamageMod{
		Expiry: -1,
		Key:    "yaemiko-a2",
		Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
			if atk.Info.AttackTag != core.AttackTagElementalArt {
				// only trigger on elemental art damage
				return nil, false
			}
			val := make([]float64, core.EndStatType)
			val[core.DmgP] = c.Stats[core.EM] * 0.0015
			return val, true
		},
	})
}

func (c *char) c4() {
	// c4
	// When Sesshou Sakura thunderbolts hit opponents, the Electro DMG Bonus of all nearby party members is increased by 20% for 5s.
	c.Core.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		// TODO: does this trigger for yaemiko too? assuming it does
		for _, char := range c.Core.Chars {
			char.AddPreDamageMod(core.PreDamageMod{
				Expiry: c.Core.F + 5*60,
				Key:    "yaemiko-c4",
				Amount: func(atk *core.AttackEvent, t core.Target) ([]float64, bool) {
					if atk.Info.AttackTag != core.AttackTagElementalArt {
						// only trigger on elemental art damage
						return nil, false
					}
					val := make([]float64, core.EndStatType)
					val[core.ElectroP] = 0.2
					return val, true
				},
			})
		}
		return false
	}, "yaemiko-c4")
}
