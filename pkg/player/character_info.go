package player

import "github.com/genshinsim/gcsim/pkg/core/attributes"

func (c *MasterChar) Level() int {
	return c.Base.Level
}

func (c *MasterChar) Zone() ZoneType {
	return c.CharZone
}

func (c *MasterChar) Ele() attributes.Element {
	return c.Base.Element
}
