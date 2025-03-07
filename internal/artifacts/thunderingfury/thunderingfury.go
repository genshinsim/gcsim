package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
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
// and the DMG Bonus conferred by Aggravate is increased by 20%. When Quicken or the aforementioned
// Elemental Reactions are triggered, Elemental Skill CD is decreased by 1s. Can only occur once every 0.8s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ElectroP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("tf-2pc", -1),
			AffectedStat: attributes.ElectroP,
			Amount: func() ([]float64, bool) {
				return m, true
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
		Amount: func(ai combat.AttackInfo) (float64, bool) {
			if ai.Catalyzed && ai.CatalyzedType == reactions.Aggravate {
				return 0.2, false
			}
			switch ai.AttackTag {
			case attacks.AttackTagOverloadDamage,
				attacks.AttackTagECDamage,
				attacks.AttackTagSuperconductDamage,
				attacks.AttackTagHyperbloom:
				return 0.4, false
			}
			return 0, false
		},
	})

	//nolint:unparam // ignoring for now, event refactor should get rid of bool return of event sub
	reduce := func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		char.ReduceActionCooldown(action.ActionSkill, 60)
		c.Log.NewEvent("thunderfury 4pc proc", glog.LogArtifactEvent, char.Index).
			Write("reaction", atk.Info.Abil).
			Write("new cd", char.Cooldown(action.ActionSkill))
		return false
	}

	reduceNoGadget := func(args ...interface{}) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		return reduce(args...)
	}

	c.Events.Subscribe(event.OnOverload, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnElectroCharged, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSuperconduct, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnHyperbloom, reduce, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnQuicken, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnAggravate, reduceNoGadget, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
