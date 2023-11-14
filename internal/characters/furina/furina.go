package furina

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Furina, NewChar)
}

type Arkhe int

const (
	ousia Arkhe = iota
	pneuma
)

func (a Arkhe) String() string {
	switch a {
	case ousia:
		return "Ousia"
	case pneuma:
		return "Pneuma"
	}
	return "unknown"
}

type char struct {
	*tmpl.Character
	curFanfare          float64
	maxQFanfare         float64
	maxFanfare          float64
	burstBuff           []float64
	a4Buff              []float64
	a1HealsStopFrameMap []int
	a1HealsFlagMap      []bool
	lastSummonSrc       int
	arkhe               Arkhe
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 5
	c.BurstCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.arkhe = ousia

	c.a1()

	c.a4()
	c.a4Tick()

	c.burstInit()

	return nil
}