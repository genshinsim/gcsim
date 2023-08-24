package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Blackcliff struct {
	Index int
}

func (b *Blackcliff) SetIndex(idx int) { b.Index = idx }
func (b *Blackcliff) Init() error      { return nil }

func NewBlackcliff(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {

	b := &Blackcliff{}

	atk := 0.09 + float64(p.Refine)*0.03
	index := 0
	stackKey := []string{
		"blackcliff-stack-1",
		"blackcliff-stack-2",
		"blackcliff-stack-3",
	}
	m := make([]float64, attributes.EndStatType)

	amtfn := func() ([]float64, bool) {
		count := 0
		for _, v := range stackKey {
			if char.StatusIsActive(v) {
				count++
			}
		}
		m[attributes.ATKP] = atk * float64(count)
		return m, true
	}

	c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
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
		//add status to char given index
		char.AddStatus(stackKey[index], 1800, true)
		//update buff
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("blackcliff", 1800),
			AffectedStat: attributes.ATKP,
			Amount:       amtfn,
		})
		index++
		if index == 3 {
			index = 0
		}
		return false
	}, fmt.Sprintf("blackcliff-%v", char.Base.Key.String()))

	return b, nil
}
