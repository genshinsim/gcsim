package ultimateoverlordsmegamagicsword

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.UltimateOverlordsMegaMagicSword, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// ATK increased by 12/15/18/21/24%. That's not all!
// The support from all Melusines you've helped in Merusea Village fills you with strength!
// Based on the number of them you've helped, your ATK is increased by up to an additional 12/15/18/21/24%.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// 6 melusines are enough to cap out the buff
	melusinesCap := 6
	melusines, ok := p.Params["melusines"]
	if !ok {
		melusines = 6 // default is max buff
	} else {
		melusines = min(melusines, melusinesCap)
	}

	// floating point division otherwise will either be 0 or 1
	additional := max(float64(melusines)/float64(melusinesCap), 0)

	// perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = (0.09 + float64(r)*0.03) * (1 + additional)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("ultimateoverlordsmegamagicsword", -1),
		AffectedStat: attributes.ATKP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	return w, nil
}
