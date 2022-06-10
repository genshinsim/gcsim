package yunjin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Yunjin, NewChar)
}

type char struct {
	*tmpl.Character
	burstTriggers       [4]int
	partyElementalTypes int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	c.partyElementalTypes = 0
	for i := range c.burstTriggers {
		c.burstTriggers[i] = 30
	}

	w.Character = &c

	return nil
}

// Occurs after all characters are loaded, so getPartyElementalTypeCounts works properly
func (c *char) Init() error {
	c.getPartyElementalTypeCounts()
	c.burstProc()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 22)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 33)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 39)

	// charge -> x
	chargeFrames = frames.InitAbilSlice(66)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 66

	// skill -> x
	skillFrames = make([][]int, 3)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(31)

	// skill (level=1) -> x
	skillFrames[1] = frames.InitAbilSlice(81)

	// skill (level=2) -> x
	skillFrames[2] = frames.InitAbilSlice(121)

	// burst -> x
	burstFrames = frames.InitAbilSlice(53)
}

// Adds event to get the number of elemental types in the party for Yunjin A4
func (c *char) getPartyElementalTypeCounts() {
	partyElementalTypes := make(map[attributes.Element]int)
	for _, char := range c.Core.Player.Chars() {
		partyElementalTypes[char.Base.Element]++
	}
	for range partyElementalTypes {
		c.partyElementalTypes += 1
	}
	c.Core.Log.NewEvent("Yun Jin Party Elemental Types (A4)", glog.LogCharacterEvent, c.Index, "party_elements", c.partyElementalTypes)
}
