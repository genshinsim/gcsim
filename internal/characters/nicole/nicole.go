package nicole

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

type char struct {
	*tmpl.Character
	skillShield *shd
	skillBuff   []float64
	a1Buff      []float64
	projections int
	a1Src       int
	c2Buff      []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	hexerei, ok := p.Params["hexerei"]
	if !ok {
		hexerei = 1
	}
	c.IsHexerei = hexerei > 0

	return nil
}

func (c *char) Init() error {
	c.skillInit()
	c.burstInit()
	c.a1Init()
	c.a4Init()
	c.c1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	switch k {
	case info.AnimationXingqiuN0StartDelay:
		return 8
	case info.AnimationYelanN0StartDelay:
		return 5
	default:
		return c.Character.AnimationStartDelay(k)
	}
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "projections":
		if !c.StatusIsActive(burstKey) {
			return 0, nil
		}
		return c.projections, nil
	default:
		return c.Character.Condition(fields)
	}
}
