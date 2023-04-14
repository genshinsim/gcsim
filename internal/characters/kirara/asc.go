package kirara

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a1IcdStatus = "kirara-a1-icd"
)

// When Kirara is in the Urgent Neko Parcel state of Meow-teor Kick, each impact against an opponent will grant her a stack of Reinforced Packaging.
// This effect can be triggered once for each opponent hit every 0.5s. Max 3 stacks. When the Urgent Neko Parcel state ends, each stack of Reinforced
// Packaging will create 1 Shield of Safe Transport for Kirara. The shields that are created this way will have 20% of the DMG absorption that
// the Shield of Safe Transport produced by Meow-teor Kick would have. If Kirara is already protected by a Shield of Safe Transport created by
// Meow-teor Kick, its DMG absorption will stack with these shields and its duration will reset.
func (c *char) a1(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(a1IcdStatus) {
		return
	}
	if c.a1Stacks >= 3 {
		return
	}
	c.a1Stacks++

	shieldamt := c.shieldHP() * 0.2
	c.genShield("Shield of Safe Transport", shieldamt)
	c.AddStatus(a1IcdStatus, 0.5*60, true)
}

// Every 1,000 Max HP Kirara possesses will increase the DMG dealt by Meow-teor Kick by 0.4%, and the DMG dealt by Secret Art: Surprise Dispatch by 0.3%.
func (c *char) a4() {
	mSkill := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kirara-a4-skill", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalArt && atk.Info.AttackTag != attacks.AttackTagElementalArtHold {
				return nil, false
			}
			mSkill[attributes.DmgP] = c.MaxHP() * 0.001 * 0.004
			return mSkill, true
		},
	})

	mBurst := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("kirara-a4-burst", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.AttackTag != attacks.AttackTagElementalBurst {
				return nil, false
			}
			mBurst[attributes.DmgP] = c.MaxHP() * 0.001 * 0.003
			return mBurst, true
		},
	})
}
