package raiden

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

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Raiden, NewChar)
}

type char struct {
	*tmpl.Character
	burstCastF     int
	eyeICD         int
	stacksConsumed float64
	stacks         float64
	restoreICD     int
	restoreCount   int
	applyC4        bool
	c6Count        int
	c6ICD          int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro

	c.EnergyMax = 90
	c.Weapon.Class = weapon.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneInazuma

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.eyeOnDamage()
	c.onBurstStackCount()
	c.onSwapClearBurst()
	return nil
}

func initCancelFrames() {
	initAttackFrames()
	initSwordFrames()

	// charge -> x
	chargeFrames = frames.InitAbilSlice(37) //n1, skill, burst all at 37
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 36

	// charge (burst) -> x
	swordCAFrames = frames.InitAbilSlice(56)
	swordCAFrames[action.ActionDash] = swordCAHitmarks[len(swordCAHitmarks)-1]
	swordCAFrames[action.ActionJump] = swordCAHitmarks[len(swordCAHitmarks)-1]

	// skill -> x
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 17
	skillFrames[action.ActionSwap] = 36

	// burst -> x
	burstFrames = frames.InitAbilSlice(112)
	burstFrames[action.ActionAttack] = 111
	burstFrames[action.ActionCharge] = 500 //TODO: this action is illegal
	burstFrames[action.ActionSkill] = 111
	burstFrames[action.ActionDash] = 110
	burstFrames[action.ActionSwap] = 110
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.Core.Status.Duration("raidenburst") > 0 {
			return 20
		}
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	s := c.Character.Snapshot(a)

	//a1 add dmg based on ER%
	excess := int(s.Stats[attributes.ER] / 0.01)

	s.Stats[attributes.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Core.Log.NewEvent("a4 adding electro dmg", glog.LogCharacterEvent, c.Index, "stacks", excess, "final", s.Stats[attributes.ElectroP])

	return s
}
