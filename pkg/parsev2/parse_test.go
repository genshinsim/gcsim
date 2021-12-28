package parse

import (
	"fmt"
	"testing"
)

var pteststring = `
xiangling char lvl=80/90 cons=4 talent=6,9,9;
xiangling add weapon="staff of homa" lvl=80/90 refine=3;
xiangling add set="seal of insulation" count=4;
xiangling add stats hp=4780 atk=311 er=.518 pyro%=0.466 cr=0.311;
xiangling add stats atk%=.0992 cr=.1655 cd=.7282 em=39.64 er=.5510 hp%=.0992 hp=507.88 atk=33.08 def%=.124 def=39.36;

xiangling skill /if=.debuff.res.1.test==1
xiangling burst,skill;
`

func TestParse(t *testing.T) {
	p := New("test", pteststring)
	a, _, err := p.Parse()
	fmt.Println("characters:")
	for _, v := range a.Characters.Profile {
		fmt.Println(v.Base.Key.String())
		//basic stats:
		fmt.Println("\t basics", v.Base)
		fmt.Println("\t weapons", v.Weapon)
		fmt.Println("\t talents", v.Talents)
		fmt.Println("\t sets", v.Sets)
		//pretty print stats
		fmt.Println("\t stats", v.Stats)
	}
	fmt.Println("rotations:")
	for _, v := range a.Rotation {
		fmt.Println(v)
	}
	fmt.Println(a.Targets)
	if err != nil {
		t.Error(err)
	}

}
