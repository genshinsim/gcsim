package character

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type Character struct {
	Core                   *core.Core
	Index                  *int
	ActionCD               []int
	cdQueueWorkerStartedAt []int
	cdCurrentQueueWorker   []*func()
	cdQueue                [][]int
	AvailableCDCharge      []int
	additionalCDCharge     []int
}

func (c *Character) Snapshot(a *combat.AttackInfo) combat.Snapshot {

	char := c.Core.Player.ByIndex(*c.Index)

	s := combat.Snapshot{
		CharLvl:     char.Base.Level,
		ActorEle:    char.Base.Element,
		BaseAtk:     char.Base.Atk + char.Weapon.Atk,
		BaseDef:     char.Base.Def,
		SourceFrame: c.Core.F,
	}

	var evt glog.Event = nil
	var debug []interface{}

	if c.Core.Flags.LogDebug {
		evt = c.Core.Log.NewEvent(
			a.Abil, glog.LogSnapshotEvent, *c.Index,
			"abil", a.Abil,
			"mult", a.Mult,
			"ele", a.Element.String(),
			"durability", float64(a.Durability),
			"icd_tag", a.ICDTag,
			"icd_group", a.ICDGroup,
		)
	}

	//snapshot the stats
	s.Stats, debug = c.Core.Player.ByIndex(*c.Index).Stats()

	//check infusion
	var inf attributes.Element
	if !a.IgnoreInfusion {
		inf = c.Core.Player.Infused(*c.Index, a.AttackTag)
		if inf != attributes.NoElement {
			a.Element = inf
		}
	}

	//check if we need to log
	if c.Core.Flags.LogDebug {
		evt.Write(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != attributes.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}
