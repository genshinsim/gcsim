package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Blackcliff struct {
	Index int
}

func (b *Blackcliff) SetIndex(idx int) { b.Index = idx }
func (b *Blackcliff) Init() error      { return nil }

func NewBlackcliff(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (*Blackcliff, error) {
	b := &Blackcliff{}

	atk := 0.09 + float64(p.Refine)*0.03
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, attributes.EndStatType)
	char.AddStatMod(
		"blackcliff",
		-1,
		attributes.NoStat,
		func() ([]float64, bool) {
			count := 0
			for _, v := range stacks {
				if v > c.F {
					count++
				}
			}
			m[attributes.ATKP] = atk * float64(count)
			return m, true
		},
	)

	c.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		stacks[index] = c.F + 1800
		index++
		if index == 3 {
			index = 0
		}
		return false
	}, fmt.Sprintf("blackcliff-%v", char.Base.Name))

	return b, nil
}

type Favonius struct {
	Index int
}

func (b *Favonius) SetIndex(idx int) { b.Index = idx }
func (b *Favonius) Init() error      { return nil }

func NewFavonius(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (*Favonius, error) {
	b := &Favonius{}

	prob := 0.50 + float64(p.Refine)*0.1
	cd := 810 - p.Refine*90
	icd := 0
	//add on crit effect
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if icd > c.F {
			return false
		}

		if c.Rand.Float64() > prob {
			return false
		}
		c.Log.NewEvent("favonius proc'd", glog.LogWeaponEvent, char.Index)

		c.QueueParticle("favonius-"+char.Base.Name, 3, attributes.NoElement, 80)

		icd = c.F + cd

		return false
	}, fmt.Sprintf("favo-%v", char.Base.Name))

	return b, nil
}

type NoEffect struct {
	Index int
}

func (b *NoEffect) SetIndex(idx int) { b.Index = idx }
func (b *NoEffect) Init() error      { return nil }
func NewNoEffect(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (*NoEffect, error) {
	return &NoEffect{}, nil
}
