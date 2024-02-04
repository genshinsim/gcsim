package xianyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var burstFrames []int

const (
	burstStart   = 47
	burstHitmark = 78
	burstKey     = "xianyun-burst"
)

// TODO: dummy frame data from shenhe
func init() {
	burstFrames = frames.InitAbilSlice(100) // Q -> E
	burstFrames[action.ActionAttack] = 99   // Q -> N1
	burstFrames[action.ActionDash] = 78     // Q -> D
	burstFrames[action.ActionJump] = 79     // Q -> J
	burstFrames[action.ActionWalk] = 98     // Q -> Walk
	burstFrames[action.ActionSwap] = 98     // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Stars Gather at Dusk (Initial)",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}
	burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 7}, 4)
	c.Core.QueueAttack(ai, burstArea, burstHitmark, burstHitmark)

	// init heal
	stats, _ := c.Stats()
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Starwicker-Heal-Initial",
		Src:     instantHealp[c.TalentLvlBurst()]*((c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP])+stats[attributes.ATK]) + instantHealFlat[c.TalentLvlBurst()],
		Bonus:   c.Stat(attributes.Heal),
	})

	// 16 seconds duration
	burstDuration := 16 * 60

	c.AddStatus(burstKey, burstDuration, false)
	c.AddStatus(starwickerKey, burstDuration, true)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(4)

	c.plungeDoTTrigger()
	c.a4()
	for i := burstStart; i <= burstStart+burstDuration; i += 2.5 * 60 {
		c.Core.Tasks.Add(c.BurstHealDoT, i+2.5*60)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) plungeDoTTrigger() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		burstArea := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2}, 1)
		burstPos := burstArea.Shape.Pos()
		if atk.Info.AttackTag != attacks.AttackTagPlunge {
			return false
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagPlunge:
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Starwicker Plunge DoT Damage",
				AttackTag:  attacks.AttackTagElementalBurst,
				ICDTag:     attacks.ICDTagElementalBurst,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Anemo,
				Durability: 25,
				Mult:       burstdot[c.TalentLvlBurst()],
			}
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHitOnTargetFanAngle(burstPos, nil, 8, 120),
				5,
				5,
			)
			return false
		default:
			return false
		}
	}, "xianyun-starwicker-plunge-DoT-hook")
}

func (c *char) BurstHealDoT() {
	stats, _ := c.Stats()
	c.Core.Player.Heal(player.HealInfo{
		Caller:  c.Index,
		Target:  -1,
		Message: "Starwicker-Heal-DoT",
		Src:     healdotp[c.TalentLvlBurst()]*((c.Base.Atk+c.Weapon.BaseAtk)*(1+stats[attributes.ATKP])+stats[attributes.ATK]) + healdotflat[c.TalentLvlBurst()],
		Bonus:   c.Stat(attributes.Heal),
	})
}
