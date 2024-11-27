package peakpatrolsong

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
	core.RegisterWeaponFunc(keys.PeakPatrolSong, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
}

const (
	buffKey     = "peakpatrolsong-buff"
	buffDur     = 6 * 60
	teamBuffKey = "peakpatrolsong-team-buff"
	teamBuffDur = 15 * 60
	icdKey      = "peakpatrolsong-buff-icd"
	icdDur      = 0.1 * 60
)

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := float64(p.Refine)

	m := make([]float64, attributes.EndStatType)
	t := make([]float64, attributes.EndStatType)

	selfDef := 0.06 + 0.02*r
	selfBonus := 0.075 + 0.025*r
	teamBonus := 0.06 + 0.02*r
	maxTeamBonus := 0.192 + 0.064*r

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal && atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 2 {
			w.stacks++
		}

		stacks := float64(w.stacks)
		m[attributes.DEFP] = selfDef * stacks
		bonus := selfBonus * stacks
		for i := attributes.PyroP; i <= attributes.DendroP; i++ {
			m[i] = bonus
		}
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(buffKey, buffDur),
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})

		if w.stacks == 2 {
			bonus := teamBonus * char.TotalDef(false) / 1000.0
			bonus = min(bonus, maxTeamBonus)
			for i := attributes.PyroP; i <= attributes.DendroP; i++ {
				t[i] = bonus
			}
			for _, this := range c.Player.Chars() {
				this.AddStatMod(character.StatMod{
					Base: modifier.NewBaseWithHitlag(teamBuffKey, teamBuffDur),
					Amount: func() ([]float64, bool) {
						return t, true
					},
				})
			}
		}

		char.AddStatus(icdKey, icdDur, true)
		return false
	}, fmt.Sprintf("peakpatrolsong-hit-%v", char.Base.Key.String()))

	return w, nil
}
