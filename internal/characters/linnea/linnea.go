package linnea

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Linnea, NewChar)
}

type char struct {
	*tmpl.Character
	recastCount    int
	skillRecastSrc int
	skillSrc       int
	skillHitNum    int // skill Hit num is used to track the attack sequence when in Super Power form
	c1Stacks       int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.NormalCon = 3
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4Init()
	c.moonsignInit()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use E (in E recast)
	if a == action.ActionSkill && c.StatusIsActive(skillRecastKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 12
	case info.AnimationYelanN0StartDelay:
		return 5
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
