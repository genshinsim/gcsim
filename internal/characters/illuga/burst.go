package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstKey     = "haunted-night-oriole-song"
	burstHitmark = 48
)

func init() {
	burstFrames = frames.InitAbilSlice(65) // Q -> N1
	burstFrames[action.ActionDash] = 64
	burstFrames[action.ActionJump] = 64
	burstFrames[action.ActionWalk] = 64
	burstFrames[action.ActionSwap] = 62
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Shadowless Reflection",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   150,
		Element:    attributes.Geo,
		Durability: 25,
	}
	c.Core.Tasks.Add(func() {
		ai.FlatDmg += burst_em[c.TalentLvlBurst()] * c.Stat(attributes.EM)
		ai.FlatDmg += burst_def[c.TalentLvlBurst()] * c.TotalDef(false)
	}, burstHitmark)

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5)

	c.AddStatus(burstKey, 20*60, true)

	c.c2Reset()

	c.nightingalesSong = 21

	c.nightingalesSongExtraConstruct = 15

	_, constructs := c.Core.Constructs.ConstructsByType(construct.GeoConstructInvalid)

	playerPos := c.Core.Combat.Player().Pos()

	for _, construct := range constructs {
		if c.nightingalesSongExtraConstruct < 1 {
			break
		}

		if playerPos.Distance(construct.Pos()) > 30 {
			continue
		}

		c.nightingalesSongExtraConstruct -= 5

		c.nightingalesSong += 5
	}

	c.a1()

	c.c4(c.Core.F)()

	c.SetCD(action.ActionBurst, 15*60)

	c.ConsumeEnergy(6)

	c.Core.QueueAttack(
		ai,
		ap,
		burstHitmark,
		burstHitmark,
	)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstBuffInit() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.Element != attributes.Geo {
			return
		}

		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if !c.StatusIsActive(burstKey) {
			return
		}

		if c.nightingalesSong < 1 {
			return
		}

		var amt float64

		switch atk.Info.AttackTag {
		case attacks.AttackTagElementalBurst,
			attacks.AttackTagElementalArt,
			attacks.AttackTagElementalArtHold,
			attacks.AttackTagNormal,
			attacks.AttackTagExtra,
			attacks.AttackTagPlunge:
			amt = (burst_buff_geo[c.TalentLvlBurst()] + c.a4GeoBonus()) * c.Stat(attributes.EM)
		case attacks.AttackTagDirectLunarCrystallize:
			amt = (burst_buff_lcr[c.TalentLvlBurst()] + c.a4LcrBonus()) * c.Stat(attributes.EM)
		default:
			return
		}

		c.nightingalesSong--

		atk.Info.FlatDmg += amt

		c.c2Increment()

		if c.nightingalesSong < 1 {
			c.DeleteStatus(burstKey)
		}

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Illuga Quill proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(burstKey)).
				Write("quill_left", c.nightingalesSong)
		}
	}, "illuga-burst-quill")

	c.Core.Events.Subscribe(event.OnConstructSpawned, func(args ...any) {
		if c.nightingalesSongExtraConstruct < 1 {
			return
		}

		if c.StatusIsActive(burstKey) {
			return
		}

		construct, ok := args[0].(construct.Construct)

		if !ok {
			return
		}

		playerPos := c.Core.Combat.Player().Pos()

		dist := playerPos.Distance(construct.Pos())

		if dist > 30 {
			return
		}

		c.nightingalesSongExtraConstruct -= 5

		c.nightingalesSong += 5
	}, "illuga-burst-gain-quills-on-construct")
}
