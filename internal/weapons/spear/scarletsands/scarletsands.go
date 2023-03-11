package scarletsands

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.StaffOfTheScarletSands, NewWeapon)
}

type Weapon struct {
	stacks int
	Index  int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const skillBuff = "scarletsands-skill"

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	mATK := make([]float64, attributes.EndStatType)
	atkBuff := 0.39 + 0.13*float64(r)
	atkSkillBuff := 0.21 + 0.07*float64(r)
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("scarletsands", -1),
		AffectedStat: attributes.ATK,
		Extra:        true,
		Amount: func() ([]float64, bool) {
			em := char.NonExtraStat(attributes.EM)
			mATK[attributes.ATK] = atkBuff * em
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
		if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
			return false
		}
		// TODO: is there icd?

		// reset stacks if expired
		if !char.StatModIsActive(skillBuff) {
			w.stacks = 0
		}
		if w.stacks < 3 {
			w.stacks++
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(skillBuff, 10*60),
			AffectedStat: attributes.ATK,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				em := char.NonExtraStat(attributes.EM)
				mATK[attributes.ATK] = atkSkillBuff * em * float64(w.stacks)
				return mATK, true
			},
		})

		c.Log.NewEvent("scarletsands adding stack", glog.LogWeaponEvent, char.Index).Write("stacks", w.stacks)
		return false
	}, fmt.Sprintf("scarletsands-%v", char.Base.Key.String()))

	return w, nil
}
