package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Tmpl struct {
	*reactable.Reactable
	Core        *core.Core
	TargetType  combat.TargettableType
	TargetIndex int
	Hitbox      combat.Circle
	Tags        map[string]int

	HPMax     float64
	HPCurrent float64
}

func (t *Tmpl) Type() combat.TargettableType { return t.TargetType }
func (t *Tmpl) Index() int                   { return t.TargetIndex }
func (t *Tmpl) SetIndex(ind int)             { t.TargetIndex = ind }
func (t *Tmpl) MaxHP() float64               { return t.HPMax }
func (t *Tmpl) HP() float64                  { return t.HPCurrent }
func (t *Tmpl) Shape() combat.Shape          { return &t.Hitbox }
func (t *Tmpl) Kill()                        {} // do nothing

func (t *Tmpl) SetTag(key string, val int) {
	t.Tags[key] = val
}

func (t *Tmpl) GetTag(key string) int {
	return t.Tags[key]
}

func (t *Tmpl) RemoveTag(key string) {
	delete(t.Tags, key)
}
