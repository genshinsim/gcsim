package homa

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.StaffOfHoma, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	mHP := make([]float64, attributes.EndStatType)
	mHP[attributes.HPP] = 0.15 + float64(r)*0.05
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("homa-hp", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return mHP, true
		},
	})

	mATK := make([]float64, attributes.EndStatType)
	atkp := 0.006 + float64(r)*0.002
	lowhp := 0.008 + float64(r)*0.002
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("homa-atk-buff", -1),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			maxhp := char.MaxHP()
			per := atkp
			if char.CurrentHPRatio() <= 0.5 {
				per += lowhp
			}
			mATK[attributes.ATK] = per * maxhp
			return mATK, true
		},
	})

	return w, nil
}
