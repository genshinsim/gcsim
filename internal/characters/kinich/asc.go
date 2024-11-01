package kinich

import "github.com/genshinsim/gcsim/pkg/core/event"

const (
	desolationKey = "kinich-a1-desolation"
	a1Icd         = "kinich-a1-desolation-icd"
	a4StackKey    = "kinich-a4-stack-key"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	hook := func(args ...interface{}) bool {
		if !c.StatusIsActive(desolationKey) {
			return false
		}
		if c.StatusIsActive(a1Icd) {
			return false
		}
		c.nightsoulState.GeneratePoints(7)
		c.AddStatus(a1Icd, 0.8*60, false)
		return false
	}
	c.Core.Events.Subscribe(event.OnBurning, hook, "kinich-a1-burning")
	c.Core.Events.Subscribe(event.OnBurgeon, hook, "kinich-a1-burgeon")
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		stacks := c.Tags[a4StackKey]
		stacks = min(stacks+1, 2)
		c.AddStatus(a4StackKey, 15*60, true)
		c.SetTag(a4StackKey, stacks)
		return false
	}, "kinich-a4")
}

func (c *char) a4Amount() float64 {
	if c.Base.Ascension < 4 {
		return 0.0
	}
	stacks := c.Tags[a4StackKey]
	c.Tags[a4StackKey] = 0
	c.DeleteStatus(a4StackKey)
	return 3.2 * float64(stacks) * c.TotalAtk()
}
