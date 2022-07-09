package xiao

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterCharFunc(keys.Xiao, NewChar)
}

// Xiao specific character implementation
type char struct {
	*tmpl.Character
	qStarted int
	a4Expiry int
	c6Src    int
	c6Count  int
}

// Initializes character
// TODO: C4 is not implemented - don't really care about def
func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 70
	c.Weapon.Class = weapon.WeaponClassSpear
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	c.c6Count = 0

	c.SetNumCharges(action.ActionSkill, 2)
	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 3)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

// Xiao specific Snapshot implementation for his burst bonuses. Similar to Hu Tao
// Implements burst anemo attack damage conversion and DMG bonus
// Also implements A1:
// While under the effects of Bane of All Evil, all DMG dealt by Xiao is increased by 5%. DMG is increased by an additional 5% for every 3s the ability persists. The maximum DMG Bonus is 25%
func (c *char) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(a)

	if c.Core.Status.Duration("xiaoburst") > 0 {
		// Calculate and add A1 damage bonus - applies to all damage
		// Fraction dropped in int conversion in go - acts like floor
		stacks := 1 + int((c.Core.F-c.qStarted)/180)
		if stacks > 5 {
			stacks = 5
		}
		ds.Stats[attributes.DmgP] += float64(stacks) * 0.05
		c.Core.Log.NewEvent("a1 adding dmg %", glog.LogCharacterEvent, c.Index).
			Write("stacks", stacks).
			Write("final", ds.Stats[attributes.DmgP]).
			Write("time since burst start", c.Core.F-c.qStarted)

		// Anemo conversion and dmg bonus application to normal, charged, and plunge attacks
		// Also handle burst CA ICD change to share with Normal
		switch a.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagExtra:
			a.ICDTag = combat.ICDTagNormalAttack
		case combat.AttackTagPlunge:
		default:
			return ds
		}
		a.Element = attributes.Anemo
		bonus := burstBonus[c.TalentLvlBurst()]
		ds.Stats[attributes.DmgP] += bonus
		c.Core.Log.NewEvent("xiao burst damage bonus", glog.LogCharacterEvent, c.Index).
			Write("bonus", bonus).
			Write("final", ds.Stats[attributes.DmgP])
	}
	return ds
}
