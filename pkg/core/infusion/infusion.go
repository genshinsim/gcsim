package infusion

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type WeaponInfusion struct {
	Key         string
	Ele         attributes.Element
	Tags        []combat.AttackTag
	Expiry      int
	CanOverride bool
}

const MaxTeamSize = 4

type Handler struct {
	f        *int
	log      glog.Logger
	debug    bool
	infusion [MaxTeamSize]WeaponInfusion
}

func (m *Handler) AddWeaponInfuse(char int, key string, ele attributes.Element, dur int, canOverride bool, tags ...combat.AttackTag) {
	if !m.infusion[char].CanOverride && m.infusion[char].Expiry > *m.f {
		return
	}
	inf := WeaponInfusion{
		Key:         key,
		Ele:         ele,
		Expiry:      *m.f + dur,
		CanOverride: canOverride,
		Tags:        tags,
	}
	m.infusion[char] = inf
}

func (m *Handler) WeaponInfuseIsActive(char int, key string) bool {
	if m.infusion[char].Key != key {
		return false
	}
	//check expiry
	if m.infusion[char].Expiry < *m.f && m.infusion[char].Expiry > -1 {
		return false
	}
	return true
}

func (h *Handler) Infused(char int, a combat.AttackTag) attributes.Element {
	if h.infusion[char].Key != "" {
		ok := false
		for _, v := range h.infusion[char].Tags {
			if v == a {
				ok = true
				break
			}
		}
		if ok {
			if h.infusion[char].Expiry > *h.f || h.infusion[char].Expiry == -1 {
				return h.infusion[char].Ele
			}
		}
	}
	return attributes.NoElement
}
