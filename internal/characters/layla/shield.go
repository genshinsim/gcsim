package layla

import (
	"sort"

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

var starsTravel = []int{35, 33, 30, 28}

type ICDNightStar int

const (
	ICDNightStarNone  ICDNightStar = iota
	ICDNightStarSkill              // 0.3s
	ICDNightStarBurst              // 0.5s
)

func (c *char) removeShield() {
	if c.Tag(shootingStars) == 0 {
		c.SetTag(nightStars, 0)
	}
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

func (c *char) addNightStars(count int, icd ICDNightStar) {
	if c.Tag(shootingStars) > 0 {
		return
	}

	switch icd {
	case ICDNightStarSkill:
		if c.StatusIsActive(starSkillIcd) {
			return
		}
		c.AddStatus(starSkillIcd, 0.3*60, false)
	case ICDNightStarBurst:
		if c.StatusIsActive(starBurstIcd) {
			return
		}
		c.AddStatus(starBurstIcd, 0.5*60, false)
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

	if stars == 4 && c.Tag(shootingStars) == 0 {
		c.SetTag(shootingStars, 1)
		c.shootStarSrc = c.Core.F
		c.Core.Tasks.Add(c.shootStars(c.shootStarSrc, -1), 0.1*60)
	}
}

func (c *char) shootStars(src int, last int) func() {
	return func() {
		if c.shootStarSrc != src {
			return
		}
		if c.Tag(shootingStars) == 0 {
			return
		}

		// find near target
		nearTarget := -1
		trgs := c.Core.Combat.EnemiesWithinRadius(c.Core.Combat.Player().Pos(), 10)
		if len(trgs) > 0 {
			sort.Slice(trgs, func(i, j int) bool { return i < j })
			nearTarget = trgs[0]
		}

		if last == -1 {
			// if not found
			if nearTarget == -1 {
				c.Core.Tasks.Add(c.shootStars(src, -1), 0.1*60)
				return
			}
			c.Core.Tasks.Add(c.shootStars(src, nearTarget), 0.5*60)
			return
		}
		if nearTarget == -1 {
			nearTarget = last
		}

		stars := c.Tag(nightStars)
		if stars == 4 && c.Base.Cons >= 4 {
			for _, char := range c.Core.Player.Chars() {
				char.AddStatus(c4Key, 3*60, true)
			}
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Shooting Star",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupLayla,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       starDmg[c.TalentLvlSkill()],
			FlatDmg:    0.015 * c.MaxHP(), // A4
		}

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
				c.AddStatus(skillEnergy, 3.5*60, true)
			}
			if c.Base.Cons >= 2 {
				c.AddEnergy("layla-c2", 1)
			}
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.Enemy(nearTarget),
				nil,
				0.8,
			),
			0,
			starsTravel[len(starsTravel)-stars],
			cb,
		)

		stars--
		c.SetTag(nightStars, stars)
		if stars != 0 {
			c.Core.Tasks.Add(c.shootStars(src, nearTarget), 0.45*60)
			return
		}

		c.RemoveTag(shootingStars)
		c.starTickSrc = c.Core.F
		c.tickNightStar(c.starTickSrc, false)
	}
}

func (c *char) tickNightStar(src int, star bool) func() {
	return func() {
		if c.starTickSrc != src {
			return
		}
		exist := c.Core.Player.Shields.Get(shield.ShieldLaylaSkill)
		if exist == nil {
			return
		}

		if star {
			c.addNightStars(1, ICDNightStarSkill)
		}

		interval := 1.5 * 60
		if c.Base.Cons >= 6 {
			interval = 1.5 * 0.8 * 60
		}
		c.QueueCharTask(c.tickNightStar(src, true), int(interval))
	}
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
