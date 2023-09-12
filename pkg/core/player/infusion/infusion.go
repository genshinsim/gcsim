package infusion

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type WeaponInfusion struct {
	Key             string
	Ele             attributes.Element
	Tags            []attacks.AttackTag
	Expiry          float64
	CanBeOverridden bool
}

const MaxTeamSize = 4

type InfusionHandler struct {
	f        *int
	log      glog.Logger
	debug    bool
	infusion [MaxTeamSize]WeaponInfusion
}

func New(f *int, log glog.Logger, debug bool) InfusionHandler {
	return InfusionHandler{
		f:     f,
		log:   log,
		debug: debug,
	}
}

func (i *InfusionHandler) ExtendInfusion(char int, factor, dur float64) {
	// if infusion is active, extend it
	if i.infusion[char].Expiry < float64(*i.f) || i.infusion[char].Expiry == -1 {
		return
	}
	i.infusion[char].Expiry += dur * (1 - factor)
}

func (i *InfusionHandler) AddWeaponInfuse(char int, key string, ele attributes.Element, dur int, canBeOverriden bool, tags ...attacks.AttackTag) {
	if !i.infusion[char].CanBeOverridden && i.infusion[char].Expiry > float64(*i.f) {
		return
	}
	inf := WeaponInfusion{
		Key:             key,
		Ele:             ele,
		Expiry:          float64(*i.f + dur),
		CanBeOverridden: canBeOverriden,
		Tags:            tags,
	}
	if dur == -1 {
		inf.Expiry = -1
	}
	i.infusion[char] = inf
}

func (i *InfusionHandler) WeaponInfuseIsActive(char int, key string) bool {
	if i.infusion[char].Key != key {
		return false
	}
	// check expiry
	if i.infusion[char].Expiry < float64(*i.f) && i.infusion[char].Expiry > -1 {
		return false
	}
	return true
}

func (i *InfusionHandler) Infused(char int, a attacks.AttackTag) attributes.Element {
	if i.infusion[char].Key != "" {
		ok := false
		for _, v := range i.infusion[char].Tags {
			if v == a {
				ok = true
				break
			}
		}
		if ok {
			if i.infusion[char].Expiry > float64(*i.f) || i.infusion[char].Expiry == -1 {
				return i.infusion[char].Ele
			}
		}
	}
	return attributes.NoElement
}
