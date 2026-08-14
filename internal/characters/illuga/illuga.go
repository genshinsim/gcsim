package illuga

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character

	nightingalesSong               int
	nightingalesSongExtraConstruct int
	c2QuillCounter                 int
	c4Src                          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = 4
	c.SkillCon = 5
	c.BurstCon = 3
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstBuffInit()
	c.c1Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 11
	}
	return c.Character.AnimationStartDelay(k)
}
