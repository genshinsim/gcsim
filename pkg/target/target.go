package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

type Target struct {
	Core        *core.Core
	TargetType  combat.TargettableType
	TargetIndex int
	Hitbox      combat.Circle
	Tags        map[string]int

	HPMax     float64
	HPCurrent float64
}

func New(core *core.Core, x, y, r float64) *Target {
	t := &Target{
		Core: core,
	}
	t.Hitbox = *combat.NewCircle(x, y, r)
	t.Tags = make(map[string]int)

	return t
}

func (t *Target) Type() combat.TargettableType { return t.TargetType }
func (t *Target) Index() int                   { return t.TargetIndex }
func (t *Target) SetIndex(ind int)             { t.TargetIndex = ind }
func (t *Target) MaxHP() float64               { return t.HPMax }
func (t *Target) HP() float64                  { return t.HPCurrent }
func (t *Target) Shape() combat.Shape          { return &t.Hitbox }
func (t *Target) SetPos(x, y float64)          { t.Hitbox.SetPos(x, y) }
func (t *Target) Pos() (float64, float64)      { return t.Hitbox.Pos() }
func (t *Target) Kill()                        {} // do nothing

func (t *Target) SetTag(key string, val int) {
	t.Tags[key] = val
}

func (t *Target) GetTag(key string) int {
	return t.Tags[key]
}

func (t *Target) RemoveTag(key string) {
	delete(t.Tags, key)
}
