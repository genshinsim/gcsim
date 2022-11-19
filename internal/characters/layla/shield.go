package layla

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

const (
	nightStars    = "nightstars"
	starSkillIcd  = "nightstar-skill-icd"
	starBurstIcd  = "nightstar-burst-icd"
	shootingStars = "shooting-stars"
)

func (c *char) removeShield() {
	c.SetTag(nightStars, 0)
	c.a1Stack = 0
}

func (c *char) newShield(base float64, dur int) *shd {
	n := &shd{}
	n.Tmpl = &shield.Tmpl{}
	n.Tmpl.Src = c.Core.F
	n.Tmpl.ShieldType = shield.ShieldLaylaSkill
	n.Tmpl.Ele = attributes.Cryo
	n.Tmpl.HP = base
	n.Tmpl.Name = "Layla Skill"
	n.Tmpl.Expires = c.Core.F + dur
	n.c = c
	return n
}

func (c *char) AddNightStars(count int, icd bool) {
	if icd {
		if c.StatusIsActive(starSkillIcd) {
			return
		}
		c.AddStatus(starSkillIcd, 0.3*60, true)
	}

	if c.a1Stack < 4 {
		c.a1Stack++
	}

	stars := count + c.Tag(nightStars)
	if stars > 4 {
		stars = 4
	}
	c.SetTag(nightStars, stars)
	c.Core.Log.NewEvent("adding stars", glog.LogCharacterEvent, c.Index).
		Write("stars", stars)
}

func (c *char) TickNightStar(star bool) {
	exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
	if exist == nil {
		return
	}

	if star {
		c.AddNightStars(1, true)
	}
	delay := 0
	if c.Tag(nightStars) == 4 {
		delay = 3 * c.starTravel
		c.AddStatus(shootingStars, delay, true)
		c.SetTag(nightStars, 0)

		if c.Base.Cons >= 4 {
			for _, char := range c.Core.Player.Chars() {
				char.AddStatus(c4Key, 3*60, true)
			}
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shooting Star",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       starDmg[c.TalentLvlSkill()],
			FlatDmg:    0.015 * c.MaxHP(), // A4
		}
		// TODO: snapshot?
		snap := c.Snapshot(&ai)

		for i := 1; i <= 4; i++ {
			c.Core.Tasks.Add(func() {
				done := false
				cb := func(_ combat.AttackCB) {
					if done {
						return
					}
					done = true

					if !c.StatusIsActive(skillEnergy) {
						var count float64 = 1
						if c.Core.Rand.Float64() < 0.33 {
							count = 2
						}
						c.Core.QueueParticle("layla", count, attributes.Cryo, c.ParticleDelay)
						c.AddStatus(skillEnergy, 3*60, true)
					}
					if c.Base.Cons >= 2 {
						c.AddEnergy("layla-c2", 1)
					}
				}

				c.Core.QueueAttackWithSnap(ai, snap, combat.NewCircleHit(c.Core.Combat.PrimaryTarget(), 0.5), 0, cb)
			}, i*c.starTravel)
		}
	}

	interval := 1.5
	if c.Base.Cons >= 6 {
		interval = 1.5 * 0.2
	}
	c.QueueCharTask(func() { c.TickNightStar(true) }, int(interval*60)+delay)
}

type shd struct {
	*shield.Tmpl
	c *char
}

func (s *shd) OnExpire() {
	s.c.removeShield()
}

func (s *shd) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	taken, ok := s.Tmpl.OnDamage(dmg, ele, bonus)
	if !ok {
		s.c.removeShield()
	}
	return taken, ok
}
