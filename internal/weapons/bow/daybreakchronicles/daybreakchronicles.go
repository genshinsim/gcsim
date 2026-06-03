package daybreakchronicles

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheDaybreakChronicles, NewWeapon)
}

type Weapon struct {
	Index                                int
	stirringDawnBreezeNormalAttackStacks int
	stirringDawnBreezeNormalAttackSource int
	stirringDawnBreezeSkillStacks        int
	stirringDawnBreezeSkillSource        int
	stirringDawnBreezeBurstStacks        int
	stirringDawnBreezeBurstSource        int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	w.stirringDawnBreezeNormalAttackStacks = 6
	w.stirringDawnBreezeSkillStacks = 6
	w.stirringDawnBreezeBurstStacks = 6
	w.stirringDawnBreezeNormalAttackSource = c.F
	w.stirringDawnBreezeSkillSource = c.F
	w.stirringDawnBreezeBurstSource = c.F

	m := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("stirring-dawn-breeze", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
				m[attributes.DmgP] = (0.075 + float64(r)*0.025) * float64(w.stirringDawnBreezeNormalAttackStacks)
			case attacks.AttackTagElementalArt:
				m[attributes.DmgP] = (0.075 + float64(r)*0.025) * float64(w.stirringDawnBreezeSkillStacks)
			case attacks.AttackTagElementalArtHold:
				m[attributes.DmgP] = (0.075 + float64(r)*0.025) * float64(w.stirringDawnBreezeSkillStacks)
			case attacks.AttackTagElementalBurst:
				m[attributes.DmgP] = (0.075 + float64(r)*0.025) * float64(w.stirringDawnBreezeBurstStacks)
			default:
				return nil
			}

			return m
		},
	})

	daybreakAddBuff := func(stacks *int, stacksSrc *int) {
		if *stacks >= 6 {
			*stacksSrc = c.F
			c.Tasks.Add(w.decreaseStacks(c, c.F, stacks, stacksSrc), 60)
			return
		}
		*stacks++

		if c.Player.GetHexereiCount() < 2 {
			return
		}

		if *stacks >= 6 {
			*stacksSrc = c.F
			c.Tasks.Add(w.decreaseStacks(c, c.F, stacks, stacksSrc), 60)
			return
		}
		*stacks++
	}

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			daybreakAddBuff(&w.stirringDawnBreezeNormalAttackStacks, &w.stirringDawnBreezeNormalAttackSource)
			c.Log.NewEvent("daybreak-chronicles-normal", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.stirringDawnBreezeNormalAttackStacks).
				Write("src", w.stirringDawnBreezeNormalAttackSource)
		case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			daybreakAddBuff(&w.stirringDawnBreezeSkillStacks, &w.stirringDawnBreezeSkillSource)
			c.Log.NewEvent("daybreak-chronicles-skill", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.stirringDawnBreezeSkillStacks).
				Write("src", w.stirringDawnBreezeSkillSource)
		case attacks.AttackTagElementalBurst:
			daybreakAddBuff(&w.stirringDawnBreezeBurstStacks, &w.stirringDawnBreezeBurstSource)
			c.Log.NewEvent("daybreak-chronicles-burst", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.stirringDawnBreezeBurstStacks).
				Write("src", w.stirringDawnBreezeBurstSource)
		}
	}, fmt.Sprintf("daybreak-chronicles-buff-refresh-%v", char.Base.Key.String()))

	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.stirringDawnBreezeNormalAttackStacks, &w.stirringDawnBreezeNormalAttackSource), 60)
	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.stirringDawnBreezeSkillStacks, &w.stirringDawnBreezeSkillSource), 60)
	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.stirringDawnBreezeBurstStacks, &w.stirringDawnBreezeBurstSource), 60)

	return w, nil
}

func (w *Weapon) decreaseStacks(c *core.Core, src int, stacks, stacksSrc *int) func() {
	return func() {
		if *stacks == 0 || src != *stacksSrc {
			return
		}

		*stacks--

		if *stacks != 0 {
			c.Tasks.Add(w.decreaseStacks(c, src, stacks, stacksSrc), 60)
		}
	}
}
