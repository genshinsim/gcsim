package echoes

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
	core.RegisterSetFunc(keys.EchoesOfAnOffering, NewSet)
}

type Set struct {
	prob        float64
	icd         int
	procExpireF int
	Index       int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2pc - ATK +18%.
// 4pc - When Normal Attacks hit opponents, there is a 36% chance that it will trigger Valley Rite, which will increase Normal Attack DMG by 70% of ATK.
//
//	This effect will be dispelled 0.05s after a Normal Attack deals DMG.
//	If a Normal Attack fails to trigger Valley Rite, the odds of it triggering the next time will increase by 20%.
//	This trigger can occur once every 0.2s.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	procDuration := 3 // 0.05s

	s := Set{}
	s.prob = 0.36
	s.icd = 0
	s.procExpireF = 0

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.18
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("echoes-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	if count < 4 {
		return &s, nil
	}

	var dmgAdded float64

	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		// if the active char is not the equipped char then ignore
		if c.Player.Active() != char.Index {
			return false
		}

		// If attack does not belong to the equipped character then ignore
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// If this is not a normal attack then ignore
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return false
		}

		// if buff is already active then buff attack
		snATK := char.TotalAtk()
		if c.F < s.procExpireF {
			dmgAdded = snATK * 0.7
			atk.Info.FlatDmg += dmgAdded
			c.Log.NewEvent("echoes 4pc adding dmg", glog.LogArtifactEvent, char.Index).
				Write("dmg_added", dmgAdded).
				Write("buff_expiry", s.procExpireF).
				Write("icd_up", s.icd)

			return false
		}

		// If Artifact set effect is still on CD then ignore
		if c.F < s.icd {
			c.Log.NewEvent("echoes 4pc failed to proc due icd", glog.LogArtifactEvent, char.Index).
				Write("icd_up", s.icd)
			return false
		}

		if c.Rand.Float64() > s.prob {
			s.icd = c.F + 12 // 0.2s
			s.prob += 0.2
			c.Log.NewEvent("echoes 4pc failed to proc due to chance", glog.LogArtifactEvent, char.Index).
				Write("probability_now", s.prob).
				Write("icd_up", s.icd)
			return false
		}

		dmgAdded = snATK * 0.7
		atk.Info.FlatDmg += dmgAdded

		s.procExpireF = c.F + procDuration
		s.icd = c.F + 12 // 0.2s

		s.prob = 0.36

		c.Log.NewEvent("echoes 4pc adding dmg", glog.LogArtifactEvent, char.Index).
			Write("dmg_added", dmgAdded).
			Write("buff_expiry", s.procExpireF).
			Write("icd_up", s.icd)

		return false
	}, fmt.Sprintf("echoes-4pc-%v", char.Base.Key.String()))

	return &s, nil
}
