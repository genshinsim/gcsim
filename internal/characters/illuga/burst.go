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
	orioleSongKey = "haunted-night-oriole-song"
	burstHitmark  = 48
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
		Abil:       "Dawnbearing Songbird Tap",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 25,
		UseEM:      true,
		Mult:       burst_em[c.TalentLvlBurst()],
	}

	ai.FlatDmg += burst_def[c.TalentLvlBurst()] + c.TotalDef(false)

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6.5)

	c.AddStatus(orioleSongKey, 20*60, true)

	c.c2QuillCounter = 0

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

	c.c4Src = c.Core.F

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
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
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

		if !c.StatusIsActive(orioleSongKey) {
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
			amt = burst_buff_geo[c.TalentLvlSkill()]*c.Stat(attributes.EM) + c.a4GeoBonus()
		case attacks.AttackTagDirectLunarCrystallize,
			attacks.AttackTagReactionLunarCrystallize:
			amt = burst_buff_lcr[c.TalentLvlSkill()]*c.Stat(attributes.EM) + c.a4LcrBonus()
		default:
			return
		}

		c.nightingalesSong--

		c.c2QuillCounter++

		atk.Info.FlatDmg += amt

		if c.c2QuillCounter >= c2QuillThreshold {
			c.c2QuillCounter -= c2QuillThreshold
			c.c2()
		}

		if c.nightingalesSong < 1 {
			c.DeleteStatus(orioleSongKey)
		}

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Illuga Quill proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(orioleSongKey)).
				Write("quill_left", c.nightingalesSong)
		}
	}, "illuga-burst-quill")

	c.Core.Events.Subscribe(event.OnConstructSpawned, func(args ...any) {
		if c.nightingalesSongExtraConstruct < 1 {
			return
		}

		if c.StatusIsActive(orioleSongKey) {
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
