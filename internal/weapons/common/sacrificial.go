package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Sacrificial struct {
	Index int
	data  *model.WeaponData
}

func (s *Sacrificial) SetIndex(idx int)        { s.Index = idx }
func (s *Sacrificial) Init() error             { return nil }
func (s *Sacrificial) Data() *model.WeaponData { return s.data }

func NewSacrificial(data *model.WeaponData) core.NewWeaponFunc {
	s := &Sacrificial{data: data}
	return s.NewWeapon
}

func (s *Sacrificial) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	r := p.Refine

	const icdKey = "sacrificial-cd"

	prob := 0.3 + float64(r)*0.1

	cd := (34 - r*4) * 60

	if r >= 4 {
		cd = (19 - (r-4)*3) * 60
	}

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if char.Cooldown(action.ActionSkill) == 0 {
			return false
		}
		if dmg == 0 {
			return false
		}
		if c.Rand.Float64() < prob {
			char.ResetActionCooldown(action.ActionSkill)
			char.AddStatus(icdKey, cd, true)
			c.Log.NewEvent("sacrificial proc'd", glog.LogWeaponEvent, char.Index)
		}
		return false
	}, fmt.Sprintf("sac-%v", char.Base.Key.String()))

	return s, nil
}
