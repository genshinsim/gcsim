package core

import "github.com/genshinsim/gcsim/pkg/coretype"

type status struct {
	expiry int
	evt    coretype.LogEvent
}

func (c *Core) StatusDuration(key string) int {
	a, ok := c.status[key]
	if !ok {
		return 0
	}
	if a.expiry > c.Frame {
		return a.expiry - c.Frame
	}
	return 0
}

func (c *Core) AddStatus(key string, dur int) {
	//check if exists
	a, ok := c.status[key]

	//if ok we want to reuse the old evt
	if ok && a.expiry > c.Frame {
		//just reuse the old and update expiry + evt.Ended
		a.expiry = c.Frame + dur
		a.evt.SetEnded(a.expiry)
		c.status[key] = a
		//log an entry for refreshing
		//TODO: this line may not be needed
		if c.Flags.LogDebug {
			c.NewEvent("status refreshed: ", coretype.LogStatusEvent, -1, "key", key, "expiry", c.Frame+dur)
		}
		return
	}

	//otherwise create a new event
	a.evt = c.NewEvent("status added: ", coretype.LogStatusEvent, -1, "key", key, "expiry", c.Frame+dur)
	a.expiry = c.Frame + dur
	a.evt.SetEnded(a.expiry)

	c.status[key] = a
}

func (c *Core) ExtendStatus(key string, dur int) {
	a, ok := c.status[key]

	//do nothing if status doesn't exist
	if !ok || a.expiry <= c.Frame {
		return
	}

	a.expiry += dur
	a.evt.SetEnded(a.expiry)
	c.status[key] = a

	//TODO: this line may not be needed
	if c.Flags.LogDebug {
		c.NewEvent("status refreshed: ", coretype.LogStatusEvent, -1, "key", key, "expiry", a.expiry)
	}
}

func (s *Core) DeleteStatus(key string) {
	//check if it exists first
	a, ok := s.status[key]
	if ok && a.expiry > s.Frame {
		a.evt.SetEnded(s.Frame)
		//TODO: this line may not be needed
		if s.Flags.LogDebug {
			s.NewEvent("status deleted: ", coretype.LogStatusEvent, -1, "key", key)
		}
	}
	delete(s.status, key)
}
