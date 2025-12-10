package silkenmoonsserenade

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	setKey2                      = "silken-moon-2pc"
	setKey4                      = "silken-moon-4pc"
	gleamingMoonDevotionEMKey    = "gleaming-moon-devotion-em"
	gleamingMoonDevotionReactKey = "gleaming-moon-devotion-reaction"
)

func init() {
	core.RegisterSetFunc(keys.SilkenMoonsSerenade, NewSet)
}

type Set struct {
	char  *character.CharWrapper
	core  *core.Core
	count int
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	if s.count < 2 {
		return nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(setKey2, -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if s.count < 4 {
		return nil
	}

	m2 := make([]float64, attributes.EndStatType)
	switch s.getMoonsignLevel() {
	case 0:
	case 1:
		m2[attributes.EM] = 60
	default:
		m2[attributes.EM] = 120
	}
	hook := func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.Element {
		case attributes.Pyro:
		case attributes.Hydro:
		case attributes.Electro:
		case attributes.Cryo:
		case attributes.Anemo:
		case attributes.Geo:
		case attributes.Dendro:
		default:
			return false
		}

		for _, char := range s.core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(gleamingMoonDevotionEMKey, 8*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					m2[attributes.EM] = 120
					return m2, true
				},
			})
		}

		return false
	}
	s.core.Events.Subscribe(event.OnEnemyDamage, hook, setKey4+"-dmg-"+s.char.Base.Key.String())

	for _, char := range s.core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(gleamingMoonDevotionReactKey, -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				switch ai.AttackTag {
				case attacks.AttackTagDirectLunarCharged:
				case attacks.AttackTagReactionLunarCharge:
				default:
					return 0, false
				}

				hasGleamingMoonDevotion := false
				for _, char1 := range s.core.Player.Chars() {
					if char1.StatModIsActive(gleamingMoonDevotionEMKey) {
						hasGleamingMoonDevotion = true
						break
					}
				}

				if !hasGleamingMoonDevotion {
					return 0, false
				}
				return 0.1, false
			},
		})
	}

	return nil
}

func (s *Set) getMoonsignLevel() int {
	count := 0
	for _, c := range s.core.Player.Chars() {
		count += c.Moonsign
	}
	return count
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		core:  c,
		count: count,
		Count: count,
	}
	return &s, nil
}
