package aubade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.AubadeOfMorningstarAndMoon, NewSet)
}

type Set struct {
	core  *core.Core
	char  *character.CharWrapper
	Index int
	Count int
}

const aubade4pcOnFieldKey = "aubade-4pc-onfield"

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }

func (s *Set) Init() error {
	// NewSet is called before all characters added, so we don't initialize buffs there due to MoonsignLevel being incorrect
	// Init gets called after all characters added, so MoonsignLevel is correct.

	if s.Count < 2 {
		return nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 80
	s.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("aubade-2pc", -1),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			return m
		},
	})

	if s.Count < 4 {
		return nil
	}

	buff := 0.2
	if s.core.Player.GetMoonsignLevel() >= 2 {
		buff += 0.4
	}

	s.core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == s.char.Index() {
			s.char.DeleteStatus(aubade4pcOnFieldKey)
		} else if next == s.char.Index() {
			s.char.AddStatus(aubade4pcOnFieldKey, 3*60, true)
		}
	}, fmt.Sprintf("aubade-4pc-%v", s.char.Base.Key.String()))

	s.char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("aubade-4pc-dmg", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if s.core.Player.Active() == s.char.Index() && !s.char.StatusIsActive(aubade4pcOnFieldKey) {
				return 0
			}
			if attacks.AttackTagIsLunar(ai.AttackTag) {
				return buff
			}
			return 0
		},
	})

	return nil
}

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:  core,
		char:  char,
		Count: count,
	}
	return &s, nil
}
