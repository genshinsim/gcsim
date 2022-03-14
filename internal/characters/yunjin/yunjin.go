package yunjin

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/tmpl/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func init() {
	core.RegisterCharFunc(core.Yunjin, NewChar)
}

type char struct {
	*character.Tmpl

	burstTriggers       [4]int
	partyElementalTypes int
}

func NewChar(s *core.Core, p coretype.CharacterProfile) (coretype.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Base.Element = core.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
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
func (c *char) Init() {
	c.Tmpl.Init()

	c.partyElementalTypes = 0
	c.getPartyElementalTypeCounts()
}

// Helper function to update tags that can be used in configs
// Should be run whenever c.burstTriggers is updated
func (c *char) updateBuffTags() {
	for _, char := range c.Core.Chars {
		c.Tags["burststacks_"+char.Name()] = c.burstTriggers[char.Index()]
		c.Tags[fmt.Sprintf("burststacks_%v", char.Index())] = c.burstTriggers[char.Index()]
	}
}

// Adds event to get the number of elemental types in the party for Yunjin A4
func (c *char) getPartyElementalTypeCounts() {
	partyElementalTypes := make(map[coretype.EleType]int)
	for _, char := range c.Core.Chars {
		partyElementalTypes[char.Ele()]++
	}
	for i := range partyElementalTypes {
		c.partyElementalTypes += 1
		// Is there a more elegant way to get go to not complain about variable not used?
		i += 0
	}
	c.coretype.Log.NewEvent("Yun Jin Party Elemental Types (A4)", coretype.LogCharacterEvent, c.Index, "party_elements", c.partyElementalTypes)
}

// When Yun Jin triggers the Crystallize Reaction, her DEF is increased by 20% for 12s.
func (c *char) c4() {
	charModFunc := func(args ...interface{}) bool {
		ae := args[1].(*coretype.AttackEvent)

		if ae.Info.ActorIndex != c.Index() {
			return false
		}

		val := make([]float64, core.EndStatType)
		val[core.DEFP] = .2
		c.AddMod(coretype.CharStatMod{
			Key:    "yunjin-c4",
			Expiry: c.Core.Frame + 12*60,
			Amount: func() ([]float64, bool) {
				return val, true
			},
		})
		return false
	}
	c.Core.Subscribe(core.OnCrystallizeCryo, charModFunc, "yunjin-c4")
	c.Core.Subscribe(core.OnCrystallizeElectro, charModFunc, "yunjin-c4")
	c.Core.Subscribe(core.OnCrystallizePyro, charModFunc, "yunjin-c4")
	c.Core.Subscribe(core.OnCrystallizeHydro, charModFunc, "yunjin-c4")
}

func (c *char) ActionStam(a core.ActionType, p map[string]int) float64 {
	switch a {
	case core.ActionDash:
		return 18
	default:
		c.coretype.Log.NewEvent("ActionStam not implemented", coretype.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}
