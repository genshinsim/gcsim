package hutao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c6ICDKey = "hutao-c6-icd"
)

func (c *char) c6() {
	c.c6buff = make([]float64, attributes.EndStatType)
	c.c6buff[attributes.CR] = 1
	// check for C6 proc on hurt
	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...interface{}) bool {
		di := args[0].(player.DrainInfo)
		if di.Amount <= 0 {
			return false
		}
		c.checkc6(false)
		return false
	}, "hutao-c6")
	// check for C6 proc every 2s regardless of hurt
	c.QueueCharTask(func() {
		c.checkc6(true)
	}, 120)
}

func (c *char) checkc6(check1HP bool) {
	// check for C6 proc every 2s regardless of hurt and c6 icd
	c.QueueCharTask(func() {
		c.checkc6(true)
	}, 120)
	// check if c6 is on icd
	if c.StatusIsActive(c6ICDKey) {
		return
	}
	// check if hp less than 25%
	if c.CurrentHPRatio() > 0.25 {
		return
	}
	// check if hp is less than 2 for the 2s check
	if check1HP && c.CurrentHP() >= 2 {
		return
	}
	// if dead, revive back to 1 hp
	if c.CurrentHPRatio() <= 0 {
		c.SetHPByAmount(1)
	}

	//increase crit rate to 100%
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("hutao-c6", 600),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			return c.c6buff, true
		},
	})

	c.AddStatus(c6ICDKey, 3600, false)
}

// Upon defeating an enemy affected by a Blood Blossom that Hu Tao applied
// herself, all nearby allies in the party (excluding Hu Tao herself) will have
// their CRIT Rate increased by 12% for 15s.
func (c *char) c4() {
	c.c4buff = make([]float64, attributes.EndStatType)
	c.c4buff[attributes.CR] = 0.12
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		//do nothing if not an enemy
		if !ok {
			return false
		}
		if !t.StatusIsActive(bbDebuff) {
			return false
		}
		for i, char := range c.Core.Player.Chars() {
			//does not affect hutao
			if c.Index == i {
				continue
			}
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag("hutao-c4", 900),
				AffectedStat: attributes.CR,
				Amount: func() ([]float64, bool) {
					return c.c4buff, true
				},
			})
		}

		return false
	}, "hutao-c4")
}
