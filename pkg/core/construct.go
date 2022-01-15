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

var ConstructString = [...]string{
	"NingSkill",
	"ZhongliSkill",
	"TravellerSkill",
	"TravellerBurst",
	"AlbedoSkill",
}

func (c GeoConstructType) String() string {
	return ConstructString[c]
}

type Construct interface {
	OnDestruct()
	Key() int
	Type() GeoConstructType
	Expiry() int
	IsLimited() bool
	Count() int
}

type ConstructHandler interface {
	New(c Construct, refresh bool)
	NewNoLimitCons(c Construct, refresh bool)
	Count() int
	CountByType(t GeoConstructType) int
	Destroy(key int) bool
	Has(key int) bool
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

func (s *ConstructCtrl) New(c Construct, refresh bool) {

	//if refresh, we nil out the old one if any
	ind := -1
	if refresh {
		for i, v := range s.constructs {
			if v.Type() == c.Type() {
				ind = i
			}
		}
	}
	if ind > -1 {
		s.core.Log.Debugw("construct replaced - New: "+c.Type().String(), "event", LogConstructEvent, "frame", s.core.F, "key", s.constructs[ind].Key(), "prev type", s.constructs[ind].Type(), "next type", c.Type())
		s.constructs[ind].OnDestruct()
		s.constructs[ind] = c

	} else {
		//add this one to the end
		s.constructs = append(s.constructs, c)
		s.core.Log.Debugw("construct created: "+c.Type().String(), "event", LogConstructEvent, "frame", s.core.F, "key", c.Key(), "type", c.Type())
	}

	//if length > 3, then destruct the beginning ones
	for i := 0; i < len(s.constructs)-3; i++ {
		s.constructs[i].OnDestruct()
		s.core.Log.Debugw("construct destroyed: "+s.constructs[i].Type().String(), "event", LogConstructEvent, "frame", s.core.F, "key", s.constructs[i].Key(), "type", s.constructs[i].Type())
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
		ind := -1
		for i, v := range s.consNoLimit {
			//if expired already, set to nil and ignore
			if v.Key() == c.Key() {
				ind = i
			}
		}
		if ind > -1 {
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
func (s *ConstructCtrl) Count() int {
	count := 0
	for _, v := range s.constructs {
		count += v.Count()
	}
	for _, v := range s.consNoLimit {
		count += v.Count()
	}
	return count
}

func (s *ConstructCtrl) CountByType(t GeoConstructType) int {
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

func (s *ConstructCtrl) Has(key int) bool {
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
