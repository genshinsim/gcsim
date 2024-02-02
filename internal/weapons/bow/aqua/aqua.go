package aqua

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.AquaSimulacra, NewWeapon)
}

type Weapon struct {
	Index   int
	dmgBuff []float64
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	v := make([]float64, attributes.EndStatType)
	v[attributes.HPP] = 0.12 + float64(r)*0.04
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("aquasimulacra-hp", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return v, true
		},
	})

	w.dmgBuff = make([]float64, attributes.EndStatType)
	w.dmgBuff[attributes.DmgP] = 0.15 + float64(r)*0.05
	// queue up first tick of the dmg buff
	char.QueueCharTask(w.enemyCheck(char, c), 30)

	return w, nil
}

func (w *Weapon) enemyCheck(char *character.CharWrapper, c *core.Core) func() {
	return func() {
		enemies := c.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Combat.Player(), nil, 8), nil)
		if enemies != nil {
			char.AddAttackMod(character.AttackMod{
				Base: modifier.NewBaseWithHitlag("aquasimulacra-dmg", 72),
				Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
					return w.dmgBuff, true
				},
			})
		}
		char.QueueCharTask(w.enemyCheck(char, c), 30)
	}
}
