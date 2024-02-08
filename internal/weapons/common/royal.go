package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const icdKey = "royal-icd"

type Royal struct {
	Index int
	data  *model.WeaponData
}

func (r *Royal) SetIndex(idx int)        { r.Index = idx }
func (r *Royal) Init() error             { return nil }
func (r *Royal) Data() *model.WeaponData { return r.data }

func NewRoyal(data *model.WeaponData) *Royal {
	return &Royal{data: data}
}

func (r *Royal) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	refine := p.Refine

	stacks := 0

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		crit := args[3].(bool)
		if dmg == 0 {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if crit {
			stacks = 0
		} else if !char.StatusIsActive(icdKey) {
			stacks++
			if stacks > 5 {
				stacks = 5
			}
			char.AddStatus(icdKey, 18, true)
			c.Log.NewEvent("royal stacked", glog.LogWeaponEvent, char.Index).
				Write("stacks", stacks)
		}
		return false
	}, fmt.Sprintf("royal-%v", char.Base.Key.String()))

	rate := 0.06 + float64(refine)*0.02
	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("royal", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			m[attributes.CR] = float64(stacks) * rate
			return m, true
		},
	})

	return r, nil
}
