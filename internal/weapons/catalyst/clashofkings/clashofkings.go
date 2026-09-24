package clashofkings

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	chargeExtensionKey = "clash-of-kings-ca-extension"
	buffKey            = "clash-of-kings"
	buffICDKey         = "clash-of-kings-icd"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15 + float64(r)*0.05
	m[attributes.EM] = 75 + float64(r)*25
	onSkill := func(args ...any) {
		if char.Index() == c.Player.Active() {
			return
		}

		if char.StatusIsActive(buffICDKey) {
			return
		}

		char.AddStatus(buffICDKey, 12*60, true)
		char.AddStatus(chargeExtensionKey, 12*60, true)
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(buffKey, 6*60),
			Amount: func() []float64 {
				return m
			},
		})
	}

	c.Events.Subscribe(event.OnSkill, onSkill, "clash-of-kings-on-skill-"+char.Base.Key.String())

	onHit := func(args ...any) {
		atk, ok := args[1].(*info.AttackEvent)

		if !ok {
			return
		}

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if char.Index() != c.Player.Active() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return
		}

		if !char.StatusIsActive(buffKey) {
			return
		}

		if !char.StatusIsActive(chargeExtensionKey) {
			return
		}

		char.DeleteStatus(chargeExtensionKey)
		char.ExtendStatus(buffKey, 6*60)
	}

	c.Events.Subscribe(event.OnEnemyHit, onHit, "clash-of-kings-on-charge-hit-"+char.Base.Key.String())

	return w, nil
}
