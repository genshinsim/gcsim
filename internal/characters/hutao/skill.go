package hutao

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillStart = 14

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 29
	skillFrames[action.ActionBurst] = 28
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 37
}

func (c *char) Skill(p map[string]int) action.ActionInfo {

	c.applyA1 = true
	c.Core.Tasks.Add(c.a1, 540+skillStart)
	c.Core.Status.Add("paramita", 540+skillStart) //to account for animation
	c.Core.Log.NewEvent("paramita activated", glog.LogCharacterEvent, c.Index).
		Write("expiry", c.Core.F+540+skillStart)

	//increase based on hp at cast time.
	//figure out atk buff
	c.ppBonus = ppatk[c.TalentLvlSkill()] * c.MaxHP()
	max := (c.Base.Atk + c.Weapon.Atk) * 4
	if c.ppBonus > max {
		c.ppBonus = max
	}

	//remove some hp
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Paramita Papilio",
		Amount:     .30 * c.HPCurrent,
	})
	c.checkc6()

	c.SetCDWithDelay(action.ActionSkill, 960, 14)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) ppParticles(ac combat.AttackCB) {
	if c.Core.Status.Duration("paramita") <= 0 {
		return
	}
	if c.paraParticleICD < c.Core.F {
		c.paraParticleICD = c.Core.F + 300 //5 seconds
		var count float64 = 2
		if c.Core.Rand.Float64() < 0.5 {
			count = 3
		}
		//TODO: this used to be 80
		c.Core.QueueParticle("hutao", count, attributes.Pyro, c.Core.Flags.ParticleDelay)
	}
}

//TODO: this needs to be multi target
func (c *char) applyBB(a combat.AttackCB) {
	c.Core.Log.NewEvent("Applying Blood Blossom", glog.LogCharacterEvent, c.Index).
		Write("current dur", c.Core.Status.Duration("htbb"))
	//check if blood blossom already active, if active extend duration by 8 second
	//other wise start first tick func
	if !c.tickActive {
		//TODO: does BB tick immediately on first application?
		c.Core.Tasks.Add(c.bbtickfunc(c.Core.F, a.Target.Index()), 240)
		c.tickActive = true
		c.Core.Log.NewEvent("Blood Blossom applied", glog.LogCharacterEvent, c.Index).
			Write("expected end", c.Core.F+570).
			Write("next expected tick", c.Core.F+240)
	}
	// c.CD["bb"] = c.Core.F + 570 //TODO: no idea how accurate this is, does this screw up the ticks?
	c.Core.Status.Add("htbb", 570)
	c.Core.Log.NewEvent("Blood Blossom duration extended", glog.LogCharacterEvent, c.Index).
		Write("new expiry", c.Core.Status.Duration("htbb"))
}

func (c *char) bbtickfunc(src, trg int) func() {
	return func() {
		c.Core.Log.NewEvent("Blood Blossom checking for tick", glog.LogCharacterEvent, c.Index).
			Write("cd", c.Core.Status.Duration("htbb")).
			Write("src", src)
		if c.Core.Status.Duration("htbb") == 0 {
			c.tickActive = false
			return
		}
		//queue up one damage instance
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Blood Blossom",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       bb[c.TalentLvlSkill()],
		}
		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			ai.FlatDmg += c.MaxHP() * 0.1
		}
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(1, combat.TargettableEnemy), 0, 0)
		c.Core.Log.NewEvent("Blood Blossom ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+240).
			Write("dur", c.Core.Status.Duration("htbb")).
			Write("src", src)
		//only queue if next tick buff will be active still
		// if c.Core.F+240 > c.CD["bb"] {
		// 	return
		// }
		//queue up next instance
		c.Core.Tasks.Add(c.bbtickfunc(src, trg), 240)

	}
}

func (c *char) ppHook() {
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("hutao-paramita", -1),
		AffectedStat: attributes.ATK,
		Amount: func() ([]float64, bool) {
			if c.Core.Status.Duration("paramita") == 0 {
				return nil, false
			}
			m[attributes.ATK] = c.ppBonus
			return m, true
		},
	})
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.Core.Status.Duration("paramita") > 0 {
			c.a1()
		}
		c.Core.Status.Delete("paramita")
		return false
	}, "hutao-exit")
}
