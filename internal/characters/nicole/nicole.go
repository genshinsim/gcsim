package nicole

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Nicole, NewChar)
}

type char struct {
	*tmpl.Character
	skillShield *shd
	skillBuff   []float64
	a1Buff      []float64
	burstHits   int
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
