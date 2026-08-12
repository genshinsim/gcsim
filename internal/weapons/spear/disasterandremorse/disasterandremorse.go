package disasterandremorse

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	pathOfConflict = "path-of-conflict"
	unforgivable   = "unforgivable"
	irreparable    = "irreparable"
	procICD        = "disaster-remorse-proc"
	naCaICD        = "disaster-remorse-na-ca-icd"
	skillBurstICD  = "disaster-remorse-skill-burst-icd"
)

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine
	bonus := 0.40 + 0.10*float64(r-1)

	naCaVal := make([]float64, attributes.EndStatType)
	skillBurstVal := make([]float64, attributes.EndStatType)

	key := "disaster-remorse-" + char.Base.Key.String()

	unforgivableMod := character.AttackMod{
		Base: modifier.NewBase(unforgivable, 3*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagNormal, attacks.AttackTagExtra:
			default:
				return nil
			}

			if c.Player.GetHexereiCount() >= 2 {
				naCaVal[attributes.DmgP] = bonus * 1.75
			} else {
				naCaVal[attributes.DmgP] = bonus
			}

			return naCaVal
		},
	}

	irreparableMod := character.AttackMod{
		Base: modifier.NewBase(irreparable, 3*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt,
				attacks.AttackTagElementalArtHold,
				attacks.AttackTagElementalBurst:
			default:
				return nil
			}

			if c.Player.GetHexereiCount() >= 2 {
				skillBurstVal[attributes.DmgP] = bonus * 1.75
			} else {
				skillBurstVal[attributes.DmgP] = bonus
			}

			return skillBurstVal
		},
	}

	c.Events.Subscribe(event.OnSkill, func(args ...any) {
		if c.Player.Active() != char.Index() {
			return
		}

		if char.StatusIsActive(procICD) {
			return
		}

		char.AddStatus(procICD, 18*60, true)
		char.AddStatus(pathOfConflict, 17*60, true)

		char.AddAttackMod(unforgivableMod)
		char.AddAttackMod(irreparableMod)

		char.QueueCharTask(func() {
			char.DeleteAttackMod(unforgivable)
			char.DeleteAttackMod(irreparable)
		}, 17*60)
	}, key)

	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)

		if atk.Info.ActorIndex != char.Index() {
			return
		}

		if !char.StatusIsActive(pathOfConflict) {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal, attacks.AttackTagExtra:
			if char.StatusIsActive(naCaICD) {
				return
			}

			char.AddStatus(naCaICD, 0.1*60, true)
			char.ExtendStatus(irreparable, 60)

		case attacks.AttackTagElementalArt,
			attacks.AttackTagElementalArtHold,
			attacks.AttackTagElementalBurst:
			if char.StatusIsActive(skillBurstICD) {
				return
			}

			char.AddStatus(skillBurstICD, 0.1*60, true)
			char.ExtendStatus(unforgivable, 60)
		}
	}, key+"-extend")

	c.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		prev := args[0].(int)

		if prev != char.Index() {
			return
		}

		char.DeleteAttackMod(unforgivable)
		char.DeleteAttackMod(irreparable)
		char.DeleteStatus(pathOfConflict)
	}, key+"-swap")

	return w, nil
}
