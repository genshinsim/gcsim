package core

type GeoConstructType int

const (
	GeoConstructNingSkill GeoConstructType = iota
	GeoConstructZhongliSkill
	GeoConstructTravellerSkill
	GeoConstructTravellerBurst
	GeoConstructAlbedoSkill
	EndGeoConstructType
)

type Construct interface {
	OnDestruct()
	Key() int
	Type() GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
}

type ConstructHandler interface {
	NewConstruct(c Construct, refresh bool)
	NewNoLimitCons(c Construct, refresh bool)
	ConstructCount() int
	ConstructCountType(t GeoConstructType) int
	Destroy(key int) bool
	HasConstruct(key int) bool
	Tick()
}

type ConstructCtrl struct {
	constructs  []Construct
	consNoLimit []Construct
	core        *Core
}

func NewConstructCtrl(c *Core) *ConstructCtrl {
	return &ConstructCtrl{
		constructs:  make([]Construct, 0, 3),
		consNoLimit: make([]Construct, 0, 3),
		core:        c,
	}
}

func (s *ConstructCtrl) NewConstruct(c Construct, refresh bool) {

	//if refresh, we nil out the old one if any
	ind := len(s.constructs)
	if refresh {
		for i, v := range s.constructs {
			if v.Type() == c.Type() {
				ind = i
			}
		}
	}
	if ind != 0 && ind != len(s.constructs) {
		s.core.Log.Debugw("construct replaced", "event", LogConstructEvent, "frame", s.core.F, "key", s.constructs[ind].Key(), "prev type", s.constructs[ind].Type(), "next type", c.Type())
		s.constructs[ind].OnDestruct()
		s.constructs[ind] = c

	} else {
		//add this one to the end
		s.constructs = append(s.constructs, c)
	}

	//if length > 3, then destruct the beginning ones
	for i := 0; i < len(s.constructs)-3; i++ {
		s.constructs[i].OnDestruct()
		s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", s.constructs[ind].Key(), "type", s.constructs[ind].Type())
		s.constructs[i] = nil
	}

	//clean out any nils
	n := 0
	for _, x := range s.constructs {
		if x != nil {
			s.constructs[n] = x
			n++
		}
	}
	s.constructs = s.constructs[:n]
}

func (s *ConstructCtrl) NewNoLimitCons(c Construct, refresh bool) {
	if refresh {
		ind := len(s.consNoLimit)
		for i, v := range s.consNoLimit {
			//if expired already, set to nil and ignore
			if v.Key() == c.Key() {
				ind = i
			}
		}
		if ind != 0 && ind != len(s.consNoLimit) {
			//destroy the existing by setting expiry
			s.consNoLimit[ind].OnDestruct()
			s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", s.consNoLimit[ind].Key(), "type", s.consNoLimit[ind].Type())
			s.consNoLimit[ind] = nil

		}
	}
	s.consNoLimit = append(s.consNoLimit, c)
}

func (s *ConstructCtrl) Tick() {
	//clean out expired
	n := 0
	for _, v := range s.constructs {
		if v.Expiry() == s.core.F {
			v.OnDestruct()
			s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", v.Key(), "type", v.Type())
		} else {
			s.constructs[n] = v
			n++
		}
	}
	s.constructs = s.constructs[:n]
	n = 0
	for i, v := range s.consNoLimit {
		if v.Expiry() == s.core.F {
			s.consNoLimit[i].OnDestruct()
			s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", v.Key(), "type", v.Type())
		} else {
			s.consNoLimit[n] = v
			n++
		}
	}
	s.consNoLimit = s.consNoLimit[:n]

}

//how many of the given
func (s *ConstructCtrl) ConstructCount() int {
	count := 0
	for _, v := range s.constructs {
		count += v.Count()
	}
	for _, v := range s.consNoLimit {
		count += v.Count()
	}
	return count
}

func (s *ConstructCtrl) ConstructCountType(t GeoConstructType) int {
	count := 0
	for _, v := range s.constructs {
		if v.Type() == t {
			count++
		}
	}
	for _, v := range s.consNoLimit {
		if v.Type() == t {
			count++
		}
	}
	return count
}

func (s *ConstructCtrl) HasConstruct(key int) bool {
	for _, v := range s.constructs {
		if v.Key() == key {
			return true
		}
	}
	for _, v := range s.consNoLimit {
		if v.Key() == key {
			return true
		}
	}
	return false
}

//destroy key if exist, return true if destroyed
func (s *ConstructCtrl) Destroy(key int) bool {
	ok := false
	//clean out expired
	n := 0
	for _, v := range s.constructs {
		if v.Key() == key {
			v.OnDestruct()
			ok = true
			s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", v.Key(), "type", v.Type())
		} else {
			s.constructs[n] = v
			n++
		}
	}
	s.constructs = s.constructs[:n]
	if ok {
		return ok
	}
	n = 0
	for i, v := range s.consNoLimit {
		if v.Key() == key {
			s.consNoLimit[i].OnDestruct()
			ok = true
			s.core.Log.Debugw("destroyed", "event", LogConstructEvent, "frame", s.core.F, "key", v.Key(), "type", v.Type())
		} else {
			s.consNoLimit[n] = v
			n++
		}
	}
	s.consNoLimit = s.consNoLimit[:n]
	return ok
}
