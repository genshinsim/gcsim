package tartaglia

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	riptideFlashICDKey = "riptide-flash-icd"
	riptideKey         = "riptide"
	riptideSlashICDKey = "riptide-slash-icd"
	particleICDKey     = "tartaglia-particle-icd"
	MeleeKey           = "tartagliamelee"
)

func init() {
	core.RegisterCharFunc(keys.Tartaglia, NewChar)
}

// tartaglia specific character implementation
type char struct {
	*tmpl.Character
	riptideDuration int
	eCast           int  // the frame tartaglia casts E to enter melee stance
	c4Src           int  // used for c4
	mlBurstUsed     bool // used for c6
}

// Initializes character
func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = normalHitNum

	c.riptideDuration = 10 * 60

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	c.onExitField()
	c.onDefeatTargets()

	for _, char := range c.Core.Player.Chars() {
		char.SetTag(keys.ChildePassive, 1)
	}

	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge {
		return 20
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) NextQueueItemIsValid(a action.Action, p map[string]int) error {
	if a == action.ActionCharge && c.Core.Player.LastAction.Type != action.ActionAttack {
		return player.ErrInvalidChargeAction
	}
	return c.NextQueueItemIsValid(a, p)
}
