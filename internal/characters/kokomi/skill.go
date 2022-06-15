package kokomi

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const skillHitmark = 24

func init() {
	skillFrames = frames.InitAbilSlice(61)
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 29
}

// Skill handling - Handles primary damage instance
// Deals Hydro DMG to surrounding opponents and heal nearby active characters once every 2s. This healing is based on Kokomi's Max HP.
func (c *char) Skill(p map[string]int) action.ActionInfo {
	// skill duration is ~12.5s
	// Plus 1 to avoid same frame issues with skill ticks
	c.Core.Status.Add("kokomiskill", 12*60+30+1)

	d := c.createSkillSnapshot()

	// You get 1 tick immediately, then 1 tick every 2 seconds for a total of 7 ticks
	c.skillLastUsed = c.Core.F
	c.Core.Tasks.Add(func() { c.skillTick(d) }, skillHitmark)
	c.Core.Tasks.Add(c.skillTickTask(d, c.Core.F), skillHitmark+126)

	c.SetCDWithDelay(action.ActionSkill, 20*60, 20)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillHitmark,

		State: action.SkillState,
	}
}

// Helper function since this needs to be created both on skill use and burst use
func (c *char) createSkillSnapshot() *combat.AttackEvent {

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Bake-Kurage",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skillDmg[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)

	return (&combat.AttackEvent{
		Info:        ai,
		Pattern:     combat.NewDefCircHit(5, false, combat.TargettableEnemy),
		SourceFrame: c.Core.F,
		Snapshot:    snap,
	})

}

// Helper function that handles damage, healing, and particle components of every tick of her E
func (c *char) skillTick(d *combat.AttackEvent) {

	// check if skill has burst bonus snapshot
	// max swap frame should be 40 frame before 2nd tick
	if c.swapEarlyF > c.skillLastUsed && c.swapEarlyF < (c.skillLastUsed+120-40) {
		d.Info.FlatDmg = c.skillFlatDmg
	} else {
		d.Info.FlatDmg = c.burstDmgBonus(d.Info.AttackTag)
	}

	maxhp := c.MaxHP()

	c.Core.QueueAttackEvent(d, 0)
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  c.Core.Player.Active(),
		Message: "Bake-Kurage",
		Src:     skillHealPct[c.TalentLvlSkill()]*maxhp + skillHealFlat[c.TalentLvlSkill()],
		Bonus:   d.Snapshot.Stats[attributes.Heal],
	})

	// Particles are 0~1 (1:2) on every damage instance
	if c.Core.Rand.Float64() < .6667 {
		c.Core.QueueParticle("kokomi", 1, attributes.Hydro, 100)
	}

	// C2 handling - believe this is an additional instance of flat healing
	// Sangonomiya Kokomi gains the following Healing Bonuses with regard to characters with 50% or less HP via the following methods:
	// Kurage's Oath Bake-Kurage: 4.5% of Kokomi's Max HP.
	if c.Base.Cons >= 2 {
		active := c.Core.Player.ActiveChar()
		if active.HPCurrent/active.MaxHP() <= .5 {
			c.Core.Player.Heal(player.HealInfo{
				Caller:  c.Index,
				Target:  c.Core.Player.Active(),
				Message: "The Clouds Like Waves Rippling",
				Src:     0.045 * maxhp,
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}
}

// Handles repeating skill damage ticks. Split into a separate function as you can only have 1 jellyfish on field at once
// Skill snapshots, so inputs into the function are the originating snapshot
func (c *char) skillTickTask(originalSnapshot *combat.AttackEvent, src int) func() {
	return func() {
		c.Core.Log.NewEvent("Skill Tick Debug", glog.LogCharacterEvent, c.Index, "current dur", c.Core.Status.Duration("kokomiskill"), "skilllastused", c.skillLastUsed, "src", src)
		if c.Core.Status.Duration("kokomiskill") == 0 {
			return
		}

		// Basically stops "old" casts of E from working, and also stops further ticks from that source
		if c.skillLastUsed > src {
			return
		}

		c.skillTick(originalSnapshot)

		c.Core.Tasks.Add(c.skillTickTask(originalSnapshot, src), 120)
	}
}
