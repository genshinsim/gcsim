package reminiscence

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ShimenawasReminiscence, NewSet)
}

type Set struct {
	cd    int
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}
	s.cd = -1

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("shim-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	//11:51 AM] Episodde｜ShimenawaChildePeddler: Basically I found out that the fox set energy tax have around a 10 frame delay.
	//so I was testing if you can evade the fox set 15 energy tax by casting burst within those 10 frame after using an elemental
	//skill (not on hit). Turn out it work with childe :Childejoy:
	//The finding is now in #energy-drain-effects-have-a-delay if you want to take a closer look
	if count >= 4 {
		const icdKey = "shim-4pc-icd"
		icd := 600 // 10s * 60

		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.50
		c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
			if c.Player.Active() != char.Index {
				return false
			}
			if char.Energy < 15 {
				return false
			}
			if char.StatusIsActive(icdKey) {
				return false
			}
			char.AddStatus(icdKey, icd, true)

			// consume 15 energy, increased normal/charge/plunge dmg by 50%
			c.Tasks.Add(func() {
				char.AddEnergy("shim-4pc", -15)
			}, 10)

			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("shim-4pc", 60*10),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					switch atk.Info.AttackTag {
					case attacks.AttackTagNormal:
					case attacks.AttackTagExtra:
					case attacks.AttackTagPlunge:
					default:
						return nil, false
					}
					return m, true
				},
			})

			return false
		}, fmt.Sprintf("shim-4pc-%v", char.Base.Key.String()))

	}

	return &s, nil
}
