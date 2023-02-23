package keyofkhajnisut

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.KeyOfKhajNisut, NewWeapon)
}

// HP increased by 20%. When an Elemental Skill hits opponents, you gain the Grand Hymn effect for 20s.
// This effect increases the equipping character’s Elemental Mastery by 0.12% of their Max HP. This effect can trigger once every 0.3s.
// Max 3 stacks. When this effect gains 3 stacks, or when the third stack’s duration is refreshed, the Elemental Mastery
// of all nearby party members will be increased by 0.2% of the equipping character’s max HP for 20s.
type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey     = "khaj-nisut-buff"
	teamBuffKey = "khaj-nisut-team-buff"
	icdKey      = "khaj-nisut-icd"
)

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	duration := 1200 // 20*60
	cd := 18         // 0.3 * 60

	hp := 0.15 + 0.05*float64(r)
	em := 0.0009 + 0.0003*float64(r)
	emTeam := 0.0015 + 0.0005*float64(r)

	m := make([]float64, attributes.EndStatType)
	m[attributes.HPP] = hp
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("khaj-nisut", -1),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Player.Active() != char.Index {
			return false
		}
		if char.StatusIsActive(icdKey) {
			return false
		}
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}

		if !char.StatModIsActive(buffKey) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}

		val := make([]float64, attributes.EndStatType)
		val[attributes.EM] = char.MaxHP() * em * float64(w.stacks)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, duration),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})

		if w.stacks == 3 {
			val := make([]float64, attributes.EndStatType)
			val[attributes.EM] = char.MaxHP() * emTeam
			for _, this := range c.Player.Chars() {
				this.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(teamBuffKey, duration),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return val, true
					},
				})
			}
		}

		char.AddStatus(icdKey, cd, true)
		return false
	}, "khaj-nisut")

	return w, nil
}
