package thunderingfury

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ThunderingFury, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}
	icd := 0

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
	if count >= 4 {

		// add +0.4 reaction damage
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("tf-4pc", -1),
			Amount: func(ai combat.AttackInfo) (float64, bool) {
				// overload dmg can't melt or vape so it's fine
				switch ai.AttackTag {
				case combat.AttackTagOverloadDamage:
				case combat.AttackTagECDamage:
				case combat.AttackTagSuperconductDamage:
				default:
					return 0, false
				}
				return 0.4, false
			},
		})

		reduce := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if c.Player.Active() != char.Index {
				return false
			}
			if icd > c.F {
				return false
			}
			icd = c.F + 48
			char.ReduceActionCooldown(action.ActionSkill, 60)
			c.Log.NewEvent("thunderfury 4pc proc", glog.LogArtifactEvent, char.Index).
				Write("reaction", atk.Info.Abil).
				Write("new cd", char.Cooldown(action.ActionSkill))
			return false
		}

		c.Events.Subscribe(event.OnOverload, reduce, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnElectroCharged, reduce, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnSuperconduct, reduce, fmt.Sprintf("tf-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
