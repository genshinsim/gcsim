package sethos

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Sethos, NewChar)
}

type char struct {
	*tmpl.Character
	lastSkillFrame int
	a4Count        int
	c4Buff         []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}

	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.NormalCon = 3

	c.lastSkillFrame = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillRefundHook()
	c.a4()
	c.c1()
	c.c2()
	c.c4()
	c.onExitField()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 5
	case info.AnimationYelanN0StartDelay:
		return 4
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...any) bool {
		if c.StatusIsActive(burstBuffKey) {
			c.DeleteStatus(burstBuffKey)
		}
		return false
	}, "sethos-exit")
}
