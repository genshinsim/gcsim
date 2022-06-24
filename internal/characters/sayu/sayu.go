package sayu

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Sayu, NewChar)
}

type char struct {
	*tmpl.Character
	eInfused            attributes.Element
	eInfusedTag         combat.ICDTag
	eDuration           int
	infuseCheckLocation combat.AttackPattern
	c2Bonus             float64
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneInazuma

	c.eInfused = attributes.NoElement
	c.eDuration = -1
	c.c2Bonus = .0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.rollAbsorb()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func (c *char) rollAbsorb() {
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalArt && atk.Info.AttackTag != combat.AttackTagElementalArtHold {
			return false
		}
		if atk.Info.Element != attributes.Anemo || c.eInfused == attributes.NoElement {
			return false
		}
		if c.Core.F > c.eDuration {
			return false
		}

		switch atk.Info.AttackTag {
		case combat.AttackTagElementalArt:
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  combat.AttackTagElementalArtHold,
				ICDTag:     c.eInfusedTag,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    c.eInfused,
				Durability: 25,
				Mult:       skillAbsorb[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 1, 1)
		case combat.AttackTagElementalArtHold:
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Yoohoo Art: Fuuin Dash (Elemental DMG)",
				AttackTag:  combat.AttackTagElementalArt,
				ICDTag:     combat.ICDTagNone,
				ICDGroup:   combat.ICDGroupDefault,
				Element:    c.eInfused,
				Durability: 25,
				Mult:       skillAbsorbEnd[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, combat.NewDefCircHit(0.1, false, combat.TargettableEnemy), 1, 1)
		}

		return false
	}, "sayu-absorb-check")
}
