package gambler

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.Gambler, NewSet)
}

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// 4-Piece Bonus: Resets Skill CD after defeating an enemy
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}

	// 2 Piece: Increases Elemental Skill DMG by 20%.
	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.20
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("gambler-2pc", -1),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagElementalArt {
					return nil, false
				}
				return m, true
			},
		})
	}

	// 4 Piece: Defeating an opponent has 100% chance to remove Elemental Skill CD. Can only occur once every 15s.
	if count >= 4 {
		const icdKey = "gambler-4pc-icd"
		c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
			// don't proc if on icd
			if char.StatusIsActive(icdKey) {
				return false
			}
			_, ok := args[0].(*enemy.Enemy)
			// ignore if not an enemy
			if !ok {
				return false
			}
			atk := args[1].(*combat.AttackEvent)
			// don't proc if someone else defeated the enemy
			if atk.Info.ActorIndex != char.Index {
				return false
			}
			// don't proc if off-field
			if c.Player.Active() != char.Index {
				return false
			}

			// reset skill cd
			char.ResetActionCooldown(action.ActionSkill)
			c.Log.NewEvent("gambler-4pc proc'd", glog.LogArtifactEvent, char.Index)

			// set icd
			char.AddStatus(icdKey, 900, true) // 15s

			return false
		}, fmt.Sprintf("gambler-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
