package tartaglia

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

const (
	riptideFlashICDKey = "riptide-flash-icd"
	riptideKey         = "riptide"
	riptideSlashICDKey = "riptide-slash-icd"
	energyICDKey       = "riptide-energy-icd"
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
func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
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
	switch a {
	case action.ActionCharge:
		return 20
	}
	return c.Character.ActionStam(a, p)
}
