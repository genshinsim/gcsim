package songofdayspast

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
	yearningKey        = "yearning-effect"
	healStacksKey      = "sodp-heal-stacks"
	healSnapKey        = "sodp-heal-snap"
	wavesOfDaysPastKey = "waves-of-days-past"
)

func init() {
	core.RegisterSetFunc(keys.SongOfDaysPast, NewSet)
}

type Set struct {
	Index int
	core  *core.Core
	char  *character.CharWrapper
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		core: c,
		char: char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.Heal] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("sodp-2pc", -1),
			AffectedStat: attributes.Heal,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		c.Events.Subscribe(event.OnHeal, s.OnHeal(), fmt.Sprintf("sodp-4pc-heal-accumulation-%v", char.Base.Key.String()))
		c.Events.Subscribe(event.OnEnemyHit, s.OnEnemyHit(), fmt.Sprintf("waves-of-days-past-%v", char.Base.Key.String()))
	}

	return &s, nil
}

func (s *Set) OnHeal() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		src := args[0].(*info.HealInfo)
		healAmt := args[2].(float64)
		if src.Caller != s.char.Index {
			return false
		}
		s.core.Flags.Custom[healStacksKey] += healAmt
		if s.core.Flags.Custom[healStacksKey] >= 15000 {
			s.core.Flags.Custom[healStacksKey] = 15000
		}
		if s.core.Status.Duration(yearningKey) == 0 {
			s.core.Status.Add(yearningKey, 6*60)
			s.core.Tasks.Add(func() {
				s.core.Flags.Custom[healSnapKey] = s.core.Flags.Custom[healStacksKey]
				s.core.Flags.Custom[healStacksKey] = 0
				s.core.Status.Add(wavesOfDaysPastKey, 10*60)
				s.core.Flags.Custom[wavesOfDaysPastKey] = 5
			}, 6*60)
		}
		return false
	}
}

func (s *Set) OnEnemyHit() func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}
		if s.core.Status.Duration(wavesOfDaysPastKey) == 0 {
			return false
		}
		if atk.Info.ActorIndex != s.core.Player.Active() {
			return false
		}
		if s.core.Flags.Custom[wavesOfDaysPastKey] > 0 {
			s.core.Flags.Custom[wavesOfDaysPastKey]--
			amt := s.core.Flags.Custom[healSnapKey] * 0.08
			atk.Info.FlatDmg += amt
			s.core.Log.NewEvent("sodp 4pc adding dmg", glog.LogArtifactEvent, atk.Info.ActorIndex).
				Write("dmg_added", amt)
		}
		if s.core.Flags.Custom[wavesOfDaysPastKey] == 0 {
			s.core.Status.Delete(wavesOfDaysPastKey)
		}
		return false
	}
}
