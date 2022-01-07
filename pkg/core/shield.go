package core

type ShieldType int

const (
	ShieldCrystallize ShieldType = iota //lasts 15 seconds
	ShieldNoelleSkill
	ShieldNoelleA2
	ShieldZhongliJadeShield
	ShieldDionaSkill
	ShieldBeidouThunderShield
	ShieldXinyanSkill
	ShieldXinyanC2
	ShieldKaeyaC4
	ShieldYanfeiC4
	ShieldBell
	EndShieldType
)

type Shield interface {
	Key() int
	Type() ShieldType
	OnDamage(dmg float64, ele EleType, bonus float64) (float64, bool) //return dmg taken and shield stays
	OnExpire()
	OnOverwrite()
	Expiry() int
	CurrentHP() float64
	Element() EleType
	Desc() string
}

type ShieldHandler interface {
	Add(shd Shield)
	IsShielded() bool
	Get(t ShieldType) Shield
	AddBonus(f func() float64)
	OnDamage(dmg float64, ele EleType) float64
	Count() int
	Tick()
}

type ShieldCtrl struct {
	shields   []Shield
	bonusFunc []func() float64
	core      *Core
}

func NewShieldCtrl(c *Core) *ShieldCtrl {
	return &ShieldCtrl{
		shields:   make([]Shield, 0, EndShieldType),
		bonusFunc: make([]func() float64, 0, 10),
		core:      c,
	}
}

func (s *ShieldCtrl) Count() int { return len(s.shields) }

func (s *ShieldCtrl) IsShielded() bool { return len(s.shields) > 0 }

func (s *ShieldCtrl) Get(t ShieldType) Shield {
	for _, v := range s.shields {
		if v.Type() == t {
			return v
		}
	}
	return nil
}

func (s *ShieldCtrl) AddBonus(f func() float64) {
	s.bonusFunc = append(s.bonusFunc, f)
}

func (s *ShieldCtrl) Add(shd Shield) {
	//we always assume over write of the same type
	ind := -1
	for i, v := range s.shields {
		if v.Type() == shd.Type() {
			ind = i
		}
	}
	if ind > -1 {
		s.core.Log.Debugw("shield overridden", "frame", s.core.F, "event", LogShieldEvent, "frame", s.core.F, "overwrite", true, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
		s.shields[ind].OnOverwrite()
		s.shields[ind] = shd
	} else {
		s.shields = append(s.shields, shd)
		s.core.Log.Debugw("shield added", "frame", s.core.F, "event", LogShieldEvent, "frame", s.core.F, "overwrite", false, "name", shd.Desc(), "hp", shd.CurrentHP(), "ele", shd.Element(), "expiry", shd.Expiry())
	}
	s.core.Events.Emit(OnShielded, shd)
}

func (s *ShieldCtrl) OnDamage(dmg float64, ele EleType) float64 {
	var bonus float64
	//find shield bonuses
	for _, f := range s.bonusFunc {
		bonus += f()
	}
	min := dmg //min of damage taken
	n := 0
	for _, v := range s.shields {
		taken, ok := v.OnDamage(dmg, ele, bonus)
		if taken < min {
			min = taken
		}
		if ok {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
	return min
}

func (s *ShieldCtrl) Tick() {
	n := 0
	for _, v := range s.shields {
		if v.Expiry() == s.core.F {
			v.OnExpire()
			s.core.Log.Debugw("shield expired", "frame", s.core.F, "event", LogShieldEvent, "frame", s.core.F, "name", v.Desc(), "hp", v.CurrentHP())
		} else {
			s.shields[n] = v
			n++
		}
	}
	s.shields = s.shields[:n]
}
