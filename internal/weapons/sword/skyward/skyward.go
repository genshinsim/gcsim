package skyward

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
	core.RegisterWeaponFunc(keys.SkywardBlade, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

const (
	buffKey = "skyward-blade"
)

// CRIT Rate increased by 4%. Gains Skypiercing Might upon using an Elemental
// Burst: Increases Movement SPD by 10%, increases ATK SPD by 10%, and Normal and
// Charged hits deal additional DMG equal to 20% of ATK. Skypiercing Might lasts
// for 12s.
func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	//perm buff
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.03 + float64(r)*0.01
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("skyward-blade-crit", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})

	atkspdBuff := make([]float64, attributes.EndStatType)
	atkspdBuff[attributes.AtkSpd] = 0.1
	c.Events.Subscribe(event.OnBurst, func(args ...interface{}) bool {
		if c.Player.Active() != char.Index {
			return false
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey, 720),
			AffectedStat: attributes.NoStat,
			Amount: func() ([]float64, bool) {
				return atkspdBuff, true
			},
		})
		return false
	}, fmt.Sprintf("skyward-blade-%v", char.Base.Key.String()))

	//deals damage proc on normal/charged attacks. i dont know why description in game sucks
	dmgper := .15 + .05*float64(r)
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		dmg := args[2].(float64)
		//check if char is correct?
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		//check if buff up
		if !char.StatModIsActive(buffKey) {
			return false
		}
		if dmg == 0 {
			return false
		}
		//add a new action that deals % dmg immediately
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Skyward Blade Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       dmgper,
		}
		trg := args[0].(combat.Target)
		c.QueueAttack(ai, combat.NewSingleTargetHit(trg.Key()), 0, 1)
		return false

	}, fmt.Sprintf("skyward-blade-%v", char.Base.Key.String()))

	return w, nil
}
