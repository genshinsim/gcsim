package candace

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Candace, NewChar)
}

type char struct {
	*tmpl.Character
	waveCount int
	infuseSrc int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.infuseSrc = -1

	w.Character = &c

	return nil
}

// Occurs after all characters are loaded, so getPartyElementalTypeCounts works properly
func (c *char) Init() error {
	c.burstSwap()
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}
