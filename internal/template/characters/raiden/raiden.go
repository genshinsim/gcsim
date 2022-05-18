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

func init() {
	core.RegisterCharFunc(keys.Raiden, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 90
	}
	c.Energy = float64(e)
	c.EnergyMax = 90
	c.Weapon.Class = weapon.WeaponClassSpear
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 5
	c.CharZone = character.ZoneInazuma

	w.Character = &c

	return nil
}

func (c *char) Init() error {

	c.InitCancelFrames()
	c.eyeOnDamage()
	c.onBurstStackCount()
	c.onSwapClearBurst()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return nil
}

func (c *char) InitCancelFrames() {
	c.initNormalCancels()
	c.initBurstAttackCancels()

	frames.InitAbilSlice(&chargeFrames, 37) //n1, skill, burst all at 37
	chargeFrames[action.ActionSwap] = 36

	frames.InitAbilSlice(&swordCAFrames, 56)
	swordCAFrames[action.ActionDash] = 35
	swordCAFrames[action.ActionJump] = 35
	swordCAFrames[action.ActionSwap] = 55

	frames.InitAbilSlice(&skillFrames, 37)
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 17
	skillFrames[action.ActionSwap] = 17

	frames.InitAbilSlice(&burstFrames, 112)
	burstFrames[action.ActionAttack] = 111
	burstFrames[action.ActionSkill] = 111
	burstFrames[action.ActionCharge] = 500 //TODO: this action is illegal
	burstFrames[action.ActionDash] = 110
	burstFrames[action.ActionSwap] = 110
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionDash:
		return 18
	case action.ActionCharge:
		if c.Core.Status.Duration("raidenburst") == 0 {
			return 25
		}
		return 20
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", glog.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	s := c.Character.Snapshot(a)

	//a1 add dmg based on ER%
	excess := int(s.Stats[attributes.ER] / 0.01)

	s.Stats[attributes.ElectroP] += float64(excess) * 0.004 /// 0.4% extra dmg
	c.Core.Log.NewEvent("a4 adding electro dmg", glog.LogCharacterEvent, c.Index, "stacks", excess, "final", s.Stats[attributes.ElectroP])

	return s
}
