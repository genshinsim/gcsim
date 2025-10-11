package ineffa

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

const (
	strStackKey = "strStack"
)

func init() {
	core.RegisterCharFunc(keys.Ineffa, NewChar)
}

type char struct {
	*tmpl.Character
	birgittaSrc int
	skillShield *shd
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Snapshot(a *info.AttackInfo) info.Snapshot {
	s := info.Snapshot{
		CharLvl:     c.Base.Level,
		SourceFrame: c.Core.F,
	}

	var evt glog.Event
	var debug []any

	if c.Core.Flags.LogDebug {
		evt = c.Core.Log.NewEvent(a.Abil, glog.LogSnapshotEvent, c.Index()).
			Write("abil", a.Abil).
			Write("mult", a.Mult).
			Write("ele", a.Element.String()).
			Write("durability", float64(a.Durability)).
			Write("icd_tag", a.ICDTag).
			Write("icd_group", a.ICDGroup)
	}

	// snapshot the stats
	s.Stats, debug = c.Stats()

	// check infusion
	var inf attributes.Element
	if !a.IgnoreInfusion {
		inf = c.Core.Player.Infused(c.Index(), a.AttackTag)
		if inf != attributes.NoElement {
			a.Element = inf
		}
	}

	// check if we need to log
	if c.Core.Flags.LogDebug {
		evt.WriteBuildMsg(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != attributes.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}

func (c *char) Init() error {
	c.a4Init()
	c.lunarchargeInit()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 27
	}
	return c.Character.AnimationStartDelay(k)
}
