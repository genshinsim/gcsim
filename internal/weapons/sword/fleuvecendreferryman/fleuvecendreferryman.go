package fleuvecendreferryman

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.FleuveCendreFerryman, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Increases Elemental Skill CRIT Rate by 8/10/12/14/16%. Additionally, increases Energy Recharge by 16/20/24/28/32% for 5s after using an Elemental Skill.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	// skill crit rate
	mCR := make([]float64, attributes.EndStatType)
	mCR[attributes.CR] = 0.06 + 0.02*float64(r)
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("fleuvecendreferryman-cr", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag == attacks.AttackTagElementalArt || atk.Info.AttackTag == attacks.AttackTagElementalArtHold {
				return mCR, true
			}
			return nil, false
		},
	})

	// er up after skill
	mER := make([]float64, attributes.EndStatType)
	mER[attributes.ER] = 0.12 + 0.04*float64(r)
	c.Events.Subscribe(event.OnSkill, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("fleuvecendreferryman-er", 5*60),
			AffectedStat: attributes.ER,
			Amount: func() ([]float64, bool) {
				return mER, true
			},
		})
		return false
	}, fmt.Sprintf("fleuvecendreferryman-%v", char.Base.Key.String()))

	return w, nil
}
