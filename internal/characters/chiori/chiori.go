package chiori

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Chiori, NewChar)
}

type char struct {
	*tmpl.Character

	// dolls
	skillDoll        *ticker
	constructChecker *ticker
	rockDoll         *ticker

	// a1 tracking
	a1Triggered   bool
	a1AttackCount int

	a4buff []float64

	// cons
	geoCount int
	kinus    []*ticker
	c2Ticker *ticker

	c4count int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1init()
	c.a4init()

	c.c2init()
	c.c4init()

	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if stiletto is on-field
	if a == action.ActionSkill && c.StatusIsActive(a1TailorMadeWindowKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}
