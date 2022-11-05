package alley

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.TheAlleyFlash, NewWeapon)
}

// Upon damaging an opponent, increases CRIT Rate by 8/10/12/14/16%. Max 5 stacks. A CRIT Hit removes all stacks.
type Weapon struct {
	Index int
	c     *core.Core
	char  *character.CharWrapper
}

const lockoutKey = "alley-flash-lockout"

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func (w *Weapon) selfDisable(lambda float64) func() {
	return func() {
		//disable for 5 sec
		w.char.AddStatus(lockoutKey, 300, true)
		//-ln(U)/lambda` (where U~Uniform[0,1]).
		next := int(math.Log(w.c.Rand.Float64()) / lambda)
		w.c.Tasks.Add(w.selfDisable(lambda), next)
	}
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{
		c:    c,
		char: char,
	}
	r := p.Refine

	//allow user to periodically lock out this weapon (just to screw around with bennett)
	//follows poisson distribution, user provides lambda:
	//https://stackoverflow.com/questions/6527345/simulating-poisson-waiting-times
	if lambda, ok := p.Params["lambda"]; ok {
		//user supplied lambda should be per min, so we need to scale this down by *60*60
		l := float64(lambda) / 3600.0
		//queue tasks to disable
		next := int(-math.Log(1-w.c.Rand.Float64()) / l)
		c.Tasks.Add(w.selfDisable(l), next)
	}

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		w.char.AddStatus(lockoutKey, 300, true)
		return false
	}, fmt.Sprintf("alleyflash-%v", char.Base.Key.String()))

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 0.09 + 0.03*float64(r)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("alleyflash", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, !char.StatusIsActive(lockoutKey)
		},
	})

	return w, nil
}
