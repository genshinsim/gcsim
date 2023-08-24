package primordial

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.PrimordialJadeWingedSpear, NewWeapon)
}

type Weapon struct {
	Index  int
	stacks int
	buff   []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	//On hit, increases ATK by 3.2% for 6s. Max 7 stacks. This effect can only
	//occur once every 0.3s. While in possession of the maximum possible stacks,
	//DMG dealt is increased by 12%.
	w := &Weapon{}
	r := p.Refine
	const icdKey = "primordial-jade-spear-icd"
	const buffKey = "primordial"
	w.buff = make([]float64, attributes.EndStatType)
	perStackBuff := float64(r)*0.007 + 0.025
	dmgBuffAtMax := 0.09 + float64(r)*0.03

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		//check if char is correct?
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}
		//check if spear is on icd
		if char.StatusIsActive(icdKey) {
			return false
		}
		//check if buff expired; if so reset the stacks
		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		//every 0.3s
		char.AddStatus(icdKey, 18, true)

		if w.stacks < 7 {
			w.stacks++
			//check if it's max or amt
			if w.stacks == 7 {
				w.buff[attributes.DmgP] = dmgBuffAtMax
			}
			w.buff[attributes.ATKP] = float64(w.stacks) * perStackBuff
		}

		//refresh mod
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 6*60),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return w.buff, true
			},
		})

		return false
	}, fmt.Sprintf("primordial-%v", char.Base.Key.String()))
	return w, nil
}
