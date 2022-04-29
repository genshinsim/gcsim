package snapshot

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type SnapshotHandler struct {
	index int
	core  *core.Core
}

func (h *SnapshotHandler) Snapshot(a *combat.AttackInfo) combat.Snapshot {

	char := h.core.Player.ByIndex(h.index)

	s := combat.Snapshot{
		CharLvl:     char.Base.Level,
		ActorEle:    char.Base.Element,
		BaseAtk:     char.Base.Atk + char.Weapon.Atk,
		BaseDef:     char.Base.Def,
		SourceFrame: h.core.F,
	}

	var evt glog.Event = nil
	var debug []interface{}

	if h.core.Flags.LogDebug {
		evt = h.core.Log.NewEvent(
			a.Abil, glog.LogSnapshotEvent, h.index,
			"abil", a.Abil,
			"mult", a.Mult,
			"ele", a.Element.String(),
			"durability", float64(a.Durability),
			"icd_tag", a.ICDTag,
			"icd_group", a.ICDGroup,
		)
	}

	//snapshot the stats
	s.Stats, debug = h.core.Player.ByIndex(h.index).Stats()

	//check infusion
	var inf attributes.Element
	if !a.IgnoreInfusion {
		inf = h.core.Player.Infused(h.index, a.AttackTag)
		if inf != attributes.NoElement {
			a.Element = inf
		}
	}

	//check if we need to log
	if h.core.Flags.LogDebug {
		evt.Write(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != attributes.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}
