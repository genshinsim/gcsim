package tartaglia

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	normalHitNum       = 6
	riptideDuration    = 18 * 60
	riptideFlashICDKey = "riptide-flash-icd"
	riptideKey         = "riptide"
	riptideSlashICDKey = "riptide-slash-icd"
)

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Tartaglia, NewChar)
}

// tartaglia specific character implementation
type char struct {
	*tmpl.Character
	eCast         int // the frame tartaglia casts E to enter melee stance
	rtParticleICD int
	mlBurstUsed   bool // used for c6
}

// Initializes character
func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassBow
	c.SkillCon = 3
	c.BurstCon = 5
	c.NormalHitNum = normalHitNum

	c.eCast = 0
	c.rtParticleICD = 0
	c.mlBurstUsed = false

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.onDefeatTargets()

	for _, char := range c.Core.Player.Chars() {
		char.SetTag(keys.ChildePassive, 1)
	}

	return nil
}

func initCancelFrames() {
	initRangedFrames()
	initMeleeFrames()
}

func initRangedFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 17)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 13)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 34)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 37)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 22)
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5], 39)

	// aimed -> x
	aimedFrames = frames.InitAbilSlice(84)

	// skill -> x
	skillRangedFrames = frames.InitAbilSlice(28)

	// burst -> x
	burstRangedFrames = frames.InitAbilSlice(52)
}

func initMeleeFrames() {
	// melee cancels

	// NA cancels (melee)
	meleeFrames = make([][]int, normalHitNum)

	meleeFrames[0] = frames.InitNormalCancelSlice(meleeFrames[0][0], 7)
	meleeFrames[1] = frames.InitNormalCancelSlice(meleeFrames[1][0], 13)
	meleeFrames[2] = frames.InitNormalCancelSlice(meleeFrames[2][0], 28)
	meleeFrames[3] = frames.InitNormalCancelSlice(meleeFrames[3][0], 32)
	meleeFrames[4] = frames.InitNormalCancelSlice(meleeFrames[4][0], 36)
	meleeFrames[5] = frames.InitNormalCancelSlice(meleeFrames[5][1], 49)

	// charge -> x
	chargeFrames = frames.InitAbilSlice(73)
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]

	// skill -> x
	skillMeleeFrames = frames.InitAbilSlice(20)

	// burst -> x
	burstMeleeFrames = frames.InitAbilSlice(97)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		return 20
	}
	return c.Character.ActionStam(a, p)
}
