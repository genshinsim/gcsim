package finaleofthedeepgalleries

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

func init() {
	core.RegisterSetFunc(keys.FinaleOfTheDeepGalleries, NewSet)
}

type Set struct {
	Index int
	Count int

	c                 *core.Core
	char              *character.CharWrapper
	procNormalExpireF int
	procBurstExpireF  int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
		c:     c,
		char:  char,
	}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.CryoP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("deep-galleries-2pc", -1),
			AffectedStat: attributes.CryoP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count >= 4 {
		procDurNormal := 360 // 6s * 60
		procDurBurst := 360  // 6s * 60

		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.6

		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("deep-galleries-4pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if char.Energy != 0 {
					return nil, false
				}
				if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
					return nil, false
				}
				if atk.Info.AttackTag == attacks.AttackTagNormal && c.F < s.procNormalExpireF {
					return nil, false
				}
				if atk.Info.AttackTag == attacks.AttackTagElementalBurst && c.F < s.procBurstExpireF {
					return nil, false
				}
				return m, true
			},
		})

		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			// If attack does not belong to the equipped character then ignore
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			// If this is not a normal attack then ignore
			if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return false
			}

			s.procNormalExpireF = c.F + procDurNormal
			s.procBurstExpireF = c.F + procDurBurst

			c.Log.NewEvent("deep galleries 4pc stop playing", glog.LogArtifactEvent, char.Index).
				Write("normal_buff_stop_expiry", s.procNormalExpireF).
				Write("burst_buff_stop_expiry", s.procBurstExpireF)

			return false
		}, fmt.Sprintf("deep-galleries-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
