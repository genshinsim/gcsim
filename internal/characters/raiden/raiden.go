package raiden

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
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

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 90
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.eyeOnDamage()
	c.a1()
	c.onBurstStackCount()
	c.onSwapClearBurst()
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.StatusIsActive(BurstKey) {
			return 20
		}
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Snapshot(a *combat.AttackInfo) combat.Snapshot {
	s := c.Character.Snapshot(a)

	// A4:
	// Each 1% above 100% Energy Recharge that the Raiden Shogun possesses grants her:
	//
	// - 0.4% Electro DMG Bonus.
	if c.Base.Ascension >= 4 {
		excess := int(s.Stats[attributes.ER] / 0.01)
		s.Stats[attributes.ElectroP] += float64(excess) * 0.004
		c.Core.Log.NewEvent("a4 adding electro dmg", glog.LogCharacterEvent, c.Index).
			Write("stacks", excess).
			Write("final", s.Stats[attributes.ElectroP])
	}

	return s
}
