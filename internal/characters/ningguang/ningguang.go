package ningguang

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

const skillParticleICDKey = "ningguang-particle-icd"

func init() {
	core.RegisterCharFunc(keys.Ningguang, NewChar)
}

type char struct {
	*tmpl.Character
	c2reset       int
	jadeCount     int
	lastScreen    int
	prevAttack    attackType
	skillSnapshot combat.Snapshot
}

type attackType int

const (
	attackTypeLeft attackType = iota
	attackTypeRight
	attackTypeTwirl
	endAttackType
)

var attackTypeNames = map[attackType]string{
	attackTypeLeft:  "Left",
	attackTypeRight: "Right",
	attackTypeTwirl: "Twirl",
}

func (t attackType) String() string {
	return attackTypeNames[t]
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 40
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	// Initialize at some very low value so these happen correctly at start of sim
	c.c2reset = -9999
	c.jadeCount = 0
	c.prevAttack = attackTypeLeft

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		// A1:
		// When Ningguang is in possession of Star Jades, her Charged Attack does not consume Stamina.
		if c.Base.Ascension >= 1 && c.jadeCount > 0 {
			return 0
		}
		return 50
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "jadeCount":
		return c.jadeCount, nil
	case "prevAttack":
		return int(c.prevAttack), nil
	default:
		return c.Character.Condition(fields)
	}
}
