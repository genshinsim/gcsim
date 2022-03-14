package core

import "github.com/genshinsim/gcsim/pkg/coretype"

func (s *Core) NewConstruct(c coretype.Construct, refresh bool) {
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
		s.NewEventBuildMsg(
			coretype.LogConstructEvent,
			-1,
			"construct replaced - new: ", c.Type().String(),
		).Write(
			"key", s.constructs[ind].Key(),
			"prev type", s.constructs[ind].Type(),
			"next type", c.Type(),
		)
		s.constructs[ind].OnDestruct()
		s.constructs[ind] = c

	} else {
		//add this one to the end
		s.constructs = append(s.constructs, c)
		s.NewEventBuildMsg(coretype.LogConstructEvent, -1, "construct created: ", c.Type().String()).Write("key", c.Key(), "type", c.Type())
	}

	//if length > 3, then destruct the beginning ones
	for i := 0; i < len(s.constructs)-3; i++ {
		s.constructs[i].OnDestruct()
		s.NewEventBuildMsg(coretype.LogConstructEvent, -1, "construct destroyed: "+s.constructs[i].Type().String()).Write("key", s.constructs[i].Key(), "type", s.constructs[i].Type())
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

func (s *Core) NewConstructNoLimit(c coretype.Construct, refresh bool) {
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
			s.NewEventBuildMsg(
				coretype.LogConstructEvent, -1,
				"construct destroyed: "+s.consNoLimit[ind].Type().String(),
			).Write(
				"key", s.consNoLimit[ind].Key(),
				"type", s.consNoLimit[ind].Type(),
			)
			s.consNoLimit[ind] = nil

		}
	}
	s.consNoLimit = append(s.consNoLimit, c)
}

func (s *Core) tickConstruct() {
	//clean out expired
	n := 0
	for _, v := range s.constructs {
		if v.Expiry() == s.Frame {
			v.OnDestruct()
			s.NewEventBuildMsg(
				coretype.LogConstructEvent, -1,
				"construct destroyed: "+v.Type().String(),
			).Write(
				"key", v.Key(),
				"type", v.Type(),
			)
		} else {
			s.constructs[n] = v
			n++
		}
	}
	s.constructs = s.constructs[:n]
	n = 0
	for i, v := range s.consNoLimit {
		if v.Expiry() == s.Frame {
			s.consNoLimit[i].OnDestruct()
			s.NewEventBuildMsg(
				coretype.LogConstructEvent, -1,
				"construct destroyed: "+v.Type().String(),
			).Write(
				"key", v.Key(),
				"type", v.Type(),
			)
		} else {
			s.consNoLimit[n] = v
			n++
		}
	}
	s.consNoLimit = s.consNoLimit[:n]

}

//how many of the given
func (s *Core) CountConstruct() int {
	count := 0
	for _, v := range s.constructs {
		count += v.Count()
	}
	for _, v := range s.consNoLimit {
		count += v.Count()
	}
	return count
}

func (s *Core) CountConstructByType(t coretype.GeoConstructType) int {
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

func (s *Core) HasConstruct(key int) bool {
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

func (s *Core) ConstructExpiry(t coretype.GeoConstructType) int {
	expiry := -1
	for _, v := range s.constructs {
		if v.Type() == t {
			if expiry == -1 {
				expiry = v.Expiry()
			} else if expiry > v.Expiry() {
				expiry = v.Expiry()
			}
		}
	}
	for _, v := range s.consNoLimit {
		if v.Type() == t {
			if expiry == -1 {
				expiry = v.Expiry()
			} else if expiry > v.Expiry() {
				expiry = v.Expiry()
			}
		}
	}

	expiry = expiry - s.Frame

	if expiry < 0 {
		return 0
	}

	return expiry
}

//destroy key if exist, return true if destroyed
func (s *Core) DestroyConstruct(key int) bool {
	ok := false
	//clean out expired
	n := 0
	for _, v := range s.constructs {
		if v.Key() == key {
			v.OnDestruct()
			ok = true
			s.NewEventBuildMsg(
				coretype.LogConstructEvent, -1,
				"construct destroyed: "+v.Type().String(),
			).Write(
				"key", v.Key(),
				"type", v.Type(),
			)
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
			s.NewEventBuildMsg(
				coretype.LogConstructEvent, -1,
				"construct destroyed: "+v.Type().String(),
			).Write(
				"key", v.Key(),
				"type", v.Type(),
			)
		} else {
			s.consNoLimit[n] = v
			n++
		}
	}
	s.consNoLimit = s.consNoLimit[:n]
	return ok
}
