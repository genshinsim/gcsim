package spine

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SerpentSpine, NewWeapon)
}

type Weapon struct {
	Index  int
	char   *character.CharWrapper
	c      *core.Core
	stacks int
	dmg    float64
	buff   []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }
func (w *Weapon) stackCheck() func() {
	return func() {
		//if on field and stack < 5, add a stack
		if w.char.Index == w.c.Player.Active() {
			if w.stacks < 5 {
				w.stacks++
				w.updateBuff()
			}
		}
		w.char.QueueCharTask(w.stackCheck(), 240) //check again in 4s
	}
}
func (w *Weapon) updateBuff() {
	w.buff[attributes.DmgP] = float64(w.stacks) * w.dmg
}

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Every 4s a character is on the field, they will deal 6% more DMG and take
	//3% more DMG. This effect has a maximum of 5 stacks and will not be reset
	//if the character leaves the field, but will be reduced by 1 stack when the
	//character takes DMG.
	w := &Weapon{
		char: char,
		c:    c,
		buff: make([]float64, attributes.EndStatType),
	}
	r := p.Refine

	//the damage taken/stack reduciton has a 1s internal cooldown
	//otherwise every 4s it does a check and adds a stack if on field it looks like
	//verified this in game by checking first sword added, swap off and on, and note that the
	//next sword animation is 4s despite having been off field for ~1s

	w.dmg = 0.05 + float64(r)*.01
	//set initial
	w.stacks = p.Params["stacks"]
	c.Log.NewEvent(
		"serpent spine stack check", glog.LogWeaponEvent, char.Index,
	).
		Write("params", p.Params)

	if w.stacks > 5 {
		w.stacks = 5
	}
	w.updateBuff()

	//start ticker to check for stack increase
	char.QueueCharTask(w.stackCheck(), 240)

	//add event hook to check for dmg, subject to 1s icd
	//TODO: taking 3% more damage not implemented
	const icdKey = "spine-dmgtaken-icd"
	icd := 60
	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if !di.External {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		char.AddStatus(icdKey, icd, true)
		if w.stacks > 0 {
			w.stacks--
			w.updateBuff()
		}
		return false
	}, fmt.Sprintf("spine-%v", char.Base.Key.String()))

	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("spine", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return w.buff, w.stacks > 0
		},
	})

	return w, nil
}
