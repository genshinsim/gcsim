package shield

import "github.com/genshinsim/gcsim/pkg/core"

type Tmpl struct {
	Name       string
	Src        int
	ShieldType core.ShieldType
	Ele        core.EleType
	HP         float64
	Expires    int
}

func (t *Tmpl) Desc() string {
	return t.Name
}

func (t *Tmpl) Element() core.EleType {
	return t.Ele
}

func (t *Tmpl) CurrentHP() float64 {
	return t.HP
}

func (t *Tmpl) Expiry() int {
	return t.Expires
}

func (t *Tmpl) Key() int {
	return t.Src
}

func (t *Tmpl) Type() core.ShieldType {
	return t.ShieldType
}

func (t *Tmpl) OnDamage(dmg float64, ele core.EleType, bonus float64) (float64, bool) {
	same := 1.0
	if ele == t.Ele {
		same = 2.5
	}
	if ele == core.Geo {
		same = 1.5
	}
	block := t.HP * same * (1 + bonus)
	t.HP -= t.HP * (dmg / block)
	if t.HP < 0 {
		t.HP = 0
	}

	taken := dmg - block
	if taken < 0 {
		taken = 0
	}

	return taken, t.HP != 0
}

func (t *Tmpl) OnExpire() {
}

func (t *Tmpl) OnOverwrite() {}
