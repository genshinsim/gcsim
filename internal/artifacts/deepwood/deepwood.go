package deepwood

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.DeepwoodMemories, NewSet)
}

type Set struct {
	Index int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) Init() error      { return nil }

// 2-Piece Bonus: Dendro DMG Bonus +15%.
// 4-Piece Bonus: After Elemental Skills or Bursts hit opponents, the targets’ Dendro RES will be decreased by 30% for 8s.
// This effect can be triggered even if the equipping character is not on the field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (artifact.Set, error) {
	s := Set{}

	if count >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DendroP] = 0.15
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("dm-2pc", -1),
			AffectedStat: attributes.DendroP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}
	if count >= 4 {
		c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			atk := args[1].(*combat.AttackEvent)
			t, ok := args[0].(*enemy.Enemy)
			if !ok {
				return false
			}
			if atk.Info.ActorIndex != char.Index {
				return false
			}

			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold && atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return false
			}

			t.AddResistMod(combat.ResistMod{
				Base:  modifier.NewBaseWithHitlag("dm-4pc", 8*60),
				Ele:   attributes.Dendro,
				Value: -0.3,
			})
			c.Log.NewEvent("dm 4pc proc", glog.LogArtifactEvent, char.Index).Write("char", char.Index)

			return false
		}, fmt.Sprintf("dm-4pc-%v", char.Base.Key.String()))
	}

	return &s, nil
}
