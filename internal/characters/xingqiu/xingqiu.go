package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
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
	naHook        *common.NAHook
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
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
	c.naHook = common.NewNAHook(c.CharWrapper, c.Core, "xingqiu burst", burstKey, 60, burstICDKey, common.Get5PercentN0Delay, c.summonSwordWave)
	c.naHook.NAStateHook()
	return nil
}
