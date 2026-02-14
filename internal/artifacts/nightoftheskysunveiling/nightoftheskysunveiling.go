package nightoftheskysunveiling

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
	setKey2                    = "skys-unveiling-2pc"
	setKey4                    = "skys-unveiling-4pc"
	gleamingMoonIntentCRKey    = "gleaming-moon-intent-cr"
	gleamingMoonIntentReactKey = "gleaming-moon-intent-reaction"
)

func init() {
	core.RegisterSetFunc(keys.NightOfTheSkysUnveiling, NewSet)
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
	m[attributes.EM] = 80
	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(setKey2, -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	if s.count < 4 {
		return nil
	}

	m2 := make([]float64, attributes.EndStatType)
	switch s.core.Player.GetMoonsignLevel() {
	case 0:
	case 1:
		m2[attributes.CR] = 0.15
	default:
		m2[attributes.CR] = 0.30
	}

	lunarReactHook := func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		if s.core.Player.Active() != s.char.Index() {
			return false
		}

		s.char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(gleamingMoonIntentCRKey, 4*60),
			AffectedStat: attributes.CR,
			Amount: func() ([]float64, bool) {
				return m2, true
			},
		})

		return false
	}
	s.core.Events.Subscribe(event.OnLunarCharged, lunarReactHook, setKey4+"-lc-"+s.char.Base.Key.String())
	s.core.Events.Subscribe(event.OnLunarBloom, lunarReactHook, setKey4+"-lc-"+s.char.Base.Key.String())

	for _, char := range s.core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(gleamingMoonIntentReactKey, -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if !attacks.AttackTagIsLunar(ai.AttackTag) {
					return 0, false
				}

				hasGleamingMoonIntent := false
				for _, char1 := range s.core.Player.Chars() {
					if char1.StatModIsActive(gleamingMoonIntentCRKey) {
						hasGleamingMoonIntent = true
						break
					}
				}

				if !hasGleamingMoonIntent {
					return 0, false
				}
				return 0.1, false
			},
		})
	}

	return nil
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
