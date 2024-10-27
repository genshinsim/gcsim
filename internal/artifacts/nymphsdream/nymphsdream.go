package nymphsdream

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.NymphsDream, NewSet)
}

type Set struct {
	key   string
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.HydroP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("nd-2pc", -1),
			AffectedStat: attributes.HydroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count < 4 {
		return &s, nil
	}

	s.key = fmt.Sprintf("%v-nd-4pc", char.Base.Key.String())

	const normalKey = "nd-normal"
	const chargedKey = "nd-charged"
	const skillKey = "nd-skill"
	const burstKey = "nd-burst"
	const plungeKey = "nd-plunge"

	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("nd-4pc", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			stacks := 0
			for _, k := range []string{
				normalKey, chargedKey, plungeKey,
				skillKey, burstKey,
			} {
				if stacks < 3 && char.StatusIsActive(k) {
					stacks++
				}
			}

			if stacks > 0 {
				m[attributes.ATKP] = 0.09*float64(stacks) - 0.02
				m[attributes.HydroP] = 0.04
				if stacks > 1 {
					m[attributes.HydroP] += 0.05
				}
				if stacks > 2 {
					m[attributes.HydroP] += 0.06
				}
			} else {
				m[attributes.ATKP] = 0
				m[attributes.HydroP] = 0
			}

			return m, true
		},
	})

	const stackDuration = 480

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			char.AddStatus(normalKey, stackDuration, true)
		case attacks.AttackTagExtra:
			char.AddStatus(chargedKey, stackDuration, true)
		case attacks.AttackTagPlunge:
			char.AddStatus(plungeKey, stackDuration, true)
		case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			char.AddStatus(skillKey, stackDuration, true)
		case attacks.AttackTagElementalBurst:
			char.AddStatus(burstKey, stackDuration, true)
		}

		return false
	}, s.key)

	return &s, nil
}
