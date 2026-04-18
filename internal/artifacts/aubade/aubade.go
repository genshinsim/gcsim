package aubade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.AubadeOfMorningstarAndMoon, NewSet)
}

type Set struct {
	lastSwap int
	core     *core.Core
	char     *character.CharWrapper
	Index    int
	Count    int
}

const aubade4pcKey = "aubade-4pc"

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

	if s.core.Player.Active() != s.char.Index() {
		s.gainBuff()
	}

	s.core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)
		next := args[1].(int)
		if prev == s.char.Index() {
			s.lastSwap = -1
			s.gainBuff()
		} else if next == s.char.Index() {
			s.lastSwap = s.core.F
			s.char.AddStatus(aubade4pcKey, 3*60, false)
			s.core.Tasks.Add(s.clearBuff(s.core.F), 3*60)
		}
	}, fmt.Sprintf("aubade-4pc-%v", s.char.Base.Key.String()))

	s.char.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase("aubade-4pc-dmg", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if !s.char.StatusIsActive(aubade4pcKey) {
				return 0
			}
			if attacks.LunarReactionStartDelim < ai.AttackTag && ai.AttackTag < attacks.DirectLunarReactionEndDelim {
				return buff
			}
			return 0
		},
	})

	return nil
}

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:     core,
		char:     char,
		lastSwap: -1,
		Count:    count,
	}
	return &s, nil
}

func (s *Set) gainBuff() {
	s.char.AddStatus(aubade4pcKey, -1, false)
	s.core.Log.NewEvent("aubade of morningstar and moon 4pc proc'd", glog.LogArtifactEvent, s.char.Index())
}

func (s *Set) clearBuff(src int) func() {
	return func() {
		if s.lastSwap != src {
			return
		}
		if s.core.Player.Active() != s.char.Index() {
			return
		}

		s.core.Log.NewEvent("aubade of morningstar and moon 4pc lost", glog.LogArtifactEvent, s.char.Index())
	}
}
