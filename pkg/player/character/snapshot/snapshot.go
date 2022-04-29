package snapshot

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/player"
)

type ModsHandler interface {
	StatsMods(char int) ([attributes.EndStatType]float64, []interface{})
	StatMod(char int, s attributes.Stat) float64
}

type InfusionHandler interface {
	Infused(char int, a combat.AttackTag) attributes.Element
}

type Handler struct {
	mods ModsHandler
	inf  InfusionHandler
}

func (h *Handler) Snapshot(m *player.MasterChar, a *combat.AttackInfo) combat.Snapshot {

	s := combat.Snapshot{
		CharLvl:     m.Base.Level,
		ActorEle:    m.Base.Element,
		BaseAtk:     m.Base.Atk + m.Weapon.Atk,
		BaseDef:     m.Base.Def,
		SourceFrame: m.Player.Core.F,
	}

	var evt glog.Event = nil
	var debug []interface{}

	if m.Player.Core.Flags.LogDebug {
		evt = m.Player.Core.Log.NewEvent(
			a.Abil, glog.LogSnapshotEvent, m.Index,
			"abil", a.Abil,
			"mult", a.Mult,
			"ele", a.Element.String(),
			"durability", float64(a.Durability),
			"icd_tag", a.ICDTag,
			"icd_group", a.ICDGroup,
		)
	}

	//snapshot the stats
	copy(s.Stats[:], m.Stats[:attributes.EndStatType])
	mods, debug := h.mods.StatsMods(m.Index)
	for i, v := range mods {
		s.Stats[i] += v
	}

	//check infusion
	var inf attributes.Element
	if !a.IgnoreInfusion {
		inf = h.inf.Infused(m.Index, a.AttackTag)
		if inf != attributes.NoElement {
			a.Element = inf
		}
	}

	//check if we need to log
	if m.Player.Core.Flags.LogDebug {
		evt.Write(debug...)
		evt.Write("final_stats", attributes.PrettyPrintStatsSlice(s.Stats[:]))
		if inf != attributes.NoElement {
			evt.Write("infused_ele", inf.String())
		}
		s.Logs = debug
	}
	return s
}
