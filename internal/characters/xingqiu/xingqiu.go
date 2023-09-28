package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Xingqiu, NewChar)
}

type char struct {
	*tmpl.Character
	numSwords     int
	nextRegen     bool
	burstCounter  int
	orbitalActive bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	(&common.NAHook{
		C:           c.CharWrapper,
		AbilName:    "xingqiu burst",
		Core:        c.Core,
		AbilKey:     burstKey,
		AbilProcICD: 60,
		AbilICDKey:  burstICDKey,
		DelayFunc:   common.Get5PercentN0Delay,
		SummonFunc:  c.summonSwordWave,
	}).Enable()
	return nil
}
