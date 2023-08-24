package widsith

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheWidsith, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	mATK := make([]float64, attributes.EndStatType)
	mATK[attributes.ATKP] = .45 + float64(r)*0.15

	mEM := make([]float64, attributes.EndStatType)
	mEM[attributes.EM] = 180 + float64(r)*60

	mDmg := make([]float64, attributes.EndStatType)
	dmg := .36 + float64(r)*.12
	mDmg[attributes.PyroP] = dmg
	mDmg[attributes.HydroP] = dmg
	mDmg[attributes.CryoP] = dmg
	mDmg[attributes.ElectroP] = dmg
	mDmg[attributes.AnemoP] = dmg
	mDmg[attributes.GeoP] = dmg
	mDmg[attributes.DendroP] = dmg

	const icdKey = "widsith-icd"
	icd := 1800 // 30s * 60
	state := -1
	stats := []string{"em", "dmg%", "atk%"}
	buff := [][]float64{mEM, mDmg, mATK}

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		next := args[1].(int)

		if next != char.Index {
			return false
		}

		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)

		state = c.Rand.Intn(3)

		expiry := c.F + 60*10
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("widsith", 600),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				//sanity check; should never happen
				if state == -1 {
					return nil, false
				}
				return buff[state], true
			},
		})
		c.Log.NewEvent("widsith proc'd", glog.LogWeaponEvent, char.Index).
			Write("stat", stats[state]).
			Write("expiring (without hitlag)", expiry)

		return false
	}, fmt.Sprintf("width-%v", char.Base.Key.String()))

	return w, nil

}
