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
	count int
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func (s *Set) getMoonsignLevel(core *core.Core) int {
	count := 0
	for _, c := range core.Player.Chars() {
		count += c.Moonsign
	}
	return count
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		char:  char,
		count: count,
		Count: count,
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(setKey2, -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if count < 4 {
		return &s, nil
	}

	m2 := make([]float64, attributes.EndStatType)
	switch s.getMoonsignLevel(c) {
	case 0:
		return &s, nil
	case 1:
		m2[attributes.EM] = 60
	default:
		m2[attributes.EM] = 120
	}

	lunarReactHook := func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		for _, char := range c.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBase(gleamingMoonDevotionEMKey, 8*60),
				AffectedStat: attributes.EM,
				Amount: func() ([]float64, bool) {
					return m2, true
				},
			})
		}

		return false
	}
	c.Events.Subscribe(event.OnLunarCharged, lunarReactHook, setKey4+"-lc-"+char.Base.Key.String())

	for _, char := range c.Player.Chars() {
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
				for _, char1 := range c.Player.Chars() {
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

	return &s, nil
}
