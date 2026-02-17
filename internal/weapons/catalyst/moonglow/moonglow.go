package moonglow

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.EverlastingMoonglow, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	mheal := make([]float64, attributes.EndStatType)
	mheal[attributes.Heal] = 0.075 + float64(r)*0.025
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("moonglow-heal-bonus", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() []float64 {
			return mheal
		},
	})

	nabuff := 0.005 + float64(r)*0.005
	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}

		flatdmg := char.MaxHP() * nabuff
		atk.Info.FlatDmg += flatdmg

		c.Log.NewEvent("moonglow add damage", glog.LogPreDamageMod, char.Index()).
			Write("damage_added", flatdmg)
	}, fmt.Sprintf("moonglow-nabuff-%v", char.Base.Key.String()))

	const buffKey = "moonglow-postburst"
	buffDuration := 720 // 12s * 60
	const icdKey = "moonglow-energy-icd"
	icd := 6 // 0.1s * 60

	c.Events.Subscribe(event.OnBurst, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}
		char.AddStatus(buffKey, buffDuration, true)
	}, fmt.Sprintf("moonglow-onburst-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}
		if !char.StatusIsActive(buffKey) || char.StatusIsActive(icdKey) {
			return
		}

		char.AddEnergy("moonglow", 0.6)
		char.AddStatus(icdKey, icd, true)
	}, fmt.Sprintf("moonglow-energy-%v", char.Base.Key.String()))

	return w, nil
}
