package polarstar

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
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

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Elemental Skill and Elemental Burst DMG increased by 12%. After a Normal Attack, Charged Attack,
	//Elemental Skill or Elemental Burst hits an opponent, 1 stack of Ashen Nightstar will be gained for 12s.
	//When 1/2/3/4 stacks of Ashen Nightstar are present, ATK is increased by 10/20/30/48%. The stack of Ashen
	//Nightstar created by the Normal Attack, Charged Attack, Elemental Skill or Elemental Burst will be counted
	//independently of the others.
	w := &Weapon{}
	r := p.Refine

	dmg := .09 + float64(r)*.03
	stack := .075 + float64(r)*.025
	max := .06 + float64(r)*.02

	normal := 0
	charged := 0
	skill := 0
	burst := 0

	mATK := make([]float64, attributes.EndStatType)
	char.AddStatMod(character.StatMod{Base: modifier.NewBase("polar-star", -1), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
		count := 0
		if normal > c.F {
			count++
		}
		if charged > c.F {
			count++
		}
		if skill > c.F {
			count++
		}
		if burst > c.F {
			count++
		}

		atkbonus := stack * float64(count)
		if count >= 4 {
			atkbonus += max
		}
		mATK[attributes.ATKP] = atkbonus

		return mATK, true
	}})

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.Player.Active() != char.Index {
			return false
		}

		cd := c.F + 60*12
		switch atk.Info.AttackTag {
		case combat.AttackTagNormal:
			normal = cd
		case combat.AttackTagExtra:
			charged = cd
		case combat.AttackTagElementalArt, combat.AttackTagElementalArtHold:
			skill = cd
		case combat.AttackTagElementalBurst:
			burst = cd
		}

		return false
	}, fmt.Sprintf("polar-star-%v", char.Base.Key.String()))

	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = dmg
	char.AddAttackMod(character.AttackMod{Base: modifier.NewBase("polar-star", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		switch atk.Info.AttackTag {
		case combat.AttackTagElementalArt, combat.AttackTagElementalArtHold, combat.AttackTagElementalBurst:
			return mDmg, true
		}
		return nil, false
	}})

	return w, nil
}
