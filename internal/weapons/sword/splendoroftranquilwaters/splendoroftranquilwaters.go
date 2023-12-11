package splendoroftranquilwaters

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	skillBuffKey = "splendoroftranquilwaters-skill-buff"
	skillBuffIcd = "splendoroftranquilwaters-skill-buff-icd"
	hpBuffKey    = "splendoroftranquilwaters-hp-buff"
	hpBuffIcd    = "splendoroftranquilwaters-hp-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.SplendorOfTranquilWaters, NewWeapon)
}

type Weapon struct {
	skillStacks int
	hpStacks    int
	core        *core.Core
	char        *character.CharWrapper
	refine      int
	buffSkill   []float64
	buffHp      []float64
	Index       int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{
		core:      c,
		char:      char,
		refine:    p.Refine,
		buffSkill: make([]float64, attributes.EndStatType),
		buffHp:    make([]float64, attributes.EndStatType),
	}

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.ActorIndex != char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		w.onEquipChangeHP()
		return false
	}, fmt.Sprintf("splendoroftranquilwaters-equip-drain-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if index != char.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}
		w.onEquipChangeHP()
		return false
	}, fmt.Sprintf("splendoroftranquilwaters-equip-heal-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.ActorIndex == char.Index {
			return false
		}
		if di.Amount <= 0 {
			return false
		}
		w.onOtherChangeHP()
		return false
	}, fmt.Sprintf("splendoroftranquilwaters-other-drain-%v", char.Base.Key.String()))

	c.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		index := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)
		if index == char.Index {
			return false
		}
		if amount <= 0 {
			return false
		}
		if math.Abs(amount-overheal) <= 1e-9 {
			return false
		}
		w.onOtherChangeHP()
		return false
	}, fmt.Sprintf("splendoroftranquilwaters-other-heal-%v", char.Base.Key.String()))

	return w, nil
}

func (w *Weapon) onEquipChangeHP() {
	if w.char.StatusIsActive(skillBuffIcd) {
		return
	}
	if !w.char.StatModIsActive(skillBuffKey) {
		w.skillStacks = 0
	}
	if w.skillStacks < 3 {
		w.skillStacks++
	}

	w.char.AddStatus(skillBuffIcd, 0.2*60, true)
	w.buffSkill[attributes.DmgP] = (0.06 + 0.02*float64(w.refine)) * float64(w.skillStacks)
	w.char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(skillBuffKey, 6*60),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
				return w.buffSkill, true
			case attacks.AttackTagElementalArtHold:
				return w.buffSkill, true
			default:
				return nil, false
			}
		},
	})
}

func (w *Weapon) onOtherChangeHP() {
	if w.char.StatusIsActive(hpBuffIcd) {
		return
	}
	if !w.char.StatModIsActive(hpBuffKey) {
		w.hpStacks = 0
	}
	if w.hpStacks < 2 {
		w.hpStacks++
	}

	hpp := (0.105 + float64(w.refine)*0.035) * float64(w.hpStacks)
	val := make([]float64, attributes.EndStatType)
	val[attributes.HPP] = hpp
	w.char.AddStatus(hpBuffIcd, 0.2*60, true)
	w.char.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(hpBuffKey, 6*60),
		AffectedStat: attributes.HPP,
		Amount: func() ([]float64, bool) {
			return val, true
		},
	})
}
