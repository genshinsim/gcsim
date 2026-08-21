package kachina

import (
	"github.com/genshinsim/gcsim/internal/template/character"
	nightsoultemplate "github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	playercharacter "github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	twirlyKey = "kachina-turbo-twirly"
	fieldKey  = "kachina-turbo-drill-field"
)

type char struct {
	*character.Character
	nightsoulState *nightsoultemplate.State
	twirlySrc      int
	mounted        bool
	fieldCenter    info.Point
	fieldChar      int
	c6ShieldKeys   map[string]bool
}

func NewChar(s *core.Core, w *playercharacter.CharWrapper, _ info.CharacterProfile) error {
	c := char{
		Character:    character.NewWithWrapper(s, w),
		twirlySrc:    -1,
		fieldChar:    -1,
		c6ShieldKeys: make(map[string]bool),
	}
	c.nightsoulState = nightsoultemplate.New(c.Core, c.CharWrapper)
	c.nightsoulState.MaxPoints = 60
	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5
	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.initA1()
	c.initC1()
	c.initC4()
	c.initC6()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	if len(fields) > 0 && fields[0] == "nightsoul" {
		return c.nightsoulState.Condition(fields)
	}
	if len(fields) > 0 && fields[0] == "mounted" {
		return c.mounted, nil
	}
	return c.Character.Condition(fields)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSkill && c.StatusIsActive(twirlyKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if c.mounted {
		return 0
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) startTwirly(hold bool) {
	c.startTwirlyWithPoints(nightsoul, hold)
}

func (c *char) startTwirlyWithPoints(points float64, hold bool) {
	c.nightsoulState.EnterBlessing(points)
	c.AddStatus(twirlyKey, -1, true)
	c.mounted = hold
	c.twirlySrc = c.Core.F
	if !hold {
		c.queueIndependent(c.twirlySrc, 143)
	}
}

func (c *char) endTwirly() {
	if !c.StatusIsActive(twirlyKey) && !c.nightsoulState.HasBlessing() {
		return
	}
	c.DeleteStatus(twirlyKey)
	c.nightsoulState.ExitBlessing()
	c.mounted = false
	c.twirlySrc = -1
}

func (c *char) consumeTwirlyPoints(amount float64) {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.ConsumePoints(amount)
	if c.nightsoulState.Points() <= 0.001 {
		c.endTwirly()
	}
}

func (c *char) applyTwirlyFlat(ai *info.AttackInfo) {
	if c.Base.Ascension >= 4 {
		ai.FlatDmg += c.TotalDef(false) * 0.2
	}
}

func (c *char) twirlyAttackInfo(abil string, mult float64) info.AttackInfo {
	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           abil,
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagElementalArt,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeBlunt,
		PoiseDMG:       100,
		Element:        attributes.Geo,
		Durability:     25,
		Mult:           mult,
		UseDef:         true,
		IgnoreInfusion: true,
	}
	c.applyTwirlyFlat(&ai)
	return ai
}

func (c *char) twirlyIndependentAttackInfo(abil string, mult float64) info.AttackInfo {
	ai := c.twirlyAttackInfo(abil, mult)
	ai.ICDTag = attacks.ICDTagNone
	ai.PoiseDMG = 75
	return ai
}

func (c *char) twirlyRadius() float64 {
	if c.StatusIsActive(fieldKey) {
		return 5.2
	}
	return 4
}
