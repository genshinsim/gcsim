package albedo

import "github.com/genshinsim/gcsim/pkg/core/construct"

type skillConstruct struct {
	src    int
	expiry int
	char   *char
}

func (c *char) newConstruct(dur int) *skillConstruct {
	return &skillConstruct{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

func (c *skillConstruct) OnDestruct()                      { c.char.skillActive = false }
func (c *skillConstruct) Key() int                         { return c.src }
func (c *skillConstruct) Type() construct.GeoConstructType { return construct.GeoConstructAlbedoSkill }
func (c *skillConstruct) Expiry() int                      { return c.expiry }
func (c *skillConstruct) IsLimited() bool                  { return true }
func (c *skillConstruct) Count() int                       { return 1 }
