package gorou

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

//When characters (other than Gorou) within the AoE of Gorou's General's War Banner
//or General's Glory deal Geo DMG to opponents, the CD of Gorou's Inuzaka All-Round Defense
//is decreased by 2s. This effect can occur once every 10s.
func (c *char) c1() {
	icd := -1
	c.Core.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		if c.Core.StatusDuration(generalGloryKey) == 0 && c.Core.StatusDuration(generalWarBannerKey) == 0 {
			return false
		}
		if icd > c.Core.Frame {
			return false
		}
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.ActorIndex == c.Index {
			return false
		}
		if atk.Info.Element != core.Geo {
			return false
		}
		icd = c.Core.Frame + 600
		c.ReduceActionCooldown(core.ActionSkill, 120)
		return false
	}, "gorou-c1")
}

//While General's Glory is in effect, its duration is extended by 1s when a nearby
//active character obtains an Elemental Shard from a Crystallize reaction.
//This effect can occur once every 0.1s. Max extension is 3s.
func (c *char) c2() {
	//TODO: this is currently on reaction but really should be on pick up
	cb := func(args ...interface{}) bool {
		dur := c.Core.StatusDuration(generalGloryKey)
		if dur == 0 {
			return false
		}
		if c.c2Extension == 180 {
			return false
		}
		//to simulate pickup we add 30 frames delay
		c.Core.Tasks.Add(func() {
			ext := 60
			if c.c2Extension+ext > 180 {
				ext = 180 - c.c2Extension
			}
			c.c2Extension += ext
			c.Core.AddStatus(generalGloryKey, c.Core.Frame+dur+ext)
		}, 30)
		return false
	}
	c.Core.Subscribe(core.OnCrystallizeCryo, cb, "gorou-c2")
	c.Core.Subscribe(core.OnCrystallizeElectro, cb, "gorou-c2")
	c.Core.Subscribe(core.OnCrystallizeHydro, cb, "gorou-c2")
	c.Core.Subscribe(core.OnCrystallizePyro, cb, "gorou-c2")
}

func (c *char) c6() {
	for _, char := range c.Core.Chars {
		char.AddPreDamageMod(coretype.PreDamageMod{
			Key:    c6key,
			Expiry: c.Core.Frame + 720, //12s
			Amount: func(ae *coretype.AttackEvent, t coretype.Target) ([]float64, bool) {
				if ae.Info.Element != core.Geo {
					return nil, false
				}
				return c.c6buff, true
			},
		})
	}
}
