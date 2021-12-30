package character

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func TestBasicAbilUsage(t *testing.T) {
	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	prof := core.CharacterProfile{}
	prof.Base.Element = core.Pyro
	prof.Base.Key = keys.Xiangling
	prof.Stats = make([]float64, core.EndStatType)
	prof.Base.Level = 80
	prof.Base.MaxLevel = 90
	prof.Talents.Attack = 1
	prof.Talents.Skill = 1
	prof.Talents.Burst = 1
	prof.Weapon.Key = "thecatch"
	prof.Weapon.Level = 90
	prof.Weapon.MaxLevel = 90

	x, err := NewTemplateChar(c, prof)
	//cast it to *char so we can access private members
	// this := x.(*char)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	x.SetWeaponKey("thecatch")
	err = x.CalcBaseStats()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, x)
	x.Init(0)

	if !floatApproxEqual(210, x.Base.Atk, 1) {
		t.Errorf("expecting ~210 base atk, got %v", x.Base.Atk)
	}
	if !floatApproxEqual(623, x.Base.Def, 1) {
		t.Errorf("expecting ~623 base def, got %v", x.Base.Def)
	}
	if !floatApproxEqual(10122, x.Base.HP, 1) {
		t.Errorf("expecting ~10122 base hp, got %v", x.Base.HP)
	}
	if !floatApproxEqual(96, x.Stats[core.EM], 1) {
		t.Errorf("expecting ~96 base em, got %v", x.Stats[core.EM])
	}
	if !floatApproxEqual(510, x.Weapon.Atk, 1) {
		t.Errorf("expecting ~510 base atk, got %v", x.Weapon.Atk)
	}
	if !floatApproxEqual(0.459, x.Stats[core.ER], 1) {
		t.Errorf("expecting ~45.9 base er, got %v", x.Stats[core.EM])
	}

}

func floatApproxEqual(expect, result, tol float64) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}
