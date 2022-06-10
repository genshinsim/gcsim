package xiao

import (
	"github.com/genshinsim/gcsim/internal/frames"
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

const normalHitNum = 6

func init() {
	initCancelFrames()
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
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Anemo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
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

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][1], 26)
	attackFrames[0][action.ActionAttack] = 25

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[1][action.ActionAttack] = 22

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 38)
	attackFrames[2][action.ActionAttack] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 42)
	attackFrames[3][action.ActionAttack] = 39

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 30)
	attackFrames[4][action.ActionAttack] = 24

	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5][0], 79)
	attackFrames[5][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(45)
	chargeFrames[action.ActionSkill] = 38
	chargeFrames[action.ActionBurst] = 37
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 43

	// high_plunge -> x
	highPlungeFrames = frames.InitAbilSlice(66)
	highPlungeFrames[action.ActionAttack] = 61
	highPlungeFrames[action.ActionJump] = 65
	highPlungeFrames[action.ActionSwap] = 64

	// low_plunge -> x
	lowPlungeFrames = frames.InitAbilSlice(62)
	lowPlungeFrames[action.ActionAttack] = 60
	lowPlungeFrames[action.ActionSkill] = 59
	lowPlungeFrames[action.ActionDash] = 60
	lowPlungeFrames[action.ActionJump] = 61

	// skill -> x
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionAttack] = 24
	skillFrames[action.ActionSkill] = 24
	skillFrames[action.ActionBurst] = 24
	skillFrames[action.ActionDash] = 35
	skillFrames[action.ActionSwap] = 35

	// burst -> x
	burstFrames = frames.InitAbilSlice(82)
	burstFrames[action.ActionDash] = 57
	burstFrames[action.ActionJump] = 58
	burstFrames[action.ActionSwap] = 67
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
		c.Core.Log.NewEvent("a1 adding dmg %", glog.LogCharacterEvent, c.Index, "stacks", stacks, "final", ds.Stats[attributes.DmgP], "time since burst start", c.Core.F-c.qStarted)

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
		c.Core.Log.NewEvent("xiao burst damage bonus", glog.LogCharacterEvent, c.Index, "bonus", bonus, "final", ds.Stats[attributes.DmgP])
	}
	return ds
}
