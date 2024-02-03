package common

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"

	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Lithic struct {
	Index int
	data  *model.WeaponData
}

func (l *Lithic) SetIndex(idx int)        { l.Index = idx }
func (l *Lithic) Init() error             { return nil }
func (l *Lithic) Data() *model.WeaponData { return l.data }

func NewLithic(data *model.WeaponData) *Lithic {
	return &Lithic{data: data}
}

func (l *Lithic) NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	r := p.Refine

	stacks := 0
	val := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnInitialize, func(args ...interface{}) bool {
		for _, char := range c.Player.Chars() {
			if char.CharZone == info.ZoneLiyue {
				stacks++
			}
		}
		val[attributes.CR] = (0.02 + float64(r)*0.01) * float64(stacks)
		val[attributes.ATKP] = (0.06 + float64(r)*0.01) * float64(stacks)
		return true
	}, fmt.Sprintf("lithic-%v", char.Base.Key.String()))
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("lithic", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})

	return l, nil
}
