package lisa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

func init() {
	core.RegisterCharFunc(keys.Lisa, NewChar)
}

type char struct {
	*tmpl.Character
	c6icd int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.Base.Element = attributes.Electro

	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = 4
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneMondstadt

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillHoldMult()
	if c.Base.Cons == 6 {
		c.c6()
	}

	return nil
}

func (c *char) skillHoldMult() {
	c.Core.Events.Subscribe(event.OnAttackWillLand, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.Abil != "Violet Arc (Hold)" {
			return false
		}
		stacks := t.GetTag(conductiveTag)

		atk.Info.Mult = skillHold[stacks][c.TalentLvlSkill()]

		//consume the stacks
		t.SetTag(conductiveTag, 0)

		return false
	}, "lisa-skill-hold-mul")
}
