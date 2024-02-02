package shield

import "github.com/genshinsim/gcsim/pkg/core/attributes"

type Tmpl struct {
	ActorIndex int
	Target     int
	Name       string
	Src        int
	ShieldType Type
	Ele        attributes.Element
	HP         float64
	Expires    int
}

func (t *Tmpl) ShieldOwner() int {
	return t.ActorIndex
}

func (t *Tmpl) ShieldTarget() int {
	return t.Target
}

func (t *Tmpl) Desc() string {
	return t.Name
}

func (t *Tmpl) Element() attributes.Element {
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

func (t *Tmpl) Type() Type {
	return t.ShieldType
}

func (t *Tmpl) ShieldStrength(ele attributes.Element, bonus float64) float64 {
	same := 1.0
	if ele == t.Ele {
		same = 2.5
	}
	if t.Ele == attributes.Geo {
		same = 1.5
	}
	return t.HP * same * (1 + bonus)
}

func (t *Tmpl) OnDamage(dmg float64, ele attributes.Element, bonus float64) (float64, bool) {
	same := 1.0
	if ele == t.Ele {
		same = 2.5
	}
	if ele == attributes.Geo {
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
