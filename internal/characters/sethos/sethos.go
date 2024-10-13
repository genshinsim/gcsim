package sethos

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Sethos, NewChar)
}

type char struct {
	*tmpl.Character
	a4Count  int
	c2Stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}

	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.NormalCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillRefundHook()
	c.a4()
	c.c1()
	c.c2()
	c.c4()
	return nil
}
func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 0
	}
	return c.Character.AnimationStartDelay(k)
}
