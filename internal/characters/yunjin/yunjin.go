package yunjin

import (
	"fmt"
	"github.com/genshinsim/gcsim/pkg/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func init() {
	core.RegisterCharFunc(keys.Yunjin, NewChar)
}

type char struct {
	*character.Tmpl

	burstTriggers       [4]int
	partyElementalTypes int
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassSpear
	c.NormalHitNum = 5
	c.BurstCon = 3
	c.SkillCon = 5

	for i := range c.burstTriggers {
		c.burstTriggers[i] = 30
	}

	c.getPartyElementalTypeCounts()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	c.burstProc()

	return &c, nil
}

// Occurs after all characters are loaded, so getPartyElementalTypeCounts works properly
func (c *char) Init(index int) {
	c.Tmpl.Init(index)

	c.partyElementalTypes = 0
	c.getPartyElementalTypeCounts()
}

// Helper function to update tags that can be used in configs
// Should be run whenever c.burstTriggers is updated
func (c *char) updateBuffTags() {
	for _, char := range c.Core.Chars {
		c.Tags["burststacks_"+char.Name()] = c.burstTriggers[char.CharIndex()]
		c.Tags[fmt.Sprintf("burststacks_%v", char.CharIndex())] = c.burstTriggers[char.CharIndex()]
	}
}

// Adds event to get the number of elemental types in the party for Yunjin A4
func (c *char) getPartyElementalTypeCounts() {
	partyElementalTypes := make(map[core.EleType]int)
	for _, char := range c.Core.Chars {
		partyElementalTypes[char.Ele()]++
	}
	for i := range partyElementalTypes {
		c.partyElementalTypes += 1
		// Is there a more elegant way to get go to not complain about variable not used?
		i += 0
	}
	c.Core.Log.Debugw("Yun Jin Party Elemental Types (A4)", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "party_elements", c.partyElementalTypes)
}

// When Yun Jin triggers the Crystallize Reaction, her DEF is increased by 20% for 12s.
func (c *char) c4() {
	charModFunc := func(args ...interface{}) bool {
		ae := args[1].(*core.AttackEvent)

		if ae.Info.ActorIndex != c.CharIndex() {
			return false
		}

		val := make([]float64, core.EndStatType)
		val[core.DEFP] = .2
		c.AddMod(core.CharStatMod{
			Key:    "yunjin-c4",
			Expiry: c.Core.F + 12*60,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		return false
	}
	c.Core.Events.Subscribe(core.OnCrystallizeCryo, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(core.OnCrystallizeElectro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(core.OnCrystallizePyro, charModFunc, "yunjin-c4")
	c.Core.Events.Subscribe(core.OnCrystallizeHydro, charModFunc, "yunjin-c4")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.Core.Log.Warnw("ActionStam not implemented", "character", c.Base.Key.String())
		return 0
	}
}
