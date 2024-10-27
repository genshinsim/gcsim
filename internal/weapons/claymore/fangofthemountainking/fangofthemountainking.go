package fangofthemountainking

import (
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

var stacks = []string{
	canopyFavor1,
	canopyFavor2,
	canopyFavor3,
	canopyFavor4,
	canopyFavor5,
	canopyFavor6,
}

const (
	canopyFavor1          = "canopy-favor-1"
	canopyFavor2          = "canopy-favor-2"
	canopyFavor3          = "canopy-favor-3"
	canopyFavor4          = "canopy-favor-4"
	canopyFavor5          = "canopy-favor-5"
	canopyFavor6          = "canopy-favor-6"
	maxCanopyFavorStacks  = 6
	burnBurgStackGain     = 3
	skillHitStackGain     = 1
	triggerSkillIcdKey    = "fotmk-elem-art-icd"
	triggerBurnBurgIcdKey = "fotmk-burn-burg-icd"
)

func init() {
	core.RegisterWeaponFunc(keys.FangOfTheMountainKing, NewWeapon)
}

type Weapon struct {
	Index            int
	char             *character.CharWrapper
	lastUpdatedStack int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

// Gain 1 stack of Canopy's Favor after hitting an opponent with an Elemental Skill. This can be triggered once every 0.5s.
// After a nearby party member triggers a Burning or Burgeon reaction, the equipping character will gain 3 stacks.
// This effect can be triggered once every 2s and can be triggered even when the triggering party member is off-field.
// Canopy's Favor: Elemental Skill and Burst DMG is increased by 10% for 6s. Max 6 stacks. Each stack is counted independently.
func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	w.char = char
	w.lastUpdatedStack = 0

	amt := 0.10 + float64(r-1)*0.025
	m := make([]float64, attributes.EndStatType)

	c.Events.Subscribe(event.OnBurning, w.getHook(triggerBurnBurgIcdKey, burnBurgStackGain), "fang-of-the-mountain-king")
	c.Events.Subscribe(event.OnBurgeon, w.getHook(triggerBurnBurgIcdKey, burnBurgStackGain), "fang-of-the-mountain-king")
	c.Events.Subscribe(event.OnEnemyDamage, w.getHook(triggerSkillIcdKey, skillHitStackGain), "fang-of-the-mountain-king")

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("fang-of-the-mountain-king", -1),
		Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			switch a.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return nil, false
			}
			m[attributes.DmgP] = amt * float64(w.getStacksNum())
			return m, true
		},
	})

	return w, nil
}

func (w *Weapon) getStacksNum() int {
	stacksNum := 0
	for _, stack := range stacks {
		if w.char.StatusIsActive(stack) {
			stacksNum++
		}
	}
	return stacksNum
}

func (w *Weapon) getHook(icdKey string, stacksToAdd int) func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		if stacksToAdd == 1 {
			atk := args[1].(*combat.AttackEvent)
			if atk.Info.ActorIndex != w.char.Index {
				return false
			}
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return false
			}
		}
		if w.char.StatusIsActive(icdKey) {
			return false
		}
		for i := w.lastUpdatedStack; i < w.lastUpdatedStack+stacksToAdd; i++ {
			w.char.AddStatus(stacks[i%6], 6*60, true)
		}
		w.lastUpdatedStack += stacksToAdd
		w.lastUpdatedStack %= 6
		return false
	}
}
