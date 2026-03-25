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
	buff     float64
	Index    int
	Count    int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	m := 0.2
	if s.core.Player.GetMoonsignLevel() >= 2 {
		m += 0.4
	}
	if s.Count >= 4 && s.core.Player.Active() != s.char.Index() {
		s.gainBuff(m)
	}
	return nil
}

func NewSet(core *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core:     core,
		char:     char,
		lastSwap: -1,
		Count:    count,
	}
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.EM] = 80
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("aubade-2pc", -1),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}
	if count >= 4 {
		m := 0.2
		if core.Player.GetMoonsignLevel() >= 2 {
			m += 0.4
		}

		core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
			prev := args[0].(int)
			next := args[1].(int)
			if prev == char.Index() {
				s.lastSwap = -1
				s.gainBuff(m)
			} else if next == char.Index() {
				s.lastSwap = core.F
				core.Tasks.Add(s.clearBuff(core.F), 3*60)
			}
		}, fmt.Sprintf("aubade-4pc-%v", char.Base.Key.String()))

		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("aubade-4pc", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if ai.AttackTag > attacks.LunarReactionStartDelim && ai.AttackTag < attacks.DirectLunarReactionEndDelim {
					return s.buff
				}
				return 0
			},
		})
	}

	return &s, nil
}

func (s *Set) gainBuff(m float64) {
	s.buff = m
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

		s.buff = 0
		s.core.Log.NewEvent("aubade of morningstar and moon 4pc lost", glog.LogArtifactEvent, s.char.Index())
	}
}
