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
	Index       int
	naStacks    int
	naSrc       int
	skillStacks int
	skillSrc    int
	burstStacks int
	burstSrc    int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	w.naStacks = 6
	w.skillStacks = 6
	w.burstStacks = 6

	m := make([]float64, attributes.EndStatType)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("stirring-dawn-breeze", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			dmgPerStack := 0.075 + float64(r)*0.025
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal:
				m[attributes.DmgP] = dmgPerStack * float64(w.naStacks)
			case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
				m[attributes.DmgP] = dmgPerStack * float64(w.skillStacks)
			case attacks.AttackTagElementalBurst:
				m[attributes.DmgP] = dmgPerStack * float64(w.burstStacks)
			default:
				return nil
			}

			return m
		},
	})

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			w.daybreakAddBuff(c, &w.naStacks, &w.naSrc)
			c.Log.NewEvent("daybreak-chronicles-normal", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.naStacks)
		case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			w.daybreakAddBuff(c, &w.skillStacks, &w.skillSrc)
			c.Log.NewEvent("daybreak-chronicles-skill", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.skillStacks)
		case attacks.AttackTagElementalBurst:
			w.daybreakAddBuff(c, &w.burstStacks, &w.burstSrc)
			c.Log.NewEvent("daybreak-chronicles-burst", glog.LogWeaponEvent, char.Index()).
				Write("stacks", w.burstStacks)
		}
	}, fmt.Sprintf("daybreak-chronicles-buff-refresh-%v", char.Base.Key.String()))

	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.naStacks, &w.naSrc), 60)
	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.skillStacks, &w.skillSrc), 60)
	c.Tasks.Add(w.decreaseStacks(c, c.F, &w.burstStacks, &w.burstSrc), 60)

	return w, nil
}

func (w *Weapon) decreaseStacks(c *core.Core, src int, stacks, stacksSrc *int) func() {
	return func() {
		if *stacks == 0 {
			return
		}

		if src != *stacksSrc {
			return
		}

		*stacks--

		if *stacks != 0 {
			c.Tasks.Add(w.decreaseStacks(c, src, stacks, stacksSrc), 60)
		}
	}
}

func (w *Weapon) daybreakAddBuff(c *core.Core, stacks, stacksSrc *int) {
	*stacksSrc = c.F
	c.Tasks.Add(w.decreaseStacks(c, c.F, stacks, stacksSrc), 60)

	if *stacks >= 6 {
		return
	}
	*stacks++

	if c.Player.GetHexereiCount() < 2 {
		return
	}

	if *stacks >= 6 {
		return
	}
	*stacks++
}
