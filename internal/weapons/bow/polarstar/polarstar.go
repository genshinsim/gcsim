package polarstar

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
	core.RegisterWeaponFunc(keys.PolarStar, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	// Elemental Skill and Elemental Burst DMG increased by 12%. After a Normal Attack, Charged Attack,
	// Elemental Skill or Elemental Burst hits an opponent, 1 stack of Ashen Nightstar will be gained for 12s.
	// When 1/2/3/4 stacks of Ashen Nightstar are present, ATK is increased by 10/20/30/48%. The stack of Ashen
	// Nightstar created by the Normal Attack, Charged Attack, Elemental Skill or Elemental Burst will be counted
	// independently of the others.
	w := &Weapon{}
	r := p.Refine

	dmg := .09 + float64(r)*.03
	stack := .075 + float64(r)*.025
	max := .06 + float64(r)*.02

	const normalKey = "polar-star-normal"
	const chargedKey = "polar-star-charged"
	const skillKey = "polar-star-skill"
	const burstKey = "polar-star-burst"
	stackDuration := 720 // 12s * 60

	mATK := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("polar-star-atk", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			count := 0
			if char.StatusIsActive(normalKey) {
				count++
			}
			if char.StatusIsActive(chargedKey) {
				count++
			}
			if char.StatusIsActive(skillKey) {
				count++
			}
			if char.StatusIsActive(burstKey) {
				count++
			}

			atkbonus := stack * float64(count)
			if count >= 4 {
				atkbonus += max
			}
			mATK[attributes.ATKP] = atkbonus

			return mATK, true
		},
	})

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			char.AddStatus(normalKey, stackDuration, true)
		case attacks.AttackTagExtra:
			char.AddStatus(chargedKey, stackDuration, true)
		case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold:
			char.AddStatus(skillKey, stackDuration, true)
		case attacks.AttackTagElementalBurst:
			char.AddStatus(burstKey, stackDuration, true)
		}

		return false
	}, fmt.Sprintf("polar-star-%v", char.Base.Key.String()))

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = dmg
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("polar-star-dmg", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold, attacks.AttackTagElementalBurst:
				return mDmg, true
			}
			return nil, false
		},
	})

	return w, nil
}
