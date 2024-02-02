package shield

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// Handler keeps track of all the active shields
// Note that there's no need to distinguish between characters here since the shields are shared
// However we do care which is active since only active is shielded
type Handler struct {
	shields []Shield
	log     glog.Logger
	events  event.Eventter
	f       *int

	shieldBonusMods []shieldBonusMod
}

func New(f *int, log glog.Logger, events event.Eventter) *Handler {
	h := &Handler{
		shields:         make([]Shield, 0, EndType),
		log:             log,
		events:          events,
		f:               f,
		shieldBonusMods: make([]shieldBonusMod, 0, 5),
	}
	return h
}

func (h *Handler) Count() int { return len(h.shields) }

func (h *Handler) PlayerIsShielded() bool {
	return len(h.shields) > 0
}

func (h *Handler) Get(t Type) Shield {
	for _, v := range h.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

// TODO: do shields get affected by hitlag? if so.. which timer? active char?
func (h *Handler) Add(shd Shield) {
	// we always assume over write of the same type
	ind := -1
	for i, v := range h.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind > -1 {
		h.log.NewEvent("shield overridden", glog.LogShieldEvent, -1).
			Write("overwrite", true).
			Write("name", shd.Desc()).
			Write("hp", shd.CurrentHP()).
			Write("ele", shd.Element()).
			Write("expiry", shd.Expiry())
		h.shields[ind].OnOverwrite()
		h.shields[ind] = shd
	} else {
		h.shields = append(h.shields, shd)
		h.log.NewEvent("shield added", glog.LogShieldEvent, -1).
			Write("overwrite", false).
			Write("name", shd.Desc()).
			Write("hp", shd.CurrentHP()).
			Write("ele", shd.Element()).
			Write("expiry", shd.Expiry())
	}
	h.events.Emit(event.OnShielded, shd)
}

func (h *Handler) List() []Shield {
	return h.shields
}

func (h *Handler) OnDamage(char int, dmg float64, ele attributes.Element) float64 {
	// find shield bonuses
	bonus := h.ShieldBonus()
	min := dmg // min of damage taken
	n := 0
	for _, v := range h.shields {
		preHp := v.CurrentHP()
		taken, ok := v.OnDamage(dmg, ele, bonus)
		h.log.NewEvent(
			"shield taking damage",
			glog.LogShieldEvent,
			-1,
		).Write("name", v.Desc()).
			Write("ele", v.Element()).
			Write("dmg", dmg).
			Write("previous_hp", preHp).
			Write("dmg_after_shield", taken).
			Write("current_hp", v.CurrentHP()).
			Write("expiry", v.Expiry())
		if taken < min {
			min = taken
		}
		if ok {
			h.shields[n] = v
			n++
		} else {
			// shield broken
			h.log.NewEvent(
				"shield broken",
				glog.LogShieldEvent,
				-1,
			).Write("name", v.Desc()).
				Write("ele", v.Element()).
				Write("expiry", v.Expiry())
			h.events.Emit(event.OnShieldBreak, v)
		}
	}
	h.shields = h.shields[:n]
	return min
}

func (h *Handler) Tick() {
	n := 0
	for _, v := range h.shields {
		if v.Expiry() == *h.f {
			v.OnExpire()
			h.log.NewEvent("shield expired", glog.LogShieldEvent, -1).
				Write("name", v.Desc()).
				Write("hp", v.CurrentHP())
			h.events.Emit(event.OnShieldBreak, v)
		} else {
			h.shields[n] = v
			n++
		}
	}
	h.shields = h.shields[:n]
}
