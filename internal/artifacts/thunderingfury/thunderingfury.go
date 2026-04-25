package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ThunderingFury, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 2pc - Electro DMG Bonus +15%
// 4pc - Increases DMG caused by Overloaded, Electro-Charged, Superconduct, and Hyperbloom by 40%,
// and the DMG Bonus conferred by Aggravate is increased by 20%, and the DMG caused by Lunar-Charged by 20%.
// When Quicken or the aforementioned Elemental Reactions are triggered, Elemental Skill CD is decreased by
// 1s. Can only occur once every 0.8s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ElectroP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("tf-2pc", -1),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				return m
			},
		})
	}

	if count < 4 {
		return &s, nil
	}

	const icdKey = "tf-4pc-icd"
	icd := 48 // 0.8s * 60

	char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("tf-4pc", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.Catalyzed && ai.CatalyzedType == info.ReactionTypeAggravate {
				return 0.2
			}
			switch ai.AttackTag {
			case attacks.AttackTagOverloadDamage,
				attacks.AttackTagECDamage,
				attacks.AttackTagSuperconductDamage,
				attacks.AttackTagHyperbloom:
				return 0.4
			case attacks.AttackTagDirectLunarCharged, attacks.AttackTagReactionLunarCharge:
				return 0.2
			}
			return 0
		},
	})

	reduce := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if c.Player.Active() != char.Index() {
			return
		}
		if char.StatusIsActive(icdKey) {
			return
		}
		char.AddStatus(icdKey, icd, true)
		char.ReduceActionCooldown(action.ActionSkill, 60)
		c.Log.NewEvent("thunderfury 4pc proc", glog.LogArtifactEvent, char.Index()).
			Write("reaction", atk.Info.Abil).
			Write("new cd", char.Cooldown(action.ActionSkill))
	}
	reduceNoGadget := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); ok {
			reduce(args...)
		}
	}

	c.Events.Subscribe(event.OnOverload, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnElectroCharged, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnLunarCharged, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSuperconduct, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnHyperbloom, reduce, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnQuicken, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnAggravate, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
