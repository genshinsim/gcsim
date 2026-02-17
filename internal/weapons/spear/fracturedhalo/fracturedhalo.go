package fracturedhalo

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	AtkpBuff = "fracturedhalo-atkp"
	LcBuff   = "fracturedhalo-lc-buff"
)

func init() {
	core.RegisterWeaponFunc(keys.FracturedHalo, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// After an Elemental Skill or Elemental Burst is used, ATK is increased by 24%/30%/36%/42%/48% for 20s.
// If the equipping character creates a Shield while this effect is active, they will
// gain the Electrifying Edict effect for 20s: All nearby party members
// deal 40%/50%/60%/70%/80% more Lunar-Charged DMG.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	lcBonus := 0.30 + float64(r)*0.10
	mAtk := make([]float64, attributes.EndStatType)
	mAtk[attributes.ATKP] = 0.18 + float64(r)*0.06

	atkBuffHook := func(args ...any) bool {
		if c.Player.Active() != char.Index() {
			return false
		}

		char.AddStatMod(character.StatMod{
			Base: modifier.NewBaseWithHitlag(AtkpBuff, 20*60),
			Amount: func() ([]float64, bool) {
				return mAtk, true
			},
		})

		return false
	}

	c.Events.Subscribe(event.OnBurst, atkBuffHook, fmt.Sprintf("fracturedhalo-burst-%v", char.Base.Key.String()))
	c.Events.Subscribe(event.OnSkill, atkBuffHook, fmt.Sprintf("fracturedhalo-skill-%v", char.Base.Key.String()))

	// TODO: How does buff from a different refine fractured halo work?
	for _, ch := range c.Player.Chars() {
		ch.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("fracturedhalo-lc-dmg", -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				switch ai.AttackTag {
				case attacks.AttackTagDirectLunarCharged:
				case attacks.AttackTagReactionLunarCharge:
				default:
					return 0, false
				}
				if !char.StatusIsActive(LcBuff) {
					return 0, false
				}
				return lcBonus, false
			},
		})
	}

	c.Events.Subscribe(event.OnShielded, func(args ...any) bool {
		shd := args[0].(shield.Shield)
		if shd.ShieldOwner() != char.Index() {
			return false
		}
		if !char.StatModIsActive(AtkpBuff) {
			return false
		}
		char.AddStatus(LcBuff, 20*60, true)
		return false
	}, fmt.Sprintf("fracturedhalo-shield-%v", char.Base.Key.String()))

	return w, nil
}
