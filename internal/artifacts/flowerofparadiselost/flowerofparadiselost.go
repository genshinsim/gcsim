package flowerofparadiselost

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	icdKey = "flower-4pc-icd"
	icd    = 60 // 1s

	buffKey = "flower-4pc-buff"
)

func init() {
	core.RegisterSetFunc(keys.FlowerOfParadiseLost, NewSet)
}

type Set struct {
	stacks int
	Index  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2pc - Increases Elemental Mastery by 80.
// 4pc - The equipping character's Bloom, Hyperbloom, and Burgeon reaction DMG are increased by 40%. Additionally, after the equipping
//
//	character triggers Bloom, Hyperbloom, or Burgeon, they will gain another 25% bonus to the effect mentioned prior. Each stack
//	of this lasts 10s. Max 4 stacks simultaneously. This effect can only be triggered once per second. The character who equips
//	this can still trigger its effects when not on the field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("flower-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("flower-4pc", -1),
			Amount: func(ai combat.AttackInfo) (float64, bool) {
				switch ai.AttackTag {
				case attacks.AttackTagBloom:
				case attacks.AttackTagHyperbloom:
				case attacks.AttackTagBurgeon:
				default:
					return 0, false
				}
				return 0.4, false
			},
		})

		f := func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, icd, true)

			if !char.StatusIsActive(buffKey) {
				s.stacks = 0
			}
			if s.stacks < 4 {
				s.stacks++
			}

			c.Log.NewEvent("flower of paradise lost 4pc adding stack", glog.LogArtifactEvent, char.Index).
				Write("stacks", s.stacks)

			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag(buffKey, 10*60),
				Amount: func(ai combat.AttackInfo) (float64, bool) {
					switch ai.AttackTag {
					case attacks.AttackTagBloom:
					case attacks.AttackTagHyperbloom:
					case attacks.AttackTagBurgeon:
					default:
						return 0, false
					}
					return 0.4 * float64(s.stacks) * 0.25, false
				},
			})

			return false
		}

		c.Events.Subscribe(event.OnBloom, f, fmt.Sprintf("flower-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnHyperbloom, f, fmt.Sprintf("flower-4pc-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnBurgeon, f, fmt.Sprintf("flower-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
