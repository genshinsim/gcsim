package xiao

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Xiao, NewChar)
}

// Xiao specific character implementation
type char struct {
	*tmpl.Character
	qStarted int
	a4stacks int
	a4buff   []float64
	c6Count  int
}

// Initializes character
// TODO: C4 is not implemented - don't really care about def
func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
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
	c.a4buff = make([]float64, attributes.EndStatType)
	c.onExitField()
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

// Xiao specific Snapshot implementation for his burst bonuses. Similar to Hu Tao
// Implements burst anemo attack damage conversion and DMG bonus
// Also implements A1:
// While under the effects of Bane of All Evil, all DMG dealt by Xiao is increased by 5%. DMG is increased by an additional 5% for every 3s the ability persists. The maximum DMG Bonus is 25%
func (c *char) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(a)

	if c.StatusIsActive("xiaoburst") {
		// Anemo conversion and dmg bonus application to normal, charged, and plunge attacks
		// Also handle burst CA ICD change to share with Normal
		switch a.AttackTag {
		case attacks.AttackTagNormal:
			// QN1-1 has different hitlag from N1-1
			if a.Abil == "Normal 0" {
				// this also overwrites N1-2 HitlagHaltFrames but they have the same value so it's fine
				a.HitlagHaltFrames = 0.01 * 60
			}
		case attacks.AttackTagExtra:
			// Q-CA has different hitlag from CA
			a.ICDTag = attacks.ICDTagNormalAttack
			a.HitlagHaltFrames = 0.04 * 60
		case attacks.AttackTagPlunge:
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
