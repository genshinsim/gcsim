package dummy

import "github.com/genshinsim/gsim/pkg/core"

type Char struct {
	Index   int
	Base    core.CharacterBase
	Weapon  core.WeaponProfile
	Stats   []float64
	Talents core.TalentProfile

	CDReductionFuncs []core.CDAdjust

	Energy    float64
	EnergyMax float64

	HPCurrent float64
	HPMax     float64
}

func NewChar(cfg ...func(*Char)) *Char {
	s := &Char{}
	for _, f := range cfg {
		f(s)
	}
	return s
}

func (c *Char) Init(index int)                                                 {}
func (c *Char) Tick()                                                          {}
func (c *Char) Name() string                                                   { return c.Base.Name }
func (c *Char) CharIndex() int                                                 { return c.Index }
func (c *Char) Ele() core.EleType                                              { return c.Base.Element }
func (c *Char) WeaponClass() core.WeaponClass                                  { return c.Weapon.Class }
func (c *Char) CurrentEnergy() float64                                         { return c.Energy }
func (c *Char) MaxEnergy() float64                                             { return c.EnergyMax }
func (c *Char) TalentLvlSkill() int                                            { return c.Talents.Skill }
func (c *Char) TalentLvlAttack() int                                           { return c.Talents.Attack }
func (c *Char) TalentLvlBurst() int                                            { return c.Talents.Burst }
func (c *Char) HP() float64                                                    { return c.HPCurrent }
func (c *Char) MaxHP() float64                                                 { return c.HPMax }
func (c *Char) ModifyHP(hp float64)                                            { c.HPCurrent = hp }
func (c *Char) Stat(s core.StatType) float64                                   { return c.Stats[s] }
func (c *Char) Attack(p map[string]int) int                                    { return 0 }
func (c *Char) Aimed(p map[string]int) int                                     { return 0 }
func (c *Char) ChargeAttack(p map[string]int) int                              { return 0 }
func (c *Char) HighPlungeAttack(p map[string]int) int                          { return 0 }
func (c *Char) LowPlungeAttack(p map[string]int) int                           { return 0 }
func (c *Char) Skill(p map[string]int) int                                     { return 0 }
func (c *Char) Burst(p map[string]int) int                                     { return 0 }
func (c *Char) Dash(p map[string]int) int                                      { return 0 }
func (c *Char) ActionReady(a core.ActionType, p map[string]int) bool           { return true }
func (c *Char) ActionFrames(a core.ActionType, p map[string]int) int           { return 0 }
func (c *Char) ActionStam(a core.ActionType, p map[string]int) float64         { return 0 }
func (c *Char) AddMod(mod core.CharStatMod)                                    {}
func (c *Char) AddWeaponInfuse(inf core.WeaponInfusion)                        {}
func (c *Char) SetCD(a core.ActionType, dur int)                               {}
func (c *Char) Cooldown(a core.ActionType) int                                 { return 0 }
func (c *Char) ResetActionCooldown(a core.ActionType)                          {}
func (c *Char) ReduceActionCooldown(a core.ActionType, v int)                  {}
func (c *Char) AddCDAdjustFunc(adj core.CDAdjust)                              {}
func (c *Char) Tag(key string) int                                             { return 0 }
func (c *Char) ReceiveParticle(p core.Particle, isActive bool, partyCount int) {}
func (c *Char) QueueParticle(src string, num int, ele core.EleType, delay int) {}
func (c *Char) AddEnergy(e float64)                                            {}
func (c *Char) AddTask(fun func(), name string, delay int)                     {}
func (c *Char) ResetNormalCounter()                                            {}
func (t *Char) QueueDmg(ds *core.Snapshot, delay int)                          {}
func (c *Char) Zone() core.ZoneType                                            { return core.ZoneMondstadt }

func (c *Char) Snapshot(name string, a core.AttackTag, icd core.ICDTag, g core.ICDGroup, st core.StrikeType, e core.EleType, d core.Durability, mult float64) core.Snapshot {
	return core.Snapshot{}
}
