package shield

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// Handler keeps track of all the active shields
// However we do care which is active since:
// - global shields: only active is shielded
// - 1 character shields: active only shieled if is target of shield
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

func (h *Handler) CharacterIsShielded(char, active int) bool {
	for _, v := range h.shields {
		target := v.ShieldTarget()
		if (target == -1 && char == active) || target == char {
			return true
		}
	}
	return false
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
	// we always assume over write of the same type and target
	ind := -1
	for i, v := range h.shields {
		if v.Type() == shd.Type() && v.ShieldTarget() == shd.ShieldTarget() {
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

func (h *Handler) OnDamage(char, active int, dmg float64, ele attributes.Element) float64 {
	// find shield bonuses
	bonus := h.ShieldBonus()
	mintaken := dmg // min of damage taken
	n := 0
	for _, v := range h.shields {
		target := v.ShieldTarget()
		if !((target == -1 && char == active) || target == char) {
			continue
		}
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
		if taken < mintaken {
			mintaken = taken
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
	return mintaken
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
