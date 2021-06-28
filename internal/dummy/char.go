package dummy

import "github.com/genshinsim/gsim/pkg/def"

type Char struct {
	Index   int
	Base    def.CharacterBase
	Weapon  def.WeaponProfile
	Stats   []float64
	Talents def.TalentProfile

	CDReductionFuncs []def.CDAdjust

	Energy    float64
	MaxEnergy float64

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

func (c *Char) Init(index int)                                                {}
func (c *Char) Tick()                                                         {}
func (c *Char) Name() string                                                  { return c.Base.Name }
func (c *Char) CharIndex() int                                                { return c.Index }
func (c *Char) Ele() def.EleType                                              { return c.Base.Element }
func (c *Char) WeaponClass() def.WeaponClass                                  { return c.Weapon.Class }
func (c *Char) CurrentEnergy() float64                                        { return c.Energy }
func (c *Char) TalentLvlSkill() int                                           { return c.Talents.Skill }
func (c *Char) TalentLvlAttack() int                                          { return c.Talents.Attack }
func (c *Char) TalentLvlBurst() int                                           { return c.Talents.Burst }
func (c *Char) HP() float64                                                   { return c.HPCurrent }
func (c *Char) MaxHP() float64                                                { return c.HPMax }
func (c *Char) ModifyHP(hp float64)                                           { c.HPCurrent = hp }
func (c *Char) Stat(s def.StatType) float64                                   { return c.Stats[s] }
func (c *Char) Attack(p map[string]int) int                                   { return 0 }
func (c *Char) Aimed(p map[string]int) int                                    { return 0 }
func (c *Char) ChargeAttack(p map[string]int) int                             { return 0 }
func (c *Char) HighPlungeAttack(p map[string]int) int                         { return 0 }
func (c *Char) LowPlungeAttack(p map[string]int) int                          { return 0 }
func (c *Char) Skill(p map[string]int) int                                    { return 0 }
func (c *Char) Burst(p map[string]int) int                                    { return 0 }
func (c *Char) Dash(p map[string]int) int                                     { return 0 }
func (c *Char) ActionReady(a def.ActionType, p map[string]int) bool           { return true }
func (c *Char) ActionFrames(a def.ActionType, p map[string]int) int           { return 0 }
func (c *Char) ActionStam(a def.ActionType, p map[string]int) float64         { return 0 }
func (c *Char) AddMod(mod def.CharStatMod)                                    {}
func (c *Char) AddWeaponInfuse(inf def.WeaponInfusion)                        {}
func (c *Char) SetCD(a def.ActionType, dur int)                               {}
func (c *Char) Cooldown(a def.ActionType) int                                 { return 0 }
func (c *Char) ResetActionCooldown(a def.ActionType)                          {}
func (c *Char) ReduceActionCooldown(a def.ActionType, v int)                  {}
func (c *Char) AddCDAdjustFunc(adj def.CDAdjust)                              {}
func (c *Char) Tag(key string) int                                            { return 0 }
func (c *Char) ReceiveParticle(p def.Particle, isActive bool, partyCount int) {}
func (c *Char) AddEnergy(e float64)                                           {}
func (c *Char) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d float64, mult float64) def.Snapshot {
	return def.Snapshot{}
}
func (c *Char) ResetNormalCounter() {}
