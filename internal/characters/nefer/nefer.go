package nefer

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Nefer, NewChar)
}

const (
	shadowDanceKey = "nefer-shadow-dance"
)

type char struct {
	*tmpl.Character
	ascendantGleam bool
	veilstacks     int
	maxVeilStacks  int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	c.Moonsign = 1
	c.SetNumCharges(action.ActionSkill, 2)
	c.maxVeilStacks = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.ascendantGleam = c.Core.Player.GetMoonsignLevel() >= 2
	c.lunarbloomInit()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationYelanN0StartDelay:
		return 10
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) addVeilStacks(count int) {
	if count <= 0 {
		return
	}
	c.veilstacks = min(c.veilstacks+count, c.maxVeilStacks)
}

func (c *char) consumeVeilStacks() int {
	stacks := c.veilstacks
	c.veilstacks = 0
	return stacks
}
